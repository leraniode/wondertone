package builtin

import (
	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
)

// Ember returns a warm, amber-toned dark palette.
// Like the last light of a fire — comfortable, focused, never harsh.
func Ember() *palette.Palette {
	return palette.New("Ember").
		Description("Warm amber dark palette — cozy, focused, fire-lit").
		Mood("warm").
		Author("leraniode").
		Version("1.0.0").
		Add(tone.New(tone.Light(13), tone.Vibrancy(25), tone.Hue(30), tone.Energy(0.88), tone.Named("Ember Base"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(18), tone.Vibrancy(20), tone.Hue(30), tone.Energy(0.82), tone.Named("Ember Surface"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(24), tone.Vibrancy(18), tone.Hue(30), tone.Energy(0.78), tone.Named("Ember Overlay"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(72), tone.Vibrancy(88), tone.Hue(38), tone.Energy(0.95), tone.Named("Ember Accent"), tone.Moody("vivid"))).
		Add(tone.New(tone.Light(52), tone.Vibrancy(38), tone.Hue(30), tone.Energy(0.68), tone.Named("Ember Muted"), tone.Moody("calm"))).
		Add(tone.New(tone.Light(90), tone.Vibrancy(15), tone.Hue(38), tone.Energy(0.82), tone.Named("Ember Text"), tone.Moody("clear"))).
		Add(tone.New(tone.Light(68), tone.Vibrancy(72), tone.Hue(60), tone.Energy(0.90), tone.Named("Ember Gold"), tone.Moody("warm"))).
		Add(tone.New(tone.Light(58), tone.Vibrancy(80), tone.Hue(14), tone.Energy(0.92), tone.Named("Ember Red"), tone.Moody("urgent"))).
		Add(tone.New(tone.Light(65), tone.Vibrancy(65), tone.Hue(142), tone.Energy(0.85), tone.Named("Ember Green"), tone.Moody("focused"))).
		MustBuild()
}
