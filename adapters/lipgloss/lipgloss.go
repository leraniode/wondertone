// Package lipgloss provides a first-class wondertone adapter for
// github.com/charmbracelet/lipgloss.
//
// This is the adapter for using wondertone with lipgloss — import it only if you are using lipgloss.
// The wondertone core and render packages have zero dependency on lipgloss.
//
// Usage:
//
//	import (
//	    tone  "github.com/leraniode/wondertone/core"
//	    wtlip "github.com/leraniode/wondertone/adapters/lipgloss"
//	)
//
//	wtlip.FG(colour.Unix).Bold(true).Render("hello")
//
//	wtlip.Style(colour.Unix).
//	    Background(colour.Void).
//	    Padding(0, 1).
//	    Render("hello")
//
//	wtlip.PaletteStyles(builtin.Midnight()) // map[name]lipgloss.Style
package lipgloss

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/render"
)

func init() {
	// Sync lipgloss's default renderer with whatever profile we detected.
	// This ensures SetProfile() takes effect even when called before any
	// style is created — which is the normal usage in tests and apps.
	lipgloss.SetColorProfile(toTermenvProfile(render.Detect()))
}

// profile is our wondertone render profile — drives Color() output.
var profile = render.Detect()

// toTermenvProfile maps a wondertone render.Profile to termenv.Profile.
func toTermenvProfile(p render.Profile) termenv.Profile {
	switch p {
	case render.TrueColor:
		return termenv.TrueColor
	case render.ANSI256:
		return termenv.ANSI256
	case render.ANSI16:
		return termenv.ANSI
	default:
		return termenv.Ascii
	}
}

// SetProfile overrides terminal profile detection for both wondertone
// colour conversion and lipgloss rendering.
//
// Call this in tests to guarantee coloured output regardless of TTY:
//
//	func init() { wtlip.SetProfile(render.TrueColor) }
func SetProfile(p render.Profile) {
	profile = p
	// Sync lipgloss's default renderer — this is what NewStyle() uses.
	lipgloss.SetColorProfile(toTermenvProfile(p))
}

// Color converts a wondertone Tone to a lipgloss.Color.
func Color(t tone.Tone) lipgloss.Color {
	return lipgloss.Color(render.LipglossColor(t, profile))
}

// ColorHex always returns the hex lipgloss.Color regardless of terminal profile.
func ColorHex(t tone.Tone) lipgloss.Color {
	return lipgloss.Color(t.Hex())
}

// FG returns a lipgloss.Style with the tone as foreground colour.
func FG(t tone.Tone) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(Color(t))
}

// BG returns a lipgloss.Style with the tone as background colour.
func BG(t tone.Tone) lipgloss.Style {
	return lipgloss.NewStyle().Background(Color(t))
}

// StyleBuilder wraps a lipgloss.Style with wondertone-aware setters.
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

// Render applies the style to text.
func (b *StyleBuilder) Render(text string) string {
	return b.s.Render(text)
}

// Lipgloss returns the underlying lipgloss.Style for further customisation.
func (b *StyleBuilder) Lipgloss() lipgloss.Style {
	return b.s
}

// PaletteStyles returns a map of tone name → foreground lipgloss.Style
// for every tone in the palette.
//
//	styles := wtlip.PaletteStyles(builtin.Midnight())
//	styles["Midnight Accent"].Bold(true).Render("accent text")
func PaletteStyles(p *palette.Palette) map[string]lipgloss.Style {
	out := make(map[string]lipgloss.Style, p.Len())
	for _, t := range p.All() {
		out[t.Name()] = FG(t)
	}
	return out
}

// AdaptiveStyle returns a foreground style chosen based on whether
// the background tone is light or dark — so text always contrasts.
//
//	style := wtlip.AdaptiveStyle(colour.Ink, colour.Paper, bgTone)
func AdaptiveStyle(onLight, onDark tone.Tone, bg tone.Tone) lipgloss.Style {
	if bg.IsLight() {
		return FG(onLight)
	}
	return FG(onDark)
}
