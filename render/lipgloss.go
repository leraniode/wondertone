package render

import (
	"fmt"

	"github.com/leraniode/wondertone/tone"
)

// LipglossColor converts a Tone to a lipgloss-compatible color value string.
//
// Usage with lipgloss (import lipgloss separately):
//
//	import "github.com/charmbracelet/lipgloss"
//
//	style := lipgloss.NewStyle().
//	    Foreground(lipgloss.Color(render.LipglossColor(myTone, profile)))
//
// Returns a string that lipgloss.Color() accepts:
//   - TrueColor:  "#rrggbb"
//   - ANSI256:    "42" (the 256-color index as a string)
//   - ANSI16:     "2"  (the 16-color index as a string)
//   - NoColor:    ""   (lipgloss will render unstyled)
func LipglossColor(t tone.Tone, p Profile) string {
	switch p {
	case TrueColor:
		return t.Hex()
	case ANSI256:
		return fmt.Sprintf("%d", nearestANSI256(t))
	case ANSI16:
		return fmt.Sprintf("%d", nearestANSI16(t))
	default:
		return ""
	}
}

// LipglossColorHex always returns the hex value regardless of profile.
// Use this when lipgloss itself handles profile detection (e.g. via ColorProfile adapter).
func LipglossColorHex(t tone.Tone) string {
	return t.Hex()
}
