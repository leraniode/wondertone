package palette

import (
	"fmt"

	tone "github.com/leraniode/wondertone/core"
)

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

// validate runs all quality checks on a Palette.
func validate(p *Palette) ValidationReport {
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
		PaletteName: p.name,
		Issues:      issues,
		Passed:      len(issues) == 0,
	}
}

// ContrastPair checks the WCAG contrast ratio between two named Tones in a Palette.
func ContrastPair(p *Palette, fgName, bgName string) (float64, error) {
	fg, ok := p.Get(fgName)
	if !ok {
		return 0, fmt.Errorf("wondertone/palette: no Tone named %q", fgName)
	}
	bg, ok := p.Get(bgName)
	if !ok {
		return 0, fmt.Errorf("wondertone/palette: no Tone named %q", bgName)
	}
	return fg.ContrastWith(bg), nil
}

// EnsurePairContrast returns a new Palette where the foreground Tone has been
// adjusted to meet the given WCAG level ("AA" or "AAA") against the background Tone.
// Only lightness is adjusted — hue, vibrancy, and energy are preserved.
func EnsurePairContrast(p *Palette, fgName, bgName, level string) (*Palette, error) {
	fg, ok := p.Get(fgName)
	if !ok {
		return nil, fmt.Errorf("wondertone/palette: no Tone named %q", fgName)
	}
	bg, ok := p.Get(bgName)
	if !ok {
		return nil, fmt.Errorf("wondertone/palette: no Tone named %q", bgName)
	}
	fixed := fg.EnsureContrast(bg, level)
	return p.Replace(fgName, fixed)
}

// ContrastMatrix returns the contrast ratio between every pair of Tones in the palette.
// The result is a map of "fgName/bgName" → ratio.
func ContrastMatrix(p *Palette) map[string]float64 {
	tones := p.All()
	result := make(map[string]float64, len(tones)*len(tones))
	for _, fg := range tones {
		for _, bg := range tones {
			if fg.Name() == bg.Name() {
				continue
			}
			key := fg.Name() + "/" + bg.Name()
			result[key] = fg.ContrastWith(bg)
		}
	}
	return result
}

// FindReadablePairs returns all (fg, bg) Tone pairs in the palette that meet
// the given WCAG level. Useful for discovering which color combinations work.
func FindReadablePairs(p *Palette, level string) []ContrastResult {
	target := 4.5
	if level == "AAA" {
		target = 7.0
	}
	tones := p.All()
	var results []ContrastResult
	for _, fg := range tones {
		for _, bg := range tones {
			if fg.Name() == bg.Name() {
				continue
			}
			ratio := fg.ContrastWith(bg)
			if ratio >= target {
				results = append(results, ContrastResult{
					FG:       fg,
					BG:       bg,
					Ratio:    ratio,
					PassesAA: ratio >= 4.5,
					PassesAAA: ratio >= 7.0,
				})
			}
		}
	}
	return results
}

// ContrastResult holds one FG/BG contrast result.
type ContrastResult struct {
	FG        tone.Tone
	BG        tone.Tone
	Ratio     float64
	PassesAA  bool
	PassesAAA bool
}

func (r ContrastResult) String() string {
	level := "✗ FAIL"
	if r.PassesAAA {
		level = "✓ AAA"
	} else if r.PassesAA {
		level = "✓ AA"
	}
	return fmt.Sprintf("%-20s on %-20s  %.2f:1  %s",
		r.FG.Name(), r.BG.Name(), r.Ratio, level)
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
