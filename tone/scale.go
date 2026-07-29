package tone

import "github.com/leraniode/wondertone/space"

// ToneScale is a 12-step perceptual ladder from a single base Tone.
// Hue never drifts. Every step is in sRGB gamut.
//
// Step roles (1=lightest, 12=darkest):
//
//	 1  Page background         — barely tinted white
//	 2  Subtle background       — stripes, alternating rows
//	 3  UI element background   — cards, inputs
//	 4  Hovered element bg      — hover state
//	 5  Active/selected bg      — selected state
//	 6  Subtle border           — separators
//	 7  UI border               — input borders
//	 8  Strong border           — focus rings
//	 9  Solid bg                — buttons, badges — base tone lives here
//	10  Hovered solid           — button hover
//	11  Text                    — readable on light backgrounds
//	12  High-contrast text      — headings, emphasis
type ToneScale [12]Tone

// Semantic accessors — named by their UI role.
func (s ToneScale) Background() Tone        { return s[0] }
func (s ToneScale) SubtleBackground() Tone  { return s[1] }
func (s ToneScale) ElementBackground() Tone { return s[2] }
func (s ToneScale) HoveredBackground() Tone { return s[3] }
func (s ToneScale) ActiveBackground() Tone  { return s[4] }
func (s ToneScale) SubtleBorder() Tone      { return s[5] }
func (s ToneScale) Border() Tone            { return s[6] }
func (s ToneScale) StrongBorder() Tone      { return s[7] }
func (s ToneScale) Solid() Tone             { return s[8] }
func (s ToneScale) HoveredSolid() Tone      { return s[9] }
func (s ToneScale) Text() Tone              { return s[10] }
func (s ToneScale) HighContrastText() Tone  { return s[11] }

// Step returns a tone by 1-based step number [1–12]. Clamped at edges.
func (s ToneScale) Step(n int) Tone {
	if n < 1 {
		n = 1
	}
	if n > 12 {
		n = 12
	}
	return s[n-1]
}

// All returns all 12 tones as a slice (lightest to darkest).
func (s ToneScale) All() []Tone {
	out := make([]Tone, 12)
	for i, t := range s {
		out[i] = t
	}
	return out
}

// lightnessTable: perceptual lightness per step [0–1].
// Step 9 (index 8) is where the base color lives.
var lightnessTable = [12]float64{
	0.985, 0.955, 0.915, 0.855, 0.770, 0.680,
	0.570, 0.445, 0.360, 0.270, 0.180, 0.100,
}

// vibrantTable: fraction of max chroma to use per step [0–1].
// Center steps are maximally vivid. Extremes are desaturated —
// airy at the light end, refined at the dark end.
var vibrantTable = [12]float64{
	0.08, 0.12, 0.20, 0.35, 0.50, 0.65,
	0.80, 0.92, 1.00, 0.95, 0.80, 0.65,
}

// generateScale builds the 12-step ToneScale for a base Tone.
// Hue is locked from the base. Each step computes its own maxChroma
// so vibrancy is always relative to what is actually possible at that lightness.
func generateScale(base Tone) ToneScale {
	var scale ToneScale
	for i := 0; i < 12; i++ {
		L := lightnessTable[i]
		H := base.hue

		// Max chroma at this exact (L, H) — varies per hue.
		maxC := space.MaxChromaForLH(L, H)

		// Scale by both the vibrancy table and the base tone's own vibrancy.
		// A muted base (vibrancy=40) produces a muted scale. Vivid produces vivid.
		effectiveMaxC := maxC * (base.vibrancy / 100.0)
		c := effectiveMaxC * vibrantTable[i]

		light := L * 100
		vibrancy := 0.0
		if maxC > 0 {
			vibrancy = space.Clamp((c/maxC)*100, 0, 100)
		}

		scale[i] = Tone{
			light:    light,
			vibrancy: vibrancy,
			hue:      H,
			energy:   base.energy,
			l:        L,
			c:        c,
			a:        base.a,
		}
	}
	return scale
}
