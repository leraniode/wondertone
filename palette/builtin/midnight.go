// Package builtin contains wondertone's built-in palettes.
// Each palette is ready to use — no configuration needed.
//
//	import "github.com/leraniode/wondertone/palette/builtin"
//
//	p := builtin.Midnight()
//	fmt.Println(p.MustGet("Midnight Base").Hex())
package builtin

import (
	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
)

// Midnight returns a deep dark palette — rich navy blues and glowing accents.
// Perfect for terminal UIs that want depth without harshness.
func Midnight() *palette.Palette {
	return palette.New("Midnight").
		Description("Deep dark navy palette — focused, atmospheric, alive at night").
		Mood("mystical").
		Author("leraniode").
		Version("1.0.0").
		Add(tone.New(tone.Light(12), tone.Vibrancy(28), tone.Hue(240), tone.Energy(0.90), tone.Named("Midnight Base"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(16), tone.Vibrancy(22), tone.Hue(240), tone.Energy(0.85), tone.Named("Midnight Surface"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(22), tone.Vibrancy(18), tone.Hue(240), tone.Energy(0.80), tone.Named("Midnight Overlay"), tone.Moody("deep"))).
		Add(tone.New(tone.Light(68), tone.Vibrancy(72), tone.Hue(199), tone.Energy(0.95), tone.Named("Midnight Accent"), tone.Moody("vivid"))).
		Add(tone.New(tone.Light(55), tone.Vibrancy(40), tone.Hue(240), tone.Energy(0.70), tone.Named("Midnight Muted"), tone.Moody("calm"))).
		Add(tone.New(tone.Light(88), tone.Vibrancy(12), tone.Hue(240), tone.Energy(0.85), tone.Named("Midnight Text"), tone.Moody("clear"))).
		Add(tone.New(tone.Light(62), tone.Vibrancy(65), tone.Hue(262), tone.Energy(0.88), tone.Named("Midnight Purple"), tone.Moody("mystical"))).
		Add(tone.New(tone.Light(70), tone.Vibrancy(68), tone.Hue(142), tone.Energy(0.90), tone.Named("Midnight Green"), tone.Moody("focused"))).
		Add(tone.New(tone.Light(72), tone.Vibrancy(80), tone.Hue(38), tone.Energy(0.88), tone.Named("Midnight Amber"), tone.Moody("warm"))).
		MustBuild()
}
