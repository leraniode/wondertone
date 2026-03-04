package colour

import tone "github.com/leraniode/wondertone/core"

// Void is Leraniode's deep background — not black, but the color of a room
// with the lights off and a monitor glowing. Almost nothing. Almost.
var Void = tone.New(
	tone.Light(10),
	tone.Vibrancy(15),
	tone.Hue(240),
	tone.Energy(0.80),
	tone.Named("Void"),
	tone.Moody("deep"),
)
