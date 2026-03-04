package core

import (
	"fmt"
	"math"
)

// Mix blends two Tones in OKLab space at ratio t [0–1].
// t=0 returns a, t=1 returns b, t=0.5 is the perceptual midpoint.
//
// OKLab gives straight-line interpolation — no unexpected hue swings
// through grey that you get with HSL or even OKLCH mixing.
// The result name and mood are cleared (the blend is a new thing).
func Mix(a, b Tone, t float64) Tone {
	t = clamp(t, 0, 1)

	// Convert both to OKLab
	La, aa, ba := oklchToOKLab(a.l, a.c, a.hue)
	Lb, ab, bb := oklchToOKLab(b.l, b.c, b.hue)

	// Lerp in OKLab
	L := La + t*(Lb-La)
	aL := aa + t*(ab-aa)
	bL := ba + t*(bb-ba)

	// Back to OKLCH
	l, c, h := oklabToOKLCH(L, aL, bL)

	// Blend alpha and energy
	alpha := a.a + t*(b.a-a.a)
	energy := a.energy + t*(b.energy-a.energy)

	maxC := maxChromaForLH(l, h)
	vibrancy := 0.0
	if maxC > 0 {
		vibrancy = clamp((c/maxC)*100, 0, 100)
	}

	return Tone{
		light:    l * 100,
		vibrancy: vibrancy,
		hue:      h,
		energy:   energy,
		l:        l,
		c:        c,
		a:        alpha,
	}
}

// Gradient produces n evenly spaced Tones between start and end in OKLab space.
// n must be >= 2. The first element is start, the last is end.
// No grey midpoint artifacts — OKLab interpolation guarantees perceptual smoothness.
func Gradient(start, end Tone, n int) ([]Tone, error) {
	if n < 2 {
		return nil, fmt.Errorf("wondertone: Gradient needs at least 2 steps, got %d", n)
	}
	result := make([]Tone, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		result[i] = Mix(start, end, t)
	}
	return result, nil
}

// Blend mixes multiple Tones with given weights in OKLab space.
// weights must be the same length as tones and sum > 0.
// Weights are normalized automatically — you don't need to sum to 1.
func Blend(tones []Tone, weights []float64) (Tone, error) {
	if len(tones) < 2 {
		return Tone{}, fmt.Errorf("wondertone: Blend needs at least 2 tones")
	}
	if len(tones) != len(weights) {
		return Tone{}, fmt.Errorf("wondertone: Blend: %d tones but %d weights", len(tones), len(weights))
	}

	// Normalize weights
	total := 0.0
	for _, w := range weights {
		if w < 0 {
			return Tone{}, fmt.Errorf("wondertone: Blend: weights must be non-negative")
		}
		total += w
	}
	if total == 0 {
		return Tone{}, fmt.Errorf("wondertone: Blend: weights sum to zero")
	}

	// Weighted sum in OKLab
	var sumL, sumA, sumB, sumAlpha, sumEnergy float64
	for i, tone := range tones {
		w := weights[i] / total
		L, a, b := oklchToOKLab(tone.l, tone.c, tone.hue)
		sumL += L * w
		sumA += a * w
		sumB += b * w
		sumAlpha += tone.a * w
		sumEnergy += tone.energy * w
	}

	l, c, h := oklabToOKLCH(sumL, sumA, sumB)
	maxC := maxChromaForLH(l, h)
	vibrancy := 0.0
	if maxC > 0 {
		vibrancy = clamp((c/maxC)*100, 0, 100)
	}

	return Tone{
		light:    l * 100,
		vibrancy: vibrancy,
		hue:      h,
		energy:   sumEnergy,
		l:        l,
		c:        c,
		a:        sumAlpha,
	}, nil
}

// Harmonize returns a set of Tones related by the given harmonic scheme.
// All harmony is computed by hue rotation in OKLCH — lightness and vibrancy
// are preserved from the base Tone.
//
// scheme: "complement", "triadic", "analogous", "split", "tetradic"
func Harmonize(base Tone, scheme string) ([]Tone, error) {
	switch scheme {
	case "complement":
		return []Tone{base, base.Rotate(180)}, nil
	case "triadic":
		return []Tone{base, base.Rotate(120), base.Rotate(240)}, nil
	case "analogous":
		return []Tone{base.Rotate(-30), base, base.Rotate(30)}, nil
	case "split":
		return []Tone{base, base.Rotate(150), base.Rotate(210)}, nil
	case "tetradic":
		return []Tone{base, base.Rotate(90), base.Rotate(180), base.Rotate(270)}, nil
	default:
		return nil, fmt.Errorf(
			"wondertone: unknown harmony scheme %q (want: complement, triadic, analogous, split, tetradic)",
			scheme,
		)
	}
}

// Shift nudges a Tone toward warmer or cooler territory.
// direction: "warmer" or "cooler". amount [0–1].
func Shift(t Tone, direction string, amount float64) (Tone, error) {
	amount = clamp(amount, 0, 1)
	var targetHue float64
	switch direction {
	case "warmer":
		targetHue = 30.0 // orange anchor
	case "cooler":
		targetHue = 210.0 // blue anchor
	default:
		return Tone{}, fmt.Errorf(
			"wondertone: Shift direction must be \"warmer\" or \"cooler\", got %q", direction,
		)
	}
	diff := targetHue - t.hue
	for diff > 180 {
		diff -= 360
	}
	for diff < -180 {
		diff += 360
	}
	return t.WithHue(normalizeHue(t.hue + diff*amount)), nil
}

// lerp linearly interpolates between a and b.
func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// ensure math import is used
var _ = math.Pi
