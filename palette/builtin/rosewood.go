package builtin

import (
	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
)

// Rosewood returns a rich, warm rose-and-wood dark palette.
// Expressive and personal — for when the terminal should feel like home.
func Rosewood() *palette.Palette {
	return palette.New("Rosewood").
		Description("Rich rose and warm wood dark palette — expressive, personal, alive").
		Mood("vibrant").
		Author("leraniode").
		Version("1.0.0").
		Add(tone.New(tone.Light(13), tone.Vibrancy(20), tone.Hue(345), tone.Energy(0.88), tone.Named("Rosewood Base"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(18), tone.Vibrancy(16), tone.Hue(345), tone.Energy(0.82), tone.Named("Rosewood Surface"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(24), tone.Vibrancy(14), tone.Hue(345), tone.Energy(0.78), tone.Named("Rosewood Overlay"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(65), tone.Vibrancy(80), tone.Hue(345), tone.Energy(0.95), tone.Named("Rosewood Accent"), tone.Moody("vivid"))).
		Add(tone.New(tone.Light(50), tone.Vibrancy(35), tone.Hue(345), tone.Energy(0.65), tone.Named("Rosewood Muted"), tone.Moody("calm"))).
		Add(tone.New(tone.Light(90), tone.Vibrancy(12), tone.Hue(30), tone.Energy(0.82), tone.Named("Rosewood Text"), tone.Moody("clear"))).
		Add(tone.New(tone.Light(68), tone.Vibrancy(72), tone.Hue(20), tone.Energy(0.90), tone.Named("Rosewood Warm"), tone.Moody("warm"))).
		Add(tone.New(tone.Light(60), tone.Vibrancy(68), tone.Hue(300), tone.Energy(0.88), tone.Named("Rosewood Bloom"), tone.Moody("playful"))).
		Add(tone.New(tone.Light(64), tone.Vibrancy(65), tone.Hue(142), tone.Energy(0.85), tone.Named("Rosewood Leaf"), tone.Moody("fresh"))).
		MustBuild()
}
