// Package mix provides OKLab-space colour mixing, gradients, and blending.
package mix

import (
	"fmt"

	"github.com/leraniode/wondertone/space"
	"github.com/leraniode/wondertone/tone"
)

// Mix blends two Tones in OKLab space at ratio t [0–1].
// t=0 returns a, t=1 returns b, t=0.5 is the perceptual midpoint.
func Mix(a, b tone.Tone, t float64) tone.Tone {
	t = space.Clamp(t, 0, 1)
	if t == 0 {
		return a
	}
	if t == 1 {
		return b
	}

	La, aa, ba := space.OKLCHToOKLab(a.RawL(), a.RawC(), a.Hue())
	Lb, ab, bb := space.OKLCHToOKLab(b.RawL(), b.RawC(), b.Hue())

	L := La + t*(Lb-La)
	aL := aa + t*(ab-aa)
	bL := ba + t*(bb-ba)

	l, c, h := space.OKLabToOKLCH(L, aL, bL)
	l, c, h = space.ToGamutSafe(l, c, h)

	var vibrancy float64
	maxC := space.MaxChromaForLH(l, h)
	if maxC > 0 {
		vibrancy = space.Clamp((c/maxC)*100, 0, 100)
	}

	return tone.FromOKLCH(l, c, h).
		WithLight(l * 100).
		WithVibrancy(vibrancy)
}

// Gradient produces n perceptually uniform steps from start to end in OKLab space.
func Gradient(start, end tone.Tone, n int) ([]tone.Tone, error) {
	if n < 2 {
		return nil, fmt.Errorf("wondertone/mix: Gradient needs at least 2 steps")
	}
	result := make([]tone.Tone, n)
	for i := 0; i < n; i++ {
		result[i] = Mix(start, end, float64(i)/float64(n-1))
	}
	return result, nil
}

// Blend mixes multiple Tones in OKLab space with per-tone weights.
// Weights are normalised automatically — they do not need to sum to 1.
func Blend(tones []tone.Tone, weights []float64) (tone.Tone, error) {
	if len(tones) == 0 {
		return tone.Tone{}, fmt.Errorf("wondertone/mix: Blend requires at least one tone")
	}
	if len(tones) != len(weights) {
		return tone.Tone{}, fmt.Errorf("wondertone/mix: tones and weights must be the same length")
	}

	var sumW, sumL, sumA, sumB float64
	for i, t := range tones {
		w := weights[i]
		if w < 0 {
			return tone.Tone{}, fmt.Errorf("wondertone/mix: weight %d is negative", i)
		}
		L, a, b := space.OKLCHToOKLab(t.RawL(), t.RawC(), t.Hue())
		sumL += w * L
		sumA += w * a
		sumB += w * b
		sumW += w
	}
	if sumW == 0 {
		return tone.Tone{}, fmt.Errorf("wondertone/mix: weights sum to zero")
	}

	l, c, h := space.OKLabToOKLCH(sumL/sumW, sumA/sumW, sumB/sumW)
	l, c, h = space.ToGamutSafe(l, c, h)
	return tone.FromOKLCH(l, c, h), nil
}

// Harmonize returns a set of Tones related by the given harmonic scheme.
// All harmony is computed by hue rotation in OKLCH — lightness and vibrancy
// are preserved from the base Tone.
//
// scheme: "complement", "triadic", "analogous", "split", "tetradic"
func Harmonize(base tone.Tone, scheme string) ([]tone.Tone, error) {
	switch scheme {
	case "complement":
		return []tone.Tone{base, base.Rotate(180)}, nil
	case "triadic":
		return []tone.Tone{base, base.Rotate(120), base.Rotate(240)}, nil
	case "analogous":
		return []tone.Tone{base.Rotate(-30), base, base.Rotate(30)}, nil
	case "split":
		return []tone.Tone{base, base.Rotate(150), base.Rotate(210)}, nil
	case "tetradic":
		return []tone.Tone{base, base.Rotate(90), base.Rotate(180), base.Rotate(270)}, nil
	default:
		return nil, fmt.Errorf(
			"wondertone: unknown harmony scheme %q (want: complement, triadic, analogous, split, tetradic)",
			scheme,
		)
	}
}
