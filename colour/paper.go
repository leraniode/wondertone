package colour

import tone "github.com/leraniode/wondertone/core"

// Paper is Leraniode's near-white background — not pure white, which burns
// on a dark terminal and looks clinical everywhere else. Paper is warm,
// off-white, and easy to read on all day.
var Paper = tone.New(
	tone.Light(95),
	tone.Vibrancy(6),
	tone.Hue(55),
	tone.Energy(0.65),
	tone.Named("Paper"),
	tone.Moody("airy"),
)
