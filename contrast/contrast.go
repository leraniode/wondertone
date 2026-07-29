// Package contrast provides APCA-based colour contrast and accessibility tools.
//
// Contrast ratios, AA/AAA validation, contrast matrices, and readable pair
// finding — all computed using the APCA (Advanced Perceptual Contrast Algorithm).
//
//	import "github.com/leraniode/wondertone/contrast"
//
//	ratio := contrast.APCARatio(text, bg)
//	ok    := contrast.PassesAA(text, bg)
//	pairs := contrast.FindReadablePairs(p, "AA")
package contrast

import (
	"fmt"

	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/tone"
)

// ContrastPair checks the WCAG contrast ratio between two named Tones in a Palette.
func ContrastPair(p *palette.Palette, fgName, bgName string) (float64, error) {
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
func EnsurePairContrast(p *palette.Palette, fgName, bgName, level string) (*palette.Palette, error) {
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
func ContrastMatrix(p *palette.Palette) map[string]float64 {
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
func FindReadablePairs(p *palette.Palette, level string) []ContrastResult {
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
					FG:        fg,
					BG:        bg,
					Ratio:     ratio,
					PassesAA:  ratio >= 4.5,
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
