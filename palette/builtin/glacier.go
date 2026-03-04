package builtin

import (
	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
)

// Glacier returns a cool, icy palette with cyan and blue tones.
// Calm and precise — like light through deep arctic ice.
func Glacier() *palette.Palette {
	return palette.New("Glacier").
		Description("Cool icy palette — calm, precise, arctic-inspired").
		Mood("serene").
		Author("leraniode").
		Version("1.0.0").
		Add(tone.New(tone.Light(14), tone.Vibrancy(22), tone.Hue(196), tone.Energy(0.85), tone.Named("Glacier Base"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(19), tone.Vibrancy(18), tone.Hue(196), tone.Energy(0.80), tone.Named("Glacier Surface"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(25), tone.Vibrancy(16), tone.Hue(196), tone.Energy(0.75), tone.Named("Glacier Overlay"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(74), tone.Vibrancy(68), tone.Hue(196), tone.Energy(0.90), tone.Named("Glacier Accent"), tone.Moody("vivid"))).
		Add(tone.New(tone.Light(50), tone.Vibrancy(32), tone.Hue(196), tone.Energy(0.65), tone.Named("Glacier Muted"), tone.Moody("calm"))).
		Add(tone.New(tone.Light(88), tone.Vibrancy(12), tone.Hue(200), tone.Energy(0.82), tone.Named("Glacier Text"), tone.Moody("clear"))).
		Add(tone.New(tone.Light(62), tone.Vibrancy(60), tone.Hue(220), tone.Energy(0.85), tone.Named("Glacier Blue"), tone.Moody("calm"))).
		Add(tone.New(tone.Light(66), tone.Vibrancy(62), tone.Hue(142), tone.Energy(0.88), tone.Named("Glacier Green"), tone.Moody("focused"))).
		Add(tone.New(tone.Light(68), tone.Vibrancy(70), tone.Hue(280), tone.Energy(0.85), tone.Named("Glacier Violet"), tone.Moody("mystical"))).
		MustBuild()
}
