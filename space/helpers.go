package space

import "math"

// --- Helpers ---

// normalizedSaturation returns C / C_max — the normalised saturation [0,1]
// used in the mood formulas. Same as Vibrancy/100 but computed from raw C.
func NormalizedSaturation(c, l, h float64) float64 {
	cMax := MaxChromaForLH(l, h)
	if cMax <= 0 {
		return 0
	}
	return Clamp(c/cMax, 0, 1)
}

// normalizeHue wraps hue to [0, 360).
func NormalizeHue(h float64) float64 {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	return h
}

// clamp restricts v to [min, max].
func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
