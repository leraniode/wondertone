package colour

import tone "github.com/leraniode/wondertone/core"

// Signal is electric teal-green — the colour of "on", of active status,
// of systems alive and responding. Distinct from Unix's natural green.
var Signal = tone.New(
	tone.Light(62),
	tone.Vibrancy(94),
	tone.Hue(165),
	tone.Energy(1.0),
	tone.Named("Signal"),
	tone.Moody("focused"),
)
