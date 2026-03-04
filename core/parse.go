package core

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// parseHex parses a CSS hex string into linear sRGB + alpha [0–1].
// Accepts #rgb, #rrggbb, #rrggbbaa.
func parseHex(s string) (r, g, b, a float64, err error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "#")

	switch len(s) {
	case 3:
		// Expand #rgb → #rrggbb
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	case 6:
		// canonical
	case 8:
		// with alpha
	default:
		return 0, 0, 0, 0, fmt.Errorf("wondertone: invalid hex color %q", "#"+s)
	}

	val, parseErr := strconv.ParseUint(s[:6], 16, 32)
	if parseErr != nil {
		return 0, 0, 0, 0, fmt.Errorf("wondertone: invalid hex color %q: %w", "#"+s, parseErr)
	}

	// Decode to sRGB [0–1] then to linear for internal storage
	sr := float64((val>>16)&0xff) / 255.0
	sg := float64((val>>8)&0xff) / 255.0
	sb := float64(val&0xff) / 255.0
	a = 1.0

	if len(s) == 8 {
		alpha, alphaErr := strconv.ParseUint(s[6:8], 16, 8)
		if alphaErr != nil {
			return 0, 0, 0, 0, fmt.Errorf("wondertone: invalid hex alpha in %q: %w", "#"+s, alphaErr)
		}
		a = float64(alpha) / 255.0
	}

	// sRGB → linear for the internal pipeline
	return srgbToLinear(sr), srgbToLinear(sg), srgbToLinear(sb), a, nil
}

// parseOKLCHString parses "L C H" or "L C H / A" strings.
// These are the values stored in .wtone files.
func parseOKLCHString(s string) (l, c, h, a float64, err error) {
	s = strings.TrimSpace(s)
	// Strip optional oklch() wrapper
	s = strings.TrimPrefix(s, "oklch(")
	s = strings.TrimSuffix(s, ")")

	a = 1.0
	if idx := strings.Index(s, "/"); idx >= 0 {
		alphaPart := strings.TrimSpace(s[idx+1:])
		s = s[:idx]
		a, err = strconv.ParseFloat(strings.TrimSpace(alphaPart), 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("wondertone: invalid alpha in oklch string %q", s)
		}
	}

	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) != 3 {
		return 0, 0, 0, 0, fmt.Errorf("wondertone: expected 3 values in oklch, got %d in %q", len(parts), s)
	}

	vals := [3]*float64{&l, &c, &h}
	for i, p := range parts {
		*vals[i], err = strconv.ParseFloat(p, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("wondertone: invalid oklch value %q", p)
		}
	}
	return l, c, h, a, nil
}

// formatOKLCHString formats OKLCH values into the canonical .wtone string.
// 6 decimal places: full precision, human readable.
func formatOKLCHString(l, c, h, a float64) string {
	round := func(v float64) float64 { return math.Round(v*1e6) / 1e6 }
	l, c, h = round(l), round(c), round(h)
	if a >= 1.0 {
		return fmt.Sprintf("%g %g %g", l, c, h)
	}
	return fmt.Sprintf("%g %g %g / %g", l, c, h, round(a))
}
