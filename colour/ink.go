package colour

import tone "github.com/leraniode/wondertone/core"

// Ink is the blue-black of a fine pen — dark but with unmistakable
// cool character. Distinctly different from Void's purple depth.
var Ink = tone.New(
	tone.Light(20),
	tone.Vibrancy(30),
	tone.Hue(232),
	tone.Energy(0.80),
	tone.Named("Ink"),
	tone.Moody("deep"),
)
