package core

// OKLCH ↔ OKLab ↔ Linear RGB ↔ sRGB pipeline.
//
// All math is public domain — based on Björn Ottosson's OKLab specification
// (https://bottosson.github.io/posts/oklab/) and the standard sRGB/linear
// conversion formulas. No third-party dependency required.
//
// Pipeline (forward):  OKLCH → OKLab → XYZ(D65) → Linear RGB → sRGB
// Pipeline (reverse):  sRGB → Linear RGB → XYZ(D65) → OKLab → OKLCH

import "math"

// --- OKLCH ↔ OKLab ---

// oklchToOKLab converts OKLCH to OKLab cartesian coordinates.
// H is in degrees [0–360).
func oklchToOKLab(l, c, h float64) (L, a, b float64) {
	hRad := h * math.Pi / 180.0
	return l, c * math.Cos(hRad), c * math.Sin(hRad)
}

// oklabToOKLCH converts OKLab cartesian to OKLCH polar.
// H returned in degrees [0–360).
func oklabToOKLCH(L, a, b float64) (l, c, h float64) {
	c = math.Sqrt(a*a + b*b)
	h = math.Atan2(b, a) * 180.0 / math.Pi
	if h < 0 {
		h += 360
	}
	return L, c, h
}

// --- OKLab ↔ Linear RGB ---
// Matrix coefficients from Björn Ottosson's OKLab specification.

// oklabToLinearRGB converts OKLab to linear sRGB via the OKLab transform matrices.
func oklabToLinearRGB(L, a, b float64) (r, g, bl float64) {
	// Step 1: OKLab → LMS (cube-root space)
	l_ := L + 0.3963377774*a + 0.2158037573*b
	m_ := L - 0.1055613458*a - 0.0638541728*b
	s_ := L - 0.0894841775*a - 1.2914855480*b

	// Step 2: LMS cube-root → LMS linear
	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	// Step 3: LMS → Linear RGB
	r = +4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	g = -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	bl = -0.0041960863*l - 0.7034186147*m + 1.7076147010*s
	return r, g, bl
}

// linearRGBToOKLab converts linear sRGB to OKLab.
func linearRGBToOKLab(r, g, b float64) (L, a, bl float64) {
	// Step 1: Linear RGB → LMS
	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b

	// Step 2: LMS → LMS cube-root space
	l_ := math.Cbrt(l)
	m_ := math.Cbrt(m)
	s_ := math.Cbrt(s)

	// Step 3: LMS cube-root → OKLab
	L = 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	a = 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	bl = 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_
	return L, a, bl
}

// --- Top-level pipeline (OKLCH ↔ Linear RGB) ---

// oklchToLinearRGB is the main forward pipeline: OKLCH → OKLab → Linear RGB.
func oklchToLinearRGB(l, c, h float64) (r, g, b float64) {
	L, a, bl := oklchToOKLab(l, c, h)
	return oklabToLinearRGB(L, a, bl)
}

// linearRGBToOKLCH is the main reverse pipeline: Linear RGB → OKLab → OKLCH.
func linearRGBToOKLCH(r, g, b float64) (l, c, h float64) {
	L, a, bl := linearRGBToOKLab(r, g, b)
	return oklabToOKLCH(L, a, bl)
}

// --- sRGB gamma encoding / decoding ---

// linearToSRGB applies sRGB gamma encoding (linear → display-ready [0–1]).
// Clamps input to [0,1] — gamut mapping guarantees this but float drift can
// produce tiny excursions that the clamp absorbs cleanly.
func linearToSRGB(v float64) float64 {
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	if v <= 0.0031308 {
		return v * 12.92
	}
	return 1.055*math.Pow(v, 1.0/2.4) - 0.055
}

// srgbToLinear applies sRGB gamma decoding (display-encoded → linear).
func srgbToLinear(v float64) float64 {
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// --- Gamut checking ---

// inSRGB reports whether linear RGB values are within the sRGB gamut.
// A small epsilon absorbs floating-point rounding from the color math.
func inSRGB(r, g, b float64) bool {
	const eps = 1e-6
	return r >= -eps && r <= 1+eps &&
		g >= -eps && g <= 1+eps &&
		b >= -eps && b <= 1+eps
}
