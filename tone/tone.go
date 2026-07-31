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
	return fromLinearRGB(r, g, b, a), nil
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
