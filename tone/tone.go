// Package tone is the fundamental unit of wondertone.
//
// A Tone is a named, immutable colour with a human vocabulary:
// Light, Vibrancy, Hue, Energy, and Mood. Under the hood it runs
// on OKLCH with WonderMath perceptual corrections applied at render time.
//
//	import "github.com/leraniode/wondertone/tone"
//
//	spark := tone.New(
//	    tone.Light(75),
//	    tone.Vibrancy(80),
//	    tone.Hue(30),
//	    tone.Energy(0.9),
//	    tone.Named("Spark"),
//	)
//
//	spark.Hex()              // gamut-safe hex output
//	spark.Temperature()      // "warm" / "cool" / "neutral"
//	spark.DerivedMoodValue() // mathematically derived mood
package tone

import (
	"fmt"
	"math"

	"github.com/leraniode/wondertone/space"
)

// Tone is the atomic unit of wondertone — a living, named color.
//
// Developers work with the Wondertone vocabulary:
//   - Light     [0–100]  perceptual lightness. 0=black, 100=white.
//   - Vibrancy  [0–100]  colorfulness relative to gamut max. 0=grey, 100=most vivid possible.
//   - Hue       [0–360)  color angle. 0=red, 120=green, 240=blue.
//   - Energy    [0–1]    aliveness multiplier. 0=muted/sleeping, 1=full/vivid.
//
// Under the hood, everything is OKLCH — perceptually uniform, gamut-safe,
// hue-preserving. The developer never touches L, C, H directly unless they
// want to via FromOKLCH.
//
// A Tone is immutable. Every method returns a new Tone.
type Tone struct {
	// Developer-facing vocabulary — these are what you set and read.
	light    float64 // [0–100]
	vibrancy float64 // [0–100]
	hue      float64 // [0–360)
	energy   float64 // [0–1]

	// Identity
	name string
	mood string

	// Internal OKLCH cache — computed once, never re-derived.
	// l, c are raw OKLCH values stored for math precision.
	// They are computed from light/vibrancy/hue during construction.
	l float64 // OKLCH L [0–1]
	c float64 // OKLCH C [0–~0.37]
	a float64 // alpha [0–1]
}

// --- Option pattern for construction ---

// Option configures a Tone during construction.
type Option func(*Tone)

// Light sets the perceptual lightness [0–100].
// 0 is black. 100 is white. 50 is a true mid-tone.
func Light(v float64) Option {
	return func(t *Tone) {
		t.light = space.Clamp(v, 0, 100)
		t.l = t.light / 100.0
	}
}

// Vibrancy sets how colorful the tone is [0–100], relative to the maximum
// possible saturation at this lightness and hue. 0 is pure grey. 100 is the
// most vivid color the sRGB gamut allows at this exact Light+Hue combination.
//
// This is not raw OKLCH chroma — it is a percentage of the gamut ceiling.
// This means Vibrancy(80) always looks 80% vivid regardless of hue.
func Vibrancy(v float64) Option {
	return func(t *Tone) {
		t.vibrancy = space.Clamp(v, 0, 100)
		// c is computed after all options are applied in New()
		// because it depends on both vibrancy and the final L/H values.
	}
}

// Hue sets the color angle [0–360).
// 0/360=red, 30=orange, 60=yellow, 120=green, 180=cyan, 240=blue, 300=magenta.
func Hue(v float64) Option {
	return func(t *Tone) {
		t.hue = space.NormalizeHue(v)
	}
}

// Energy sets the aliveness multiplier [0–1].
// 1.0 = full vivid (default). 0.5 = half as colorful, same hue. 0.0 = grey.
// Energy scales effective chroma at render time without changing the stored Tone.
// A palette can whisper or shout — same colors, different energy.
func Energy(v float64) Option {
	return func(t *Tone) {
		t.energy = space.Clamp(v, 0, 1)
	}
}

// Named gives the Tone a human name.
func Named(name string) Option {
	return func(t *Tone) { t.name = name }
}

// Moody gives the Tone a mood tag.
// Examples: "vibrant", "serene", "mystical", "focused", "playful", "warm".
func Moody(mood string) Option {
	return func(t *Tone) { t.mood = mood }
}

// Alpha sets the opacity [0–1]. Default is 1.0 (fully opaque).
func Alpha(a float64) Option {
	return func(t *Tone) { t.a = space.Clamp(a, 0, 1) }
}

// --- Constructors ---

// New creates a Tone from options using the Wondertone vocabulary.
//
//	spark := tone.New(
//	    tone.Light(75),
//	    tone.Vibrancy(80),
//	    tone.Hue(30),
//	    tone.Energy(0.9),
//	    tone.Named("Primary Spark"),
//	)
func New(opts ...Option) Tone {
	t := Tone{
		light:    50,
		vibrancy: 100,
		hue:      0,
		energy:   1.0,
		a:        1.0,
		l:        0.5,
	}
	for _, o := range opts {
		o(&t)
	}
	// Compute OKLCH chroma via WonderMath PerceivedChroma:
	// applies gamut ceiling, power-law V^α shaping, and k(H) hue weight.
	t.c = space.PerceivedChroma(t.vibrancy, t.l, t.hue)
	return t
}

// FromOKLCH creates a Tone directly from raw OKLCH values.
// This is the power-user escape hatch — use New() for the normal path.
// L [0–1], C [0–~0.37], H [0–360).
func FromOKLCH(l, c, h float64) Tone {
	l = space.Clamp(l, 0, 1)
	c = math.Max(0, c)
	h = space.NormalizeHue(h)
	maxC := space.MaxChromaForLH(l, h)
	vibrancy := 0.0
	if maxC > 0 {
		vibrancy = space.Clamp((c/maxC)*100, 0, 100)
	}
	return Tone{
		light:    l * 100,
		vibrancy: vibrancy,
		hue:      h,
		energy:   1.0,
		l:        l,
		c:        c,
		a:        1.0,
	}
}

// FromHex parses a CSS hex string (#rgb, #rrggbb, #rrggbbaa).
func FromHex(s string) (Tone, error) {
	r, g, b, a, err := parseHex(s)
	if err != nil {
		return Tone{}, err
	}
	return FromLinearRGB(r, g, b, a), nil
}

// MustFromHex parses a hex string and panics on error.
// Use only for known-good compile-time constants.
func MustFromHex(s string) Tone {
	t, err := FromHex(s)
	if err != nil {
		panic(fmt.Sprintf("wondertone: MustFromHex(%q): %v", s, err))
	}
	return t
}

// FromOKLCHString parses an oklch string: "0.75 0.15 30" or "0.75 0.15 30 / 0.9".
func FromOKLCHString(s string) (Tone, error) {
	l, c, h, a, err := parseOKLCHString(s)
	if err != nil {
		return Tone{}, err
	}
	t := FromOKLCH(l, c, h)
	t.a = space.Clamp(a, 0, 1)
	return t, nil
}

// --- Accessors ---

// Light returns the perceptual lightness [0–100].
func (t Tone) Light() float64 { return t.light }

// Vibrancy returns the colorfulness percentage [0–100].
func (t Tone) Vibrancy() float64 { return t.vibrancy }

// Hue returns the color angle [0–360).
func (t Tone) Hue() float64 { return t.hue }

// Energy returns the aliveness multiplier [0–1].
func (t Tone) Energy() float64 { return t.energy }

// Name returns the tone's human name.
func (t Tone) Name() string { return t.name }

// Mood returns the tone's mood tag.
func (t Tone) Mood() string { return t.mood }

// AlphaValue returns the opacity [0–1].
func (t Tone) AlphaValue() float64 { return t.a }

// --- OKLCH access ---

// OKLCH returns the raw internal OKLCH values: L [0–1], C [0–~0.37], H [0–360).
func (t Tone) OKLCH() (l, c, h float64) {
	return t.l, t.c, t.hue
}

// OKLCHString returns the canonical OKLCH string: "0.750000 0.150000 30.000000".
func (t Tone) OKLCHString() string {
	return formatOKLCHString(t.l, t.c, t.hue, t.a)
}

// RawL returns the value of L
func (t Tone) RawL() float64 {
	return t.l
}

// RawC returns the value of C
func (t Tone) RawC() float64 {
	return t.c
}

// --- Conversion output ---

// Hex returns the color as a lowercase CSS hex string.
// Gamut-safe: out-of-range OKLCH values are mapped into sRGB before encoding.
func (t Tone) Hex() string {
	r, g, b := t.ToSRGB()
	if t.a >= 1.0 {
		return fmt.Sprintf("#%02x%02x%02x",
			uint8(math.Round(r*255)),
			uint8(math.Round(g*255)),
			uint8(math.Round(b*255)),
		)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x",
		uint8(math.Round(r*255)),
		uint8(math.Round(g*255)),
		uint8(math.Round(b*255)),
		uint8(math.Round(t.a*255)),
	)
}

// RGB returns gamma-encoded sRGB values [0–255].
func (t Tone) RGB() (r, g, b uint8) {
	rf, gf, bf := t.ToSRGB()
	return uint8(math.Round(rf * 255)),
		uint8(math.Round(gf * 255)),
		uint8(math.Round(bf * 255))
}

// RGBFloat returns gamma-encoded sRGB values [0.0–1.0].
func (t Tone) RGBFloat() (r, g, b float64) { return t.ToSRGB() }

// String implements fmt.Stringer — returns hex representation.
func (t Tone) String() string { return t.Hex() }

// --- Effective chroma (energy-aware) ---

// EffectiveC returns the chroma that will actually be used for rendering.
// Uses the WonderMath Stevens' power law: C * E^γ (γ≈0.7).
// Energy=0.5 now feels half as alive, not 60% as alive.
// The stored C and Energy are never modified.
func (t Tone) EffectiveC() float64 {
	return space.EffectiveChroma(t.c, t.energy)
}

// --- Manipulation (all immutable — returns new Tone) ---

// WithLight returns a new Tone with the given lightness [0–100].
func (t Tone) WithLight(v float64) Tone {
	return New(
		Light(v), Vibrancy(t.vibrancy), Hue(t.hue), Energy(t.energy),
		Named(t.name), Moody(t.mood), Alpha(t.a),
	)
}

// WithVibrancy returns a new Tone with the given vibrancy [0–100].
func (t Tone) WithVibrancy(v float64) Tone {
	return New(
		Light(t.light), Vibrancy(v), Hue(t.hue), Energy(t.energy),
		Named(t.name), Moody(t.mood), Alpha(t.a),
	)
}

// WithHue returns a new Tone with the given hue [0–360).
func (t Tone) WithHue(h float64) Tone {
	return New(
		Light(t.light), Vibrancy(t.vibrancy), Hue(h), Energy(t.energy),
		Named(t.name), Moody(t.mood), Alpha(t.a),
	)
}

// WithEnergy returns a new Tone with the given energy [0–1].
func (t Tone) WithEnergy(e float64) Tone {
	nt := t
	nt.energy = space.Clamp(e, 0, 1)
	return nt
}

// WithName returns a new Tone with the given name.
func (t Tone) WithName(name string) Tone {
	nt := t
	nt.name = name
	return nt
}

// WithMood returns a new Tone with the given mood.
func (t Tone) WithMood(mood string) Tone {
	nt := t
	nt.mood = mood
	return nt
}

// WithAlpha returns a new Tone with the given alpha [0–1].
func (t Tone) WithAlpha(a float64) Tone {
	nt := t
	nt.a = space.Clamp(a, 0, 1)
	return nt
}

// Lighten returns a new Tone with lightness increased by amount [0–100].
func (t Tone) Lighten(amount float64) Tone {
	return t.WithLight(t.light + amount)
}

// Darken returns a new Tone with lightness decreased by amount [0–100].
func (t Tone) Darken(amount float64) Tone {
	return t.WithLight(t.light - amount)
}

// Saturate returns a new Tone with vibrancy increased by amount [0–100].
func (t Tone) Saturate(amount float64) Tone {
	return t.WithVibrancy(t.vibrancy + amount)
}

// Desaturate returns a new Tone with vibrancy decreased by amount [0–100].
func (t Tone) Desaturate(amount float64) Tone {
	return t.WithVibrancy(t.vibrancy - amount)
}

// Rotate returns a new Tone with hue rotated by degrees.
// Accepts negative values (counter-clockwise).
func (t Tone) Rotate(degrees float64) Tone {
	return t.WithHue(space.NormalizeHue(t.hue + degrees))
}

// Complement returns the Tone directly opposite on the hue wheel (+180°).
func (t Tone) Complement() Tone {
	return t.Rotate(180).WithName("").WithMood("")
}

// --- Intelligence ---

// IsLight reports whether the Tone is perceptually light (Light > 50).
func (t Tone) IsLight() bool { return t.light > 50 }

// IsDark reports whether the Tone is perceptually dark (Light <= 50).
func (t Tone) IsDark() bool { return t.light <= 50 }

// Temperature returns "warm", "cool", or "neutral".
// Upgraded in v0.2: uses WonderMath continuous formula instead of hue-range
// lookup — chroma and lightness now modulate the reading.
func (t Tone) Temperature() string {
	tv := space.TemperatureValue(t.hue, t.c, t.l)
	return space.TemperatureLabel(tv)
}

// TemperatureScalar returns the continuous warm↔cool value T ∈ [-1, +1].
// +1 = maximally warm, -1 = maximally cool, 0 = neutral.
// More precise than Temperature() which returns a label.
func (t Tone) TemperatureScalar() float64 {
	return space.TemperatureValue(t.hue, t.c, t.l)
}

// DerivedMoodValue returns the mathematically derived mood string.
// Computed from valence and arousal — independent of the stored Mood() tag.
// Use Mood() for the display label (manual override takes precedence).
func (t Tone) DerivedMoodValue() string {
	s := space.NormalizedSaturation(t.c, t.l, t.hue)
	tv := space.TemperatureValue(t.hue, t.c, t.l)
	val := space.Valence(tv, t.l, s)
	aro := space.Arousal(s, t.energy, tv)
	return space.DerivedMood(val, aro, tv)
}

// Valence returns the emotional valence of this tone ∈ [-1, +1].
// +1 = positive (bright, warm, vivid). -1 = negative (dark, cool, muted).
func (t Tone) ValenceValue() float64 {
	s := space.NormalizedSaturation(t.c, t.l, t.hue)
	tv := space.TemperatureValue(t.hue, t.c, t.l)
	return space.Valence(tv, t.l, s)
}

// ArousalValue returns the emotional arousal of this tone ∈ [-1, +1].
// +1 = activated (vivid, energetic). -1 = calm (muted, quiet).
func (t Tone) ArousalValue() float64 {
	s := space.NormalizedSaturation(t.c, t.l, t.hue)
	tv := space.TemperatureValue(t.hue, t.c, t.l)
	return space.Arousal(s, t.energy, tv)
}

// --- Accessibility ---

// Luminance returns the WCAG 2.1 relative luminance [0–1].
func (t Tone) Luminance() float64 {
	r, g, b := t.ToSRGB()
	lin := func(v float64) float64 {
		if v <= 0.04045 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(r) + 0.7152*lin(g) + 0.0722*lin(b)
}

// ContrastWith returns the WCAG 2.1 contrast ratio against another Tone [1–21].
func (t Tone) ContrastWith(other Tone) float64 {
	l1, l2 := t.Luminance(), other.Luminance()
	lighter, darker := math.Max(l1, l2), math.Min(l1, l2)
	return (lighter + 0.05) / (darker + 0.05)
}

// PassesAA reports whether this Tone meets WCAG AA (4.5:1) against bg.
func (t Tone) PassesAA(bg Tone) bool { return t.ContrastWith(bg) >= 4.5 }

// PassesAAA reports whether this Tone meets WCAG AAA (7.0:1) against bg.
func (t Tone) PassesAAA(bg Tone) bool { return t.ContrastWith(bg) >= 7.0 }

// EnsureContrast returns a Tone adjusted to meet the target contrast ratio
// against bg. Adjusts lightness only — hue, vibrancy and energy are preserved.
// level: "AA" (4.5:1) or "AAA" (7.0:1).
func (t Tone) EnsureContrast(bg Tone, level string) Tone {
	target := 4.5
	if level == "AAA" {
		target = 7.0
	}
	if t.ContrastWith(bg) >= target {
		return t
	}
	// Binary search on lightness.
	// Go lighter if we are dark, go darker if we are light.
	goLighter := t.light < 50
	lo, hi := 0.0, 100.0
	best := t
	for i := 0; i < 64; i++ {
		mid := (lo + hi) / 2
		var candidate Tone
		if goLighter {
			candidate = t.WithLight(mid)
		} else {
			candidate = t.WithLight(100 - mid)
		}
		if candidate.ContrastWith(bg) >= target {
			best = candidate
			hi = mid
		} else {
			lo = mid
		}
	}
	return best
}

// --- Tone Scale ---

// Scale returns the 12-step perceptual tone scale for this Tone.
// Step 9 is the "pure" tone — closest to the original.
func (t Tone) Scale() ToneScale {
	return generateScale(t)
}

// Step returns a single step from the 12-step scale [1–12].
// 1=lightest, 12=darkest.
func (t Tone) Step(n int) Tone {
	return t.Scale().Step(n)
}

// --- Equality ---

// Equal reports whether two Tones are perceptually identical
// within floating-point tolerance.
func (t Tone) Equal(other Tone) bool {
	const eps = 1e-4
	return math.Abs(t.l-other.l) < eps &&
		math.Abs(t.c-other.c) < eps &&
		math.Abs(t.hue-other.hue) < eps &&
		math.Abs(t.a-other.a) < eps &&
		math.Abs(t.energy-other.energy) < eps
}

// --- Internal ---

// toSRGB converts to gamma-encoded sRGB [0–1], gamut-safe.
// Full WonderMath pipeline applied at render time:
//  1. CorrectedHue       — fix blue-purple drift
//  2. EffectiveChroma    — Stevens' power law energy scaling
//  3. EffectiveLightness — subtle glow at high energy
func (t Tone) ToSRGB() (r, g, b float64) {
	h := space.CorrectedHue(t.hue, t.c, t.l)
	c := space.EffectiveChroma(t.c, t.energy)
	l := space.EffectiveLightness(t.l, t.energy)
	l, c, h = space.ToGamutSafe(l, c, h)
	lr, lg, lb := space.OKLCHToLinearRGB(l, c, h)
	return space.LinearToSRGB(lr), space.LinearToSRGB(lg), space.LinearToSRGB(lb)
}

// fromLinearRGB constructs a Tone from linear sRGB [0–1] values.
func FromLinearRGB(r, g, b, a float64) Tone {
	l, c, h := space.LinearRGBToOKLCH(r, g, b)
	t := FromOKLCH(l, c, h)
	t.a = space.Clamp(a, 0, 1)
	return t
}
