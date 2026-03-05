// Package lipgloss provides a first-class wondertone adapter for
// github.com/charmbracelet/lipgloss.
//
// This is the adapter for using wondertone with lipgloss — import it only if you are using lipgloss.
// The wondertone core and render packages have zero dependency on lipgloss.
//
// Usage:
//
//	import (
//	    tone    "github.com/leraniode/wondertone/core"
//	    "github.com/leraniode/wondertone/palette/builtin"
//	    wtlip   "github.com/leraniode/wondertone/adapters/lipgloss"
//	)
//
//	// Simplest usage — foreground style from a tone
//	style := wtlip.FG(colour.Unix)
//	fmt.Println(style.Render("hello, wondertone"))
//
//	// Full style builder
//	style := wtlip.Style(colour.Unix).
//	    Background(colour.Void).
//	    Bold(true).
//	    Padding(0, 1)
//	fmt.Println(style.Render("hello"))
//
//	// Palette → full style map
//	styles := wtlip.PaletteStyles(builtin.Midnight())
//	fmt.Println(styles["Midnight Accent"].Render("accent text"))
package lipgloss

import (
	"github.com/charmbracelet/lipgloss"
	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/render"
)

// profile is the detected terminal profile — set once at package init.
var profile = render.Detect()

// SetProfile overrides terminal profile detection.
// Call this if you manage profile detection yourself (e.g. via charmbracelet/colorprofile).
func SetProfile(p render.Profile) {
	profile = p
}

// Color converts a wondertone Tone to a lipgloss.Color.
// Respects the current terminal profile — TrueColor returns hex,
// ANSI256 returns the index string, ANSI16 returns the base color index,
// NoColor returns an empty string (lipgloss renders unstyled).
func Color(t tone.Tone) lipgloss.Color {
	return lipgloss.Color(render.LipglossColor(t, profile))
}

// ColorHex always returns the hex lipgloss.Color regardless of terminal profile.
// Use this when lipgloss itself handles profile/color downsampling.
func ColorHex(t tone.Tone) lipgloss.Color {
	return lipgloss.Color(t.Hex())
}

// FG returns a lipgloss.Style with the tone set as foreground color.
func FG(t tone.Tone) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(Color(t))
}

// BG returns a lipgloss.Style with the tone set as background color.
func BG(t tone.Tone) lipgloss.Style {
	return lipgloss.NewStyle().Background(Color(t))
}

// StyleBuilder wraps a lipgloss.Style with wondertone-aware methods.
// Returned by Style() — chain methods, call .Lipgloss() to get the final style.
type StyleBuilder struct {
	s lipgloss.Style
}

// Style starts a new StyleBuilder with the tone as foreground.
func Style(t tone.Tone) *StyleBuilder {
	return &StyleBuilder{s: FG(t)}
}

// Background sets the background tone.
func (b *StyleBuilder) Background(t tone.Tone) *StyleBuilder {
	b.s = b.s.Background(Color(t))
	return b
}

// Foreground sets the foreground tone.
func (b *StyleBuilder) Foreground(t tone.Tone) *StyleBuilder {
	b.s = b.s.Foreground(Color(t))
	return b
}

// Bold sets bold text.
func (b *StyleBuilder) Bold(v bool) *StyleBuilder {
	b.s = b.s.Bold(v)
	return b
}

// Italic sets italic text.
func (b *StyleBuilder) Italic(v bool) *StyleBuilder {
	b.s = b.s.Italic(v)
	return b
}

// Underline sets underline.
func (b *StyleBuilder) Underline(v bool) *StyleBuilder {
	b.s = b.s.Underline(v)
	return b
}

// Padding sets padding (top/bottom, left/right).
func (b *StyleBuilder) Padding(tb, lr int) *StyleBuilder {
	b.s = b.s.Padding(tb, lr)
	return b
}

// Margin sets margin (top/bottom, left/right).
func (b *StyleBuilder) Margin(tb, lr int) *StyleBuilder {
	b.s = b.s.Margin(tb, lr)
	return b
}

// Width sets the style width.
func (b *StyleBuilder) Width(w int) *StyleBuilder {
	b.s = b.s.Width(w)
	return b
}

// Render applies the style to text — shortcut so you don't need to call Lipgloss().
func (b *StyleBuilder) Render(text string) string {
	return b.s.Render(text)
}

// Lipgloss returns the underlying lipgloss.Style for further customisation.
func (b *StyleBuilder) Lipgloss() lipgloss.Style {
	return b.s
}

// PaletteStyles returns a map of tone name → foreground lipgloss.Style
// for every tone in the palette. Useful for building a style kit from
// a wondertone palette.
//
//	styles := wtlip.PaletteStyles(builtin.Midnight())
//	banner := styles["Midnight Accent"].Bold(true).Padding(0, 1)
func PaletteStyles(p *palette.Palette) map[string]lipgloss.Style {
	out := make(map[string]lipgloss.Style, p.Len())
	for _, t := range p.All() {
		out[t.Name()] = FG(t)
	}
	return out
}

// AdaptiveStyle returns different styles for light vs dark backgrounds.
// detectBg should return the current background tone so wondertone can
// pick the right variant.
//
//	style := wtlip.AdaptiveStyle(
//	    colour.Ink,       // use on light backgrounds
//	    colour.Paper,     // use on dark backgrounds
//	    myBgTone,
//	)
func AdaptiveStyle(onLight, onDark tone.Tone, bg tone.Tone) lipgloss.Style {
	if bg.IsLight() {
		return FG(onLight)
	}
	return FG(onDark)
}
