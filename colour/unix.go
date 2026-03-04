// Package colour contains Leraniode's named tones — the wondertone colour collection.
//
// Each file is one tone. Import the ones you need:
//
//	import "github.com/leraniode/wondertone/colour"
//
//	fmt.Println(colour.Unix.Hex())
package colour

import tone "github.com/leraniode/wondertone/core"

// Unix is Leraniode's signature terminal green.
// Focused, precise, alive — the color of the cursor that never blinks out.
// Named for the philosophy that earned it: do one thing, do it well.
var Unix = tone.New(
	tone.Light(68),
	tone.Vibrancy(72),
	tone.Hue(142),
	tone.Energy(0.95),
	tone.Named("Unix"),
	tone.Moody("focused"),
)
