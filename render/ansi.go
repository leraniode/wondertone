package render

import (
	"fmt"
	"math"

	"github.com/leraniode/wondertone/tone"
)

// ANSI escape codes
const (
	Reset     = "\x1b[0m"
	Bold      = "\x1b[1m"
	Dim       = "\x1b[2m"
	Italic    = "\x1b[3m"
	Underline = "\x1b[4m"
)

// FG returns the ANSI escape sequence to set the foreground to the given Tone.
// The sequence is adapted to the terminal Profile.
func FG(t tone.Tone, p Profile) string {
	return colorSequence(t, p, false)
}

// BG returns the ANSI escape sequence to set the background to the given Tone.
func BG(t tone.Tone, p Profile) string {
	return colorSequence(t, p, true)
}

// Colorize wraps text with FG color and a reset.
func Colorize(t tone.Tone, p Profile, text string) string {
	return FG(t, p) + text + Reset
}

// ColorizeOnBG wraps text with FG and BG colors and a reset.
func ColorizeOnBG(fg, bg tone.Tone, p Profile, text string) string {
	return FG(fg, p) + BG(bg, p) + text + Reset
}

// Swatch returns a colored block for palette preview.
// width is the number of full-block characters (█) to render.
func Swatch(t tone.Tone, p Profile, width int) string {
	if width <= 0 {
		width = 2
	}
	block := ""
	for i := 0; i < width; i++ {
		block += "█"
	}
	return BG(t, p) + "  " + Reset + " " + FG(t, p) + block + Reset
}

// colorSequence builds the raw ANSI escape sequence.
func colorSequence(t tone.Tone, p Profile, bg bool) string {
	switch p {
	case NoColor:
		return ""
	case TrueColor:
		return trueColorSequence(t, bg)
	case ANSI256:
		return ansi256Sequence(t, bg)
	case ANSI16:
		return ansi16Sequence(t, bg)
	default:
		return ""
	}
}

// --- TrueColor (24-bit) ---

func trueColorSequence(t tone.Tone, bg bool) string {
	r, g, b := t.RGB()
	if bg {
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// --- ANSI 256 ---

// ansi256Sequence finds the perceptually nearest color in the 256-color palette.
// Uses OKLab ΔE (Euclidean distance in OKLab space) — not RGB distance.
// This is dramatically better for perceptual accuracy.
func ansi256Sequence(t tone.Tone, bg bool) string {
	idx := nearestANSI256(t)
	if bg {
		return fmt.Sprintf("\x1b[48;5;%dm", idx)
	}
	return fmt.Sprintf("\x1b[38;5;%dm", idx)
}

// nearestANSI256 finds the index [0–255] of the perceptually closest ANSI256 color.
func nearestANSI256(t tone.Tone) int {
	tL, tA, tB := toneToOKLab(t)
	best, bestDist := 0, math.MaxFloat64
	for i, rgb := range ansi256Palette {
		cL, cA, cB := rgbToOKLab(rgb[0], rgb[1], rgb[2])
		d := oklabDist(tL, tA, tB, cL, cA, cB)
		if d < bestDist {
			best, bestDist = i, d
		}
	}
	return best
}

// --- ANSI 16 ---

func ansi16Sequence(t tone.Tone, bg bool) string {
	idx := nearestANSI16(t)
	if bg {
		if idx >= 8 {
			// Bright colors: use standard bg code + bright offset
			return fmt.Sprintf("\x1b[%dm", 100+(idx-8))
		}
		return fmt.Sprintf("\x1b[%dm", 40+idx)
	}
	if idx >= 8 {
		return fmt.Sprintf("\x1b[%dm", 90+(idx-8))
	}
	return fmt.Sprintf("\x1b[%dm", 30+idx)
}

func nearestANSI16(t tone.Tone) int {
	tL, tA, tB := toneToOKLab(t)
	best, bestDist := 0, math.MaxFloat64
	for i, rgb := range ansi16Colors {
		cL, cA, cB := rgbToOKLab(rgb[0], rgb[1], rgb[2])
		d := oklabDist(tL, tA, tB, cL, cA, cB)
		if d < bestDist {
			best, bestDist = i, d
		}
	}
	return best
}

// --- OKLab distance helpers ---

func toneToOKLab(t tone.Tone) (L, a, b float64) {
	r, g, bl := t.RGBFloat()
	return srgbToOKLab(r, g, bl)
}

func rgbToOKLab(r, g, b uint8) (L, a, bl float64) {
	// sRGB [0-255] → linear → OKLab
	rf := srgbToLinear(float64(r) / 255.0)
	gf := srgbToLinear(float64(g) / 255.0)
	bf := srgbToLinear(float64(b) / 255.0)
	return srgbToOKLab(rf, gf, bf)
}

// srgbToOKLab converts gamma-encoded sRGB [0-1] to OKLab.
func srgbToOKLab(r, g, b float64) (L, a, bl float64) {
	lr := srgbToLinear(r)
	lg := srgbToLinear(g)
	lb := srgbToLinear(b)

	// Linear RGB → LMS
	l := 0.4122214708*lr + 0.5363325363*lg + 0.0514459929*lb
	m := 0.2119034982*lr + 0.6806995451*lg + 0.1073969566*lb
	s := 0.0883024619*lr + 0.2817188376*lg + 0.6299787005*lb

	l_ := math.Cbrt(l)
	m_ := math.Cbrt(m)
	s_ := math.Cbrt(s)

	L = 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	a = 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	bl = 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_
	return
}

func srgbToLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func oklabDist(L1, a1, b1, L2, a2, b2 float64) float64 {
	dL := L1 - L2
	da := a1 - a2
	db := b1 - b2
	return dL*dL + da*da + db*db // no sqrt needed for comparison
}

// --- ANSI 16 color table ---
// Standard terminal 16 colors in sRGB [0-255].
// Order: black, red, green, yellow, blue, magenta, cyan, white,
//
//	bright-black, bright-red, bright-green, bright-yellow,
//	bright-blue, bright-magenta, bright-cyan, bright-white
var ansi16Colors = [16][3]uint8{
	{0, 0, 0},       // 0  black
	{170, 0, 0},     // 1  red
	{0, 170, 0},     // 2  green
	{170, 85, 0},    // 3  yellow (dark)
	{0, 0, 170},     // 4  blue
	{170, 0, 170},   // 5  magenta
	{0, 170, 170},   // 6  cyan
	{170, 170, 170}, // 7  white
	{85, 85, 85},    // 8  bright black (dark grey)
	{255, 85, 85},   // 9  bright red
	{85, 255, 85},   // 10 bright green
	{255, 255, 85},  // 11 bright yellow
	{85, 85, 255},   // 12 bright blue
	{255, 85, 255},  // 13 bright magenta
	{85, 255, 255},  // 14 bright cyan
	{255, 255, 255}, // 15 bright white
}

// ansi256Palette is the full 256-color ANSI palette in sRGB [0-255].
// Generated programmatically: 16 system colors + 216 color cube + 24 greys.
var ansi256Palette = func() [256][3]uint8 {
	var p [256][3]uint8
	// 0-15: system colors (same as ansi16Colors)
	for i := 0; i < 16; i++ {
		p[i] = ansi16Colors[i]
	}
	// 16-231: 6×6×6 color cube
	for i := 0; i < 216; i++ {
		r := i / 36
		g := (i / 6) % 6
		b := i % 6
		toVal := func(v int) uint8 {
			if v == 0 {
				return 0
			}
			return uint8(55 + v*40)
		}
		p[16+i] = [3]uint8{toVal(r), toVal(g), toVal(b)}
	}
	// 232-255: 24 greyscale steps
	for i := 0; i < 24; i++ {
		v := uint8(8 + i*10)
		p[232+i] = [3]uint8{v, v, v}
	}
	return p
}()
