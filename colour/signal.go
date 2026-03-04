package colour

import tone "github.com/leraniode/wondertone/core"

// Signal is Leraniode's pure success green — the color of a passing test,
// a clean build, a zero exit code. Earned, not given.
var Signal = tone.New(
	tone.Light(65),
	tone.Vibrancy(80),
	tone.Hue(142),
	tone.Energy(1.0),
	tone.Named("Signal"),
	tone.Moody("focused"),
)
