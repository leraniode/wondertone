// Package wtone handles reading and writing .wtone palette files.
//
// .wtone is wondertone's native, human-editable palette format built on TOML.
// It is the primary tool of the wondertone ecosystem — palettes live here,
// not in Go code, so designers can edit them without touching source.
//
// File format:
//
//	name        = "My Palette"
//	description = "A beautiful collection"
//	mood        = "joyful"
//	version     = "1.0.0"
//	author      = "you"
//
//	[[colors]]
//	name    = "Primary Spark"
//	l       = 0.75
//	c       = 0.15
//	h       = 30.0
//	energy  = 0.85
//	mood    = "vibrant"
//
//	[[colors]]
//	name  = "Accent Glow"
//	oklch = "0.60 0.20 200"   # shorthand: parsed into l/c/h on load
//	energy = 0.70
package wtone

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
)

// raw is the internal TOML structure — not exported.
type raw struct {
	Name        string     `toml:"name"`
	Description string     `toml:"description"`
	Mood        string     `toml:"mood"`
	Version     string     `toml:"version"`
	Author      string     `toml:"author"`
	Colors      []rawColor `toml:"colors"`
}

type rawColor struct {
	Name   string  `toml:"name"`
	L      float64 `toml:"l"`
	C      float64 `toml:"c"`
	H      float64 `toml:"h"`
	Energy float64 `toml:"energy"`
	Mood   string  `toml:"mood"`
	OKLCH  string  `toml:"oklch"` // shorthand: "L C H"
	Alpha  float64 `toml:"alpha"`
}

// LoadWTone parses a .wtone file and returns a Palette.
// Validates all colors on load — returns an error if any color is invalid.
func LoadWTone(path string) (*palette.Palette, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("wtone: cannot read %q: %w", path, err)
	}
	return ParseWTone(data)
}

// ParseWTone parses .wtone file contents from a byte slice.
// Useful for embedding .wtone files with //go:embed.
func ParseWTone(data []byte) (*palette.Palette, error) {
	var r raw
	if _, err := toml.Decode(string(data), &r); err != nil {
		return nil, fmt.Errorf("wtone: TOML parse error: %w", err)
	}
	if r.Name == "" {
		return nil, fmt.Errorf("wtone: palette must have a name")
	}
	if len(r.Colors) == 0 {
		return nil, fmt.Errorf("wtone: palette %q has no colors", r.Name)
	}

	b := palette.New(r.Name).
		Description(r.Description).
		Mood(r.Mood).
		Author(r.Author).
		Version(r.Version)

	for i, rc := range r.Colors {
		t, err := rawColorToTone(rc, r.Mood)
		if err != nil {
			return nil, fmt.Errorf("wtone: color[%d] in %q: %w", i, r.Name, err)
		}
		b.Add(t)
	}

	return b.Build()
}

// rawColorToTone converts a parsed TOML color entry to a Tone.
func rawColorToTone(rc rawColor, paletteMood string) (tone.Tone, error) {
	if rc.Name == "" {
		return tone.Tone{}, fmt.Errorf("every [[colors]] entry must have a name")
	}

	var l, c, h float64
	var err error

	if rc.OKLCH != "" {
		// Parse shorthand "L C H" or "L C H / A"
		l, c, h, _, err = parseOKLCHShorthand(rc.OKLCH)
		if err != nil {
			return tone.Tone{}, fmt.Errorf("invalid oklch shorthand %q: %w", rc.OKLCH, err)
		}
	} else {
		// Use explicit l/c/h fields
		l, c, h = rc.L, rc.C, rc.H
	}

	energy := rc.Energy
	if energy == 0 {
		energy = 1.0 // default
	}

	alpha := rc.Alpha
	if alpha == 0 {
		alpha = 1.0 // default
	}

	mood := rc.Mood
	if mood == "" {
		mood = paletteMood // inherit palette mood
	}

	t := tone.FromOKLCH(l, c, h).
		WithName(rc.Name).
		WithMood(mood).
		WithEnergy(energy).
		WithAlpha(alpha)

	return t, nil
}

// parseOKLCHShorthand parses "L C H" or "L C H / A".
func parseOKLCHShorthand(s string) (l, c, h, a float64, err error) {
	s = strings.TrimSpace(s)
	a = 1.0

	if idx := strings.Index(s, "/"); idx >= 0 {
		alphaPart := strings.TrimSpace(s[idx+1:])
		s = strings.TrimSpace(s[:idx])
		a, err = strconv.ParseFloat(alphaPart, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("invalid alpha %q", alphaPart)
		}
	}

	parts := strings.Fields(s)
	if len(parts) != 3 {
		return 0, 0, 0, 0, fmt.Errorf("expected 3 values (L C H), got %d", len(parts))
	}

	vals := [3]*float64{&l, &c, &h}
	names := [3]string{"L", "C", "H"}
	for i, p := range parts {
		*vals[i], err = strconv.ParseFloat(p, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("invalid %s value %q", names[i], p)
		}
	}
	return l, c, h, a, nil
}
