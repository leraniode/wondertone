package core

// Gamut safety — keeps every Tone inside the sRGB display gamut.
//
// Two operations:
//   - toGamutSafe: reduces chroma iteratively until in-gamut (your design)
//   - maxChromaForLH: binary search for the chroma ceiling at a given L and H
//
// Hue NEVER changes during gamut operations. This is non-negotiable —
// hue drift is the silent failure mode of naive RGB clamping.

// toGamutSafe reduces chroma iteratively until the OKLCH color fits in sRGB.
// L and H are preserved exactly. Each iteration reduces C by 2%.
//
// In practice this converges in <20 iterations for any real color.
// Worst case (a vivid hue at an extreme lightness) takes ~115 iterations.
func toGamutSafe(l, c, h float64) (float64, float64, float64) {
	r, g, b := oklchToLinearRGB(l, c, h)
	if inSRGB(r, g, b) {
		return l, c, h
	}
	for c > 0.001 {
		c *= 0.98
		r, g, b = oklchToLinearRGB(l, c, h)
		if inSRGB(r, g, b) {
			return l, c, h
		}
	}
	// Color is so achromatic that even C≈0 needs a tiny nudge — grey it out.
	return l, 0, h
}

// maxChromaForLH returns the maximum chroma that produces an in-gamut sRGB color
// at the given lightness and hue. Used to compute the Vibrancy percentage.
//
// Binary search — 24 iterations gives precision of ~0.00000006.
// This is called once at Tone construction time, not at render time.
func maxChromaForLH(l, h float64) float64 {
	const (
		iterations = 24
		ceiling    = 0.5 // practical upper bound; sRGB never exceeds ~0.37
	)
	lo, hi := 0.0, ceiling
	for i := 0; i < iterations; i++ {
		mid := (lo + hi) / 2
		r, g, b := oklchToLinearRGB(l, mid, h)
		if inSRGB(r, g, b) {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
}
