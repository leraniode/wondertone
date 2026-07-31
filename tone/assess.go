package tone

import "math"

// Luminance returns the WCAG 2.1 relative luminance [0–1].
func (t Tone) Luminance() float64 {
	r, g, b := t.toSRGB()
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
