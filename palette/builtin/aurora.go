package builtin

import (
	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
)

// Aurora returns a bright, airy light palette inspired by northern skies.
// Clean backgrounds, vivid accents, great for documentation and dashboards.
func Aurora() *palette.Palette {
	return palette.New("Aurora").
		Description("Bright airy light palette — clean, vivid, northern-sky inspired").
		Mood("serene").
		Author("leraniode").
		Version("1.0.0").
		Add(tone.New(tone.Light(97), tone.Vibrancy(8), tone.Hue(210), tone.Energy(0.80), tone.Named("Aurora Base"), tone.Moody("airy"))).
		Add(tone.New(tone.Light(93), tone.Vibrancy(12), tone.Hue(210), tone.Energy(0.78), tone.Named("Aurora Surface"), tone.Moody("airy"))).
		Add(tone.New(tone.Light(88), tone.Vibrancy(16), tone.Hue(210), tone.Energy(0.75), tone.Named("Aurora Overlay"), tone.Moody("airy"))).
		Add(tone.New(tone.Light(52), tone.Vibrancy(75), tone.Hue(199), tone.Energy(0.92), tone.Named("Aurora Accent"), tone.Moody("vivid"))).
		Add(tone.New(tone.Light(62), tone.Vibrancy(32), tone.Hue(210), tone.Energy(0.68), tone.Named("Aurora Muted"), tone.Moody("calm"))).
		Add(tone.New(tone.Light(18), tone.Vibrancy(22), tone.Hue(220), tone.Energy(0.85), tone.Named("Aurora Text"), tone.Moody("clear"))).
		Add(tone.New(tone.Light(58), tone.Vibrancy(68), tone.Hue(262), tone.Energy(0.85), tone.Named("Aurora Violet"), tone.Moody("mystical"))).
		Add(tone.New(tone.Light(55), tone.Vibrancy(70), tone.Hue(142), tone.Energy(0.88), tone.Named("Aurora Green"), tone.Moody("fresh"))).
		Add(tone.New(tone.Light(60), tone.Vibrancy(78), tone.Hue(14), tone.Energy(0.90), tone.Named("Aurora Red"), tone.Moody("urgent"))).
		MustBuild()
}
