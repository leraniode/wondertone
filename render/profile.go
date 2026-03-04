// Package render handles terminal color output for wondertone Tones.
//
// It detects the terminal's color capability and outputs the best
// possible ANSI escape sequence — TrueColor, 256-color, or 16-color.
// Downsampling is perceptual (OKLab ΔE nearest-neighbor), not RGB distance.
//
// Import as "render":
//
//	import render "github.com/leraniode/wondertone/render"
//
//	profile := render.Detect()
//	fmt.Print(render.FG(myTone, profile), "Hello, wondertone", render.Reset)
package render

import (
	"os"
	"strings"
)

// Profile represents the terminal's color capability.
type Profile int

const (
	// TrueColor supports 24-bit RGB — the full wondertone experience.
	TrueColor Profile = iota
	// ANSI256 supports 256 colors — perceptual nearest-neighbor downsampling.
	ANSI256
	// ANSI16 supports the 16 standard ANSI colors — best-effort perceptual map.
	ANSI16
	// NoColor disables all color output.
	NoColor
)

func (p Profile) String() string {
	switch p {
	case TrueColor:
		return "TrueColor"
	case ANSI256:
		return "ANSI256"
	case ANSI16:
		return "ANSI16"
	default:
		return "NoColor"
	}
}

// Detect reads environment variables to determine the terminal's color profile.
// Checks (in order): NO_COLOR, COLORTERM, TERM, TERM_PROGRAM, CI.
//
// For guaranteed detection in production, pair with charmbracelet/colorprofile.
func Detect() Profile {
	// NO_COLOR — https://no-color.org
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return NoColor
	}

	// Explicit no-color flag
	if os.Getenv("TERM") == "dumb" {
		return NoColor
	}

	// TrueColor detection
	colorterm := strings.ToLower(os.Getenv("COLORTERM"))
	if colorterm == "truecolor" || colorterm == "24bit" {
		return TrueColor
	}

	termProg := os.Getenv("TERM_PROGRAM")
	switch termProg {
	case "iTerm.app", "Hyper", "vscode", "WezTerm", "ghostty":
		return TrueColor
	}

	term := os.Getenv("TERM")
	if strings.HasSuffix(term, "-256color") || strings.HasSuffix(term, "256color") {
		return ANSI256
	}
	if strings.Contains(term, "color") || strings.HasPrefix(term, "xterm") {
		return TrueColor
	}

	// CI environments often support 256 colors
	if os.Getenv("CI") != "" {
		return ANSI256
	}

	// Conservative fallback
	return ANSI16
}

// Force returns a Profile that always returns the given level, ignoring env.
// Useful for testing or when you know the terminal capability.
func Force(p Profile) Profile { return p }
