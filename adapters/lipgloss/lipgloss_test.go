package lipgloss_test

import (
	"strings"
	"testing"

	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/colour"
	wtlip "github.com/leraniode/wondertone/adapters/lipgloss"
	"github.com/leraniode/wondertone/internal/testutil"
	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/render"
)

func init() {
	// Force TrueColor so lipgloss emits ANSI codes in the test environment.
	// SetProfile syncs both wondertone and lipgloss's default renderer.
	wtlip.SetProfile(render.TrueColor)
}

func TestColorReturnsHexInTrueColor(t *testing.T) {
	c := wtlip.Color(colour.Unix)
	testutil.Equal(t, colour.Unix.Hex(), string(c))
}

func TestFGProducesStyledText(t *testing.T) {
	style := wtlip.FG(colour.Unix)
	rendered := style.Render("hello")
	testutil.True(t, strings.Contains(rendered, "hello"),
		"rendered text should contain original string")
	testutil.True(t, len(rendered) > len("hello"),
		"styled text should be longer than raw text (ANSI codes expected)")
}

func TestBGProducesStyledText(t *testing.T) {
	style := wtlip.BG(colour.Void)
	rendered := style.Render("hello")
	testutil.True(t, strings.Contains(rendered, "hello"))
	testutil.True(t, len(rendered) > len("hello"))
}

func TestStyleBuilder(t *testing.T) {
	rendered := wtlip.Style(colour.Unix).
		Background(colour.Void).
		Bold(true).
		Padding(0, 1).
		Render("test")
	testutil.True(t, strings.Contains(rendered, "test"))
}

func TestStyleBuilderLipgloss(t *testing.T) {
	ls := wtlip.Style(colour.Unix).Lipgloss()
	result := ls.Render("ok")
	testutil.True(t, strings.Contains(result, "ok"))
}

func TestPaletteStyles(t *testing.T) {
	p, _ := palette.New("test").
		Add(tone.New(tone.Light(50), tone.Hue(30), tone.Named("Warm"))).
		Add(tone.New(tone.Light(50), tone.Hue(200), tone.Named("Cool"))).
		Build()

	styles := wtlip.PaletteStyles(p)
	testutil.Equal(t, 2, len(styles))
	testutil.True(t, len(styles["Warm"].Render("x")) > 0)
	testutil.True(t, len(styles["Cool"].Render("x")) > 0)
}

func TestAdaptiveStylePicksCorrectVariant(t *testing.T) {
	lightBg := tone.New(tone.Light(90), tone.Vibrancy(0), tone.Hue(0))
	darkBg  := tone.New(tone.Light(10), tone.Vibrancy(0), tone.Hue(0))

	onLight := colour.Ink
	onDark  := colour.Paper

	lightStyle := wtlip.AdaptiveStyle(onLight, onDark, lightBg)
	darkStyle  := wtlip.AdaptiveStyle(onLight, onDark, darkBg)

	// Both should render the text
	testutil.True(t, strings.Contains(lightStyle.Render("x"), "x"))
	testutil.True(t, strings.Contains(darkStyle.Render("x"), "x"))

	// The underlying colours should differ — Ink hex vs Paper hex
	// We verify this through Color() directly, which is reliable and
	// not subject to ANSI downsampling edge cases.
	inkColor   := string(wtlip.Color(colour.Ink))
	paperColor := string(wtlip.Color(colour.Paper))
	testutil.True(t, inkColor != paperColor,
		"Ink and Paper should produce different lipgloss.Color values")

	// In TrueColor mode, rendered output must also differ
	testutil.True(t,
		lightStyle.Render("x") != darkStyle.Render("x"),
		"adaptive style should differ for light vs dark background")
}

func TestColorHexAlwaysHex(t *testing.T) {
	wtlip.SetProfile(render.ANSI256)
	c := wtlip.ColorHex(colour.Unix) // ColorHex ignores profile
	testutil.True(t, strings.HasPrefix(string(c), "#"))
	wtlip.SetProfile(render.TrueColor) // restore for other tests
}