package tone

import (
	"fmt"
	"math"

	"github.com/leraniode/wondertone/space"
)

// Hex returns the color as a lowercase CSS hex string.
// Gamut-safe: out-of-range OKLCH values are mapped into sRGB before encoding.
func (t Tone) Hex() string {
	r, g, b := t.toSRGB()
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
	rf, gf, bf := t.toSRGB()
	return uint8(math.Round(rf * 255)),
		uint8(math.Round(gf * 255)),
		uint8(math.Round(bf * 255))
}

// RGBFloat returns gamma-encoded sRGB values [0.0–1.0].
func (t Tone) RGBFloat() (r, g, b float64) { return t.toSRGB() }

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
