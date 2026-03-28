package colour

import tone "github.com/leraniode/wondertone/core"

// Slate is cool blue-grey stone — calm and grounded, with enough colour
// identity to feel distinctly cooler than a neutral grey.
var Slate = tone.New(
	tone.Light(57),
	tone.Vibrancy(40),
	tone.Hue(218),
	tone.Energy(0.75),
	tone.Named("Slate"),
	tone.Moody("calm"),
)
