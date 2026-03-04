package render_test

import (
	"fmt"
	"strings"
	"testing"

	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/internal/testutil"
	"github.com/leraniode/wondertone/render"
)

var unix = tone.New(
	tone.Light(68), tone.Vibrancy(72), tone.Hue(142),
	tone.Energy(0.95), tone.Named("Unix"),
)

// --- TrueColor ---

func TestFGTrueColor(t *testing.T) {
	seq := render.FG(unix, render.TrueColor)
	testutil.True(t, strings.HasPrefix(seq, "\x1b[38;2;"), "TrueColor FG should start with \\x1b[38;2;")
	testutil.True(t, strings.HasSuffix(seq, "m"), "TrueColor FG should end with m")
}

func TestBGTrueColor(t *testing.T) {
	seq := render.BG(unix, render.TrueColor)
	testutil.True(t, strings.HasPrefix(seq, "\x1b[48;2;"), "TrueColor BG should start with \\x1b[48;2;")
}

func TestTrueColorContainsRGB(t *testing.T) {
	r, g, b := unix.RGB()
	seq := render.FG(unix, render.TrueColor)
	expected := strings.Join([]string{
		"\x1b[38;2",
		itoa(r),
		itoa(g),
		itoa(b) + "m",
	}, ";")
	// Remove last join artifact
	expected = "\x1b[38;2;" + strings.TrimPrefix(expected, "\x1b[38;2;")
	testutil.Equal(t, expected, seq)
}

// --- ANSI256 ---

func TestFGANSI256(t *testing.T) {
	seq := render.FG(unix, render.ANSI256)
	testutil.True(t, strings.HasPrefix(seq, "\x1b[38;5;"), "ANSI256 FG prefix wrong")
	testutil.True(t, strings.HasSuffix(seq, "m"), "ANSI256 FG suffix wrong")
}

func TestBGANSI256(t *testing.T) {
	seq := render.BG(unix, render.ANSI256)
	testutil.True(t, strings.HasPrefix(seq, "\x1b[48;5;"), "ANSI256 BG prefix wrong")
}

func TestANSI256IndexInRange(t *testing.T) {
	cases := []tone.Tone{
		tone.New(tone.Light(5), tone.Vibrancy(0), tone.Hue(0)),   // near black
		tone.New(tone.Light(95), tone.Vibrancy(0), tone.Hue(0)),  // near white
		tone.New(tone.Light(50), tone.Vibrancy(90), tone.Hue(0)), // vivid red
		unix,
	}
	for _, tc := range cases {
		seq := render.FG(tc, render.ANSI256)
		// Extract the number
		var idx int
		_, err := parseSeq(seq, &idx)
		testutil.NoError(t, err)
		testutil.True(t, idx >= 0 && idx <= 255, "ANSI256 index should be in [0,255], got %d", idx)
	}
}

// --- ANSI16 ---

func TestFGANSI16(t *testing.T) {
	seq := render.FG(unix, render.ANSI16)
	testutil.True(t, len(seq) > 0, "ANSI16 FG should not be empty")
	testutil.True(t, strings.HasPrefix(seq, "\x1b["), "ANSI16 FG should be an escape sequence")
}

// --- NoColor ---

func TestNoColor(t *testing.T) {
	testutil.Equal(t, "", render.FG(unix, render.NoColor))
	testutil.Equal(t, "", render.BG(unix, render.NoColor))
}

// --- Colorize ---

func TestColorize(t *testing.T) {
	result := render.Colorize(unix, render.TrueColor, "hello")
	testutil.True(t, strings.HasPrefix(result, "\x1b[38;2;"))
	testutil.True(t, strings.Contains(result, "hello"))
	testutil.True(t, strings.HasSuffix(result, render.Reset))
}

// --- Perceptual accuracy: black → black, white → white ---

func TestBlackMapsToBlack(t *testing.T) {
	black := tone.New(tone.Light(2), tone.Vibrancy(0), tone.Hue(0))
	seq := render.FG(black, render.ANSI256)
	var idx int
	parseSeq(seq, &idx)
	// ANSI 256 index 0 is black, 232 is also very dark grey
	testutil.True(t, idx <= 8 || idx >= 232, "near-black should map to a dark ANSI256 color, got %d", idx)
}

func TestWhiteMapsToWhite(t *testing.T) {
	white := tone.New(tone.Light(98), tone.Vibrancy(0), tone.Hue(0))
	seq := render.FG(white, render.ANSI256)
	var idx int
	parseSeq(seq, &idx)
	// Index 15 = bright white, 255 = lightest grey in 256 palette
	testutil.True(t, idx >= 231 || idx == 15 || idx == 7, "near-white should map to a light ANSI256 color, got %d", idx)
}

// --- LipglossColor ---

func TestLipglossColorTrueColor(t *testing.T) {
	s := render.LipglossColor(unix, render.TrueColor)
	testutil.True(t, strings.HasPrefix(s, "#"), "TrueColor lipgloss color should be hex")
	testutil.Equal(t, 7, len(s))
}

func TestLipglossColorANSI256(t *testing.T) {
	s := render.LipglossColor(unix, render.ANSI256)
	testutil.True(t, len(s) > 0, "ANSI256 lipgloss color should not be empty")
}

func TestLipglossColorNoColor(t *testing.T) {
	testutil.Equal(t, "", render.LipglossColor(unix, render.NoColor))
}

// --- Profile String ---

func TestProfileString(t *testing.T) {
	testutil.Equal(t, "TrueColor", render.TrueColor.String())
	testutil.Equal(t, "ANSI256", render.ANSI256.String())
	testutil.Equal(t, "ANSI16", render.ANSI16.String())
	testutil.Equal(t, "NoColor", render.NoColor.String())
}

// --- helpers ---

func itoa(v uint8) string {
	return fmt.Sprintf("%d", v)
}

// parseSeq extracts the integer from sequences like "\x1b[38;5;42m".
func parseSeq(seq string, out *int) (string, error) {
	// strip prefix and suffix
	seq = strings.TrimPrefix(seq, "\x1b[")
	seq = strings.TrimSuffix(seq, "m")
	parts := strings.Split(seq, ";")
	if len(parts) == 0 {
		return seq, fmt.Errorf("empty sequence")
	}
	last := parts[len(parts)-1]
	_, err := fmt.Sscanf(last, "%d", out)
	return seq, err
}
