// Package palette manages named collections of Tones.
//
// Import as "palette" for clean usage:
//
//	import palette "github.com/leraniode/wondertone/palette"
//
//	p := palette.New("midnight").
//	    Add(unix).
//	    Add(aurora).
//	    Build()
package palette

import (
	"fmt"

	"github.com/leraniode/wondertone/tone"
)

// Palette is an ordered, named collection of Tones.
// Order is preserved — first added, first in the list.
type Palette struct {
	name        string
	description string
	mood        string
	author      string
	version     string
	tones       []tone.Tone
	index       map[string]int // name → position, for fast lookup
}

// Builder constructs a Palette with a fluent API.
type Builder struct {
	p   *Palette
	err error
}

// New begins building a Palette with the given name.
func New(name string) *Builder {
	return &Builder{
		p: &Palette{
			name:    name,
			version: "1.0.0",
			index:   make(map[string]int),
		},
	}
}

// Description sets a human-readable description.
func (b *Builder) Description(s string) *Builder { b.p.description = s; return b }

// Mood sets the palette's overall mood — cascades to unnamed Tones.
func (b *Builder) Mood(s string) *Builder { b.p.mood = s; return b }

// Author sets the palette's author.
func (b *Builder) Author(s string) *Builder { b.p.author = s; return b }

// Version sets the palette's version string.
func (b *Builder) Version(s string) *Builder { b.p.version = s; return b }

// Add appends a Tone to the palette. Tones must have unique names.
func (b *Builder) Add(t tone.Tone) *Builder {
	if b.err != nil {
		return b
	}
	name := t.Name()
	if name == "" {
		b.err = fmt.Errorf("wondertone/palette: every Tone must have a name — use tone.Named()")
		return b
	}
	if _, exists := b.p.index[name]; exists {
		b.err = fmt.Errorf("wondertone/palette: duplicate Tone name %q", name)
		return b
	}
	b.p.index[name] = len(b.p.tones)
	b.p.tones = append(b.p.tones, t)
	return b
}

// Build finalises and validates the Palette. Returns an error if any
// constraint was violated during construction.
func (b *Builder) Build() (*Palette, error) {
	if b.err != nil {
		return nil, b.err
	}
	if len(b.p.tones) == 0 {
		return nil, fmt.Errorf("wondertone/palette: palette %q has no tones", b.p.name)
	}
	return b.p, nil
}

// MustBuild is like Build but panics on error. Safe for package-level vars.
func (b *Builder) MustBuild() *Palette {
	p, err := b.Build()
	if err != nil {
		panic(err)
	}
	return p
}

// --- Palette accessors ---

func (p *Palette) Name() string        { return p.name }
func (p *Palette) Description() string { return p.description }
func (p *Palette) Mood() string        { return p.mood }
func (p *Palette) Author() string      { return p.author }
func (p *Palette) Version() string     { return p.version }
func (p *Palette) Len() int            { return len(p.tones) }

// All returns a copy of all Tones in order.
func (p *Palette) All() []tone.Tone {
	out := make([]tone.Tone, len(p.tones))
	copy(out, p.tones)
	return out
}

// Get retrieves a Tone by name. Returns false if not found.
func (p *Palette) Get(name string) (tone.Tone, bool) {
	i, ok := p.index[name]
	if !ok {
		return tone.Tone{}, false
	}
	return p.tones[i], true
}

// MustGet retrieves a Tone by name. Panics if not found.
func (p *Palette) MustGet(name string) tone.Tone {
	t, ok := p.Get(name)
	if !ok {
		panic(fmt.Sprintf("wondertone/palette: %q has no Tone named %q", p.name, name))
	}
	return t
}

// At returns the Tone at position i (0-based). Panics if out of range.
func (p *Palette) At(i int) tone.Tone {
	if i < 0 || i >= len(p.tones) {
		panic(fmt.Sprintf("wondertone/palette: index %d out of range for palette %q (len %d)", i, p.name, len(p.tones)))
	}
	return p.tones[i]
}

// Has reports whether a Tone with the given name exists.
func (p *Palette) Has(name string) bool {
	_, ok := p.index[name]
	return ok
}

// --- Palette operations ---

// Fork creates a Builder pre-populated with all Tones from this Palette,
// ready to be modified and rebuilt under a new name.
func (p *Palette) Fork(newName string) *Builder {
	b := New(newName).
		Description(p.description).
		Mood(p.mood).
		Author(p.author).
		Version("1.0.0")
	for _, t := range p.tones {
		b.Add(t)
	}
	return b
}

// Extend returns a new Palette with additional Tones appended.
// Cannot override existing Tone names — use Fork for that.
func (p *Palette) Extend(name string, extras ...tone.Tone) (*Palette, error) {
	b := p.Fork(name)
	for _, t := range extras {
		if p.Has(t.Name()) {
			return nil, fmt.Errorf("wondertone/palette: Extend cannot override existing Tone %q — use Fork", t.Name())
		}
		b.Add(t)
	}
	return b.Build()
}

// Replace returns a new Palette with the named Tone swapped out.
func (p *Palette) Replace(name string, replacement tone.Tone) (*Palette, error) {
	if !p.Has(name) {
		return nil, fmt.Errorf("wondertone/palette: %q has no Tone named %q", p.name, name)
	}
	b := New(p.name).
		Description(p.description).
		Mood(p.mood).
		Author(p.author).
		Version(p.version)
	for _, t := range p.tones {
		if t.Name() == name {
			b.Add(replacement.WithName(name))
		} else {
			b.Add(t)
		}
	}
	return b.Build()
}

// WithEnergy returns a new Palette with Energy applied to every Tone.
// This is the "mood dial" — quieten or energise an entire palette at once.
func (p *Palette) WithEnergy(e float64) *Palette {
	b := New(p.name).
		Description(p.description).
		Mood(p.mood).
		Author(p.author).
		Version(p.version)
	for _, t := range p.tones {
		b.Add(t.WithEnergy(e))
	}
	result, _ := b.Build() // safe — same tones, same names
	return result
}

// ValidationReport summarises quality checks for a Palette.
type ValidationReport struct {
	PaletteName string
	Issues      []string
	Passed      bool
}

func (r ValidationReport) String() string {
	if r.Passed {
		return fmt.Sprintf("✓ %s — all checks passed", r.PaletteName)
	}
	out := fmt.Sprintf("✗ %s — %d issue(s):\n", r.PaletteName, len(r.Issues))
	for _, issue := range r.Issues {
		out += fmt.Sprintf("  • %s\n", issue)
	}
	return out
}

// Validate runs quality checks on the palette. Returns a Report.
func (p *Palette) Validate() ValidationReport {
	var issues []string

	tones := p.All()

	// Minimum size
	if len(tones) < 2 {
		issues = append(issues, "palette should have at least 2 tones")
	}

	// Recommended size warning (not a hard failure)
	if len(tones) > 16 {
		issues = append(issues,
			fmt.Sprintf("palette has %d tones — consider splitting into sub-palettes (recommended max: 16)", len(tones)),
		)
	}

	// All tones must be in sRGB gamut
	for _, t := range tones {
		r, g, b := t.RGBFloat()
		if r < -0.001 || r > 1.001 || g < -0.001 || g > 1.001 || b < -0.001 || b > 1.001 {
			issues = append(issues,
				fmt.Sprintf("tone %q is outside sRGB gamut — call ToGamutSafe or reduce Vibrancy", t.Name()),
			)
		}
	}

	// Adjacent tones should be perceptually distinct (ΔL ≥ 5)
	for i := 1; i < len(tones); i++ {
		prev, curr := tones[i-1], tones[i]
		if absDiff(prev.Light(), curr.Light()) < 5 {
			issues = append(issues,
				fmt.Sprintf("tones %q and %q are very similar (ΔLight < 5) — may be hard to distinguish", prev.Name(), curr.Name()),
			)
		}
	}

	return ValidationReport{
		PaletteName: p.Name(),
		Issues:      issues,
		Passed:      len(issues) == 0,
	}
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
