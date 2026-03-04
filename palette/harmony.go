package palette

import (
	"fmt"

	tone "github.com/leraniode/wondertone/core"
)

// Analogous returns a palette of tones adjacent on the hue wheel.
// count tones, spread degrees apart, centered on the base tone.
// e.g. Analogous(base, 5, 30) gives 5 tones, 30° apart, centered on base.
func Analogous(base tone.Tone, count int, spreadDeg float64) (*Palette, error) {
	if count < 2 {
		return nil, fmt.Errorf("wondertone/palette: Analogous needs at least 2 tones")
	}
	tones := make([]tone.Tone, count)
	offset := spreadDeg * float64(count-1) / 2.0
	for i := 0; i < count; i++ {
		deg := -offset + spreadDeg*float64(i)
		tones[i] = base.Rotate(deg).
			WithName(fmt.Sprintf("%s analogous %+.0f°", base.Name(), deg)).
			WithMood(base.Mood())
	}
	return fromSlice(fmt.Sprintf("%s analogous", base.Name()), tones)
}

// Complementary returns a two-tone palette of a base and its complement (+180°).
func Complementary(base tone.Tone) (*Palette, error) {
	comp := base.Complement().
		WithName(fmt.Sprintf("%s complement", base.Name())).
		WithMood(base.Mood())
	return fromSlice(fmt.Sprintf("%s complementary", base.Name()), []tone.Tone{base, comp})
}

// Triadic returns a three-tone palette evenly spaced 120° apart.
func Triadic(base tone.Tone) (*Palette, error) {
	tones := []tone.Tone{
		base,
		base.Rotate(120).WithName(fmt.Sprintf("%s triadic-2", base.Name())),
		base.Rotate(240).WithName(fmt.Sprintf("%s triadic-3", base.Name())),
	}
	return fromSlice(fmt.Sprintf("%s triadic", base.Name()), tones)
}

// SplitComplementary returns a three-tone palette — base plus two tones
// flanking its complement at ±splitDeg.
func SplitComplementary(base tone.Tone, splitDeg float64) (*Palette, error) {
	if splitDeg <= 0 || splitDeg >= 90 {
		return nil, fmt.Errorf("wondertone/palette: SplitComplementary splitDeg must be in (0, 90)")
	}
	tones := []tone.Tone{
		base,
		base.Rotate(180 - splitDeg).WithName(fmt.Sprintf("%s split-1", base.Name())),
		base.Rotate(180 + splitDeg).WithName(fmt.Sprintf("%s split-2", base.Name())),
	}
	return fromSlice(fmt.Sprintf("%s split-complementary", base.Name()), tones)
}

// Tetradic returns a four-tone palette evenly spaced 90° apart.
func Tetradic(base tone.Tone) (*Palette, error) {
	tones := []tone.Tone{
		base,
		base.Rotate(90).WithName(fmt.Sprintf("%s tetradic-2", base.Name())),
		base.Rotate(180).WithName(fmt.Sprintf("%s tetradic-3", base.Name())),
		base.Rotate(270).WithName(fmt.Sprintf("%s tetradic-4", base.Name())),
	}
	return fromSlice(fmt.Sprintf("%s tetradic", base.Name()), tones)
}

// Monochrome returns a palette of count tones from the same hue,
// evenly distributed across the lightness scale.
func Monochrome(base tone.Tone, count int) (*Palette, error) {
	if count < 2 {
		return nil, fmt.Errorf("wondertone/palette: Monochrome needs at least 2 tones")
	}
	tones := make([]tone.Tone, count)
	for i := 0; i < count; i++ {
		// Distribute lightness from 15 to 92 (avoid pure black/white)
		light := 15.0 + (77.0 * float64(i) / float64(count-1))
		tones[i] = base.WithLight(light).
			WithName(fmt.Sprintf("%s mono-%d", base.Name(), i+1))
	}
	return fromSlice(fmt.Sprintf("%s monochrome", base.Name()), tones)
}

// Rainbow returns a palette of count tones, evenly distributed around the
// full hue wheel, at the same lightness and vibrancy as the base tone.
// The first tone starts at the base hue.
func Rainbow(base tone.Tone, count int) (*Palette, error) {
	if count < 2 {
		return nil, fmt.Errorf("wondertone/palette: Rainbow needs at least 2 tones")
	}
	tones := make([]tone.Tone, count)
	step := 360.0 / float64(count)
	for i := 0; i < count; i++ {
		tones[i] = base.Rotate(step * float64(i)).
			WithName(fmt.Sprintf("rainbow-%d", i+1))
	}
	return fromSlice(fmt.Sprintf("%s rainbow-%d", base.Name(), count), tones)
}

// fromSlice builds a Palette from a pre-built slice of Tones.
func fromSlice(name string, tones []tone.Tone) (*Palette, error) {
	b := New(name)
	for _, t := range tones {
		b.Add(t)
	}
	return b.Build()
}
