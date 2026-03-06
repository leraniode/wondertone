package core

// wondermath.go — WonderSpace colour math layer.
//
// This is the layer above OKLCH. Every formula here corrects a known weakness
// in raw OKLCH or adds a new perceptual dimension that OKLCH does not have.
//
// All functions are pure — no side effects, no state. Constants are tunable.
// See docs for the research behind each formula.
//
// WonderMath pipeline (in order of application):
//
//  1. CorrectedHue        — fix blue-purple drift at high chroma
//  2. PerceivedChroma     — vibrancy + hue factor k(H) + Energy^γ
//  3. TemperatureValue    — continuous warm↔cool scalar
//  4. Valence + Arousal   — mood as a 2D vector
//  5. DerivedMood         — map valence/arousal vector to named mood

import "math"

// --- WonderMath constants ---
// These are the tunable parameters. Change here to affect all math globally.

const (
	// Hue correction (Section 1 — blue drift fix)
	wmHueCorrH0 = 250.0 // centre of problematic hue region (blue-purple)
	wmHueCorrW  = 30.0  // half-width of Gaussian correction in degrees
	wmHueCorrA  = -3.0  // peak correction magnitude, degrees (negative = shift away from purple)

	// Vibrancy → chroma (Section 2)
	wmVibrancyAlpha = 0.9 // power law exponent on V — shapes low-V perception

	// Energy (Section 6 — Stevens' power law)
	wmEnergyGamma  = 0.7  // perceptual exponent: E^γ instead of linear E
	wmEnergyLambda = 0.04 // lightness glow coefficient: subtle L boost at high energy

	// Temperature (Section 3)
	wmTempHWarm = 50.0 // warm hue axis centre (covers red-orange)
	wmTempWH    = 0.7  // weight: hue cosine (dominant)
	wmTempWC    = 0.2  // weight: chroma contribution
	wmTempWL    = 0.1  // weight: lightness shift

	// Mood coefficients (Section 4 — Valdez & Mehrabian style)
	wmValenceA1 = 0.1 // temperature → valence
	wmValenceA2 = 0.7 // lightness → valence (light = positive)
	wmValenceA3 = 0.2 // saturation → valence

	wmArousalB1 = 0.6 // saturation → arousal (vivid = excited)
	wmArousalB2 = 0.4 // energy → arousal
	wmArousalB3 = 0.0 // temperature → arousal (currently unused — set to 0)
)

// --- k(H): perceptual hue weight table ---
// Controls Section 2 — how much chroma each hue needs to appear equally vivid.
// Values < 1.0: hue appears stronger than its chroma suggests (reduce it).
// Values > 1.0: hue appears weaker than its chroma suggests (boost it).
//
// Control points: [hue_degrees, k_value]
// Interpolated via cubic Catmull-Rom spline between points.
//
// Initial estimates derived from known perceptual properties:
//   - Yellow (~60°) is naturally very vivid — reduce slightly
//   - Green (~142°) is neutral — keep near 1.0
//   - Blue (~240°) appears weaker perceptually — boost
//   - Red (~0°/360°) is vivid — reduce slightly
var kHueControlPoints = [][2]float64{
	{0, 0.90},   // red — vivid, slight reduction
	{30, 0.92},  // orange
	{60, 0.85},  // yellow — strongest hue, needs most reduction
	{90, 0.95},  // yellow-green
	{120, 0.98}, // green — near neutral
	{150, 1.00}, // cyan-green — neutral reference
	{180, 1.02}, // cyan
	{210, 1.05}, // blue-cyan
	{240, 1.10}, // blue — appears weaker, boost
	{270, 1.08}, // blue-purple
	{300, 0.95}, // magenta
	{330, 0.92}, // red-magenta
	{360, 0.90}, // red (wraps)
}

// kHue returns the perceptual hue weight k(H) for a given hue angle.
// Interpolated via linear interpolation between control points.
// (Catmull-Rom spline is ideal but linear is sufficient for initial tuning.)
func kHue(h float64) float64 {
	h = normalizeHue(h)
	pts := kHueControlPoints

	// Find surrounding control points
	for i := 0; i < len(pts)-1; i++ {
		h0, k0 := pts[i][0], pts[i][1]
		h1, k1 := pts[i+1][0], pts[i+1][1]
		if h >= h0 && h <= h1 {
			t := (h - h0) / (h1 - h0)
			return k0 + t*(k1-k0)
		}
	}
	return 1.0
}

// --- 2.1 CorrectedHue ---

// CorrectedHue applies the WonderMath blue-drift correction.
// Returns H' — the perceptually corrected hue.
//
// Formula:
//
//	H' = wrap360(H + A * exp(-(H-H0)² / (2w²)) * (C / C_max))
//
// The Gaussian is chroma-weighted: grey tones get zero correction.
// Only vivid blues near H≈250° are nudged back from purple.
func CorrectedHue(h, c, l float64) float64 {
	cMax := maxChromaForLH(l, h)
	if cMax <= 0 {
		return h
	}
	chromaWeight := c / cMax
	gaussian := math.Exp(-math.Pow(h-wmHueCorrH0, 2) / (2 * wmHueCorrW * wmHueCorrW))
	delta := wmHueCorrA * gaussian * chromaWeight
	result := math.Mod(h+delta, 360)
	if result < 0 {
		result += 360
	}
	// Ensure strict [0, 360) — floating point can land exactly on 360
	if result >= 360 {
		result -= 360
	}
	return result
}

// --- 2.2 PerceivedChroma ---

// PerceivedChroma converts Vibrancy [0–100] to raw OKLCH chroma C using
// the full WonderMath pipeline:
//
//  1. Gamut-relative base:   rawC = C_max * (V/100)^α
//  2. Hue weight correction: C = rawC * k(H)
//
// Energy is NOT applied here — it is applied separately in EffectiveC.
// This is the stored chroma value: what the tone IS, not how it renders.
func PerceivedChroma(vibrancy, l, h float64) float64 {
	cMax := maxChromaForLH(l, h)
	if cMax <= 0 {
		return 0
	}
	v := clamp(vibrancy/100.0, 0, 1)
	rawC := cMax * math.Pow(v, wmVibrancyAlpha)
	return rawC * kHue(h)
}

// --- 2.3 EffectiveChroma (Energy, Stevens' power law) ---

// EffectiveChroma applies the perceptually linear Energy scaling to chroma.
//
// Formula:
//
//	C_effective = C_base * E^γ
//
// γ ≈ 0.7 (Stevens' power law exponent for colour saturation).
// This makes Energy=0.5 feel like "half as alive", not "60% as alive".
func EffectiveChroma(c, energy float64) float64 {
	if energy <= 0 {
		return 0
	}
	if energy >= 1 {
		return c
	}
	return c * math.Pow(energy, wmEnergyGamma)
}

// EffectiveLightness applies the subtle lightness glow at high energy.
// High-energy tones get a small brightness boost — they glow slightly.
//
// Formula:
//
//	L_effective = L + λ * (E^γ - 1)
//
// λ = 0.04 means full-energy tone gets +4% lightness max.
// At Energy=0.5: L_effective = L + 0.04 * (0.5^0.7 - 1) ≈ L - 0.015
func EffectiveLightness(l, energy float64) float64 {
	glow := wmEnergyLambda * (math.Pow(clamp(energy, 0, 1), wmEnergyGamma) - 1)
	return clamp(l+glow, 0, 1)
}

// --- 2.4 TemperatureValue ---

// TemperatureValue returns a continuous warm↔cool scalar T ∈ [-1, +1].
//
//	+1.0 = maximally warm (vivid red-orange)
//	 0.0 = neutral
//	-1.0 = maximally cool (vivid blue-cyan)
//
// Formula:
//
//	raw = cos((H - H_warm) * π/180)
//	T = clamp(w_h*raw + w_c*(C/C_max) + w_l*(L-0.5), -1, 1)
//
// Hue dominates (w_h=0.7). Chroma boosts the reading (more vivid = more
// definitely warm or cool). Lightness nudges slightly (dark oranges read
// slightly cooler than bright ones).
func TemperatureValue(h, c, l float64) float64 {
	cMax := maxChromaForLH(l, h)
	chromaNorm := 0.0
	if cMax > 0 {
		chromaNorm = c / cMax
	}
	raw := math.Cos((h - wmTempHWarm) * math.Pi / 180.0)
	t := wmTempWH*raw + wmTempWC*chromaNorm + wmTempWL*(l-0.5)
	return clamp(t, -1, 1)
}

// TemperatureLabel maps a continuous temperature value to "warm", "cool",
// or "neutral". Replaces the old hue-range lookup.
func TemperatureLabel(tv float64) string {
	switch {
	case tv > 0.15:
		return "warm"
	case tv < -0.15:
		return "cool"
	default:
		return "neutral"
	}
}

// --- 2.5 Valence and Arousal ---

// Valence returns the emotional valence of a tone ∈ [-1, +1].
//
//	+1.0 = most positive (bright, light, warm)
//	-1.0 = most negative (dark, cool, desaturated)
//
// Formula (Valdez & Mehrabian style):
//
//	valence = clamp(a1*T + a2*L + a3*S, -1, 1)
func Valence(tv, l, saturation float64) float64 {
	v := wmValenceA1*tv + wmValenceA2*l + wmValenceA3*saturation
	return clamp(v, -1, 1)
}

// Arousal returns the emotional arousal/energy level of a tone ∈ [-1, +1].
//
//	+1.0 = highly activated (vivid, high energy)
//	-1.0 = calm/sleepy (muted, low energy)
//
// Formula:
//
//	arousal = clamp(b1*S + b2*E + b3*T, -1, 1)
func Arousal(saturation, energy, tv float64) float64 {
	a := wmArousalB1*saturation + wmArousalB2*energy + wmArousalB3*tv
	return clamp(a, -1, 1)
}

// --- 2.6 DerivedMood ---

// DerivedMood maps a valence/arousal vector to a named mood.
// This replaces the manual Moody() tag with a mathematically derived value.
// The manually set mood still takes precedence as a display override.
//
// Mood regions (valence, arousal):
//
//	vivid    — high arousal, mid-high valence
//	playful  — high arousal, high valence, warm
//	urgent   — high arousal, low-mid valence
//	focused  — mid arousal, mid valence
//	warm     — mid arousal, high valence, warm temperature
//	mystical — mid arousal, low valence, cool
//	calm     — low arousal, mid valence
//	deep     — low arousal, low valence
//	airy     — low arousal, high valence
func DerivedMood(valence, arousal, temperatureValue float64) string {
	switch {
	case arousal > 0.5 && valence > 0.2:
		if temperatureValue > 0.1 {
			return "playful"
		}
		return "vivid"
	case arousal > 0.5 && valence <= 0.2:
		return "urgent"
	case arousal > 0.2 && valence > 0.4 && temperatureValue > 0.1:
		return "warm"
	case arousal > 0.2 && valence < -0.1 && temperatureValue < -0.1:
		return "mystical"
	case arousal > 0.2:
		return "focused"
	case arousal <= 0.2 && valence > 0.5:
		return "airy"
	case arousal <= 0.2 && valence < -0.2:
		return "deep"
	default:
		return "calm"
	}
}

// --- Helpers ---

// normalizedSaturation returns C / C_max — the normalised saturation [0,1]
// used in the mood formulas. Same as Vibrancy/100 but computed from raw C.
func normalizedSaturation(c, l, h float64) float64 {
	cMax := maxChromaForLH(l, h)
	if cMax <= 0 {
		return 0
	}
	return clamp(c/cMax, 0, 1)
}
