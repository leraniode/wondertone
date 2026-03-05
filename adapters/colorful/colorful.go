// Package colorful provides interoperability between wondertone and
// github.com/lucasb-eyer/go-colorful.
//
// This is the adapter for using wondertone with go-colorful — import it only if your project already
// uses go-colorful. The wondertone core has zero dependency on go-colorful.
//
// Usage:
//
//	import (
//	    tone      "github.com/leraniode/wondertone/core"
//	    wcolorful "github.com/leraniode/wondertone/adapters/colorful"
//	)
//
//	// wondertone → go-colorful
//	cf := wcolorful.ToColorful(myTone)
//	lab, _ := cf.Lab()
//
//	// go-colorful → wondertone
//	t := wcolorful.FromColorful(cf)
//	fmt.Println(t.Hex())
package colorful

import (
	"github.com/lucasb-eyer/go-colorful"
	tone "github.com/leraniode/wondertone/core"
)

// ToColorful converts a wondertone Tone to a go-colorful Color.
// The conversion goes through sRGB — gamut-safe, energy-aware.
// This is the rendered value: Energy is applied, gamut is clamped.
func ToColorful(t tone.Tone) colorful.Color {
	r, g, b := t.RGBFloat()
	return colorful.Color{R: r, G: g, B: b}
}

// ToColorfulOKLCH converts a wondertone Tone to go-colorful via OKLCH directly.
// This preserves the stored OKLCH values (Energy NOT applied — raw tone truth).
// Use this when you want to continue working in color space, not render.
func ToColorfulOKLCH(t tone.Tone) colorful.Color {
	l, c, h := t.OKLCH()
	return colorful.OkLch(l, c, h)
}

// FromColorful converts a go-colorful Color to a wondertone Tone.
// Assumes the Color is sRGB (which is go-colorful's native space).
func FromColorful(c colorful.Color) tone.Tone {
	t, _ := tone.FromHex(c.Hex())
	return t
}

// FromColorfulOKLCH converts a go-colorful Color to a wondertone Tone
// via OKLCH — preserves the color space values without going through hex.
func FromColorfulOKLCH(c colorful.Color) tone.Tone {
	l, ch, h := c.OkLch()
	return tone.FromOKLCH(l, ch, h)
}

// BlendColorful mixes two wondertone Tones using go-colorful's OKLab blending.
// This is an alternative to tone.Mix() — results should be nearly identical
// since both use OKLab, but this lets you use go-colorful's implementation
// if you prefer it for consistency with the rest of your codebase.
func BlendColorful(a, b tone.Tone, t float64) tone.Tone {
	ca := ToColorful(a)
	cb := ToColorful(b)
	blended := ca.BlendOkLab(cb, t)
	return FromColorfulOKLCH(blended)
}
