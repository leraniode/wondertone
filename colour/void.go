package colour

import tone "github.com/leraniode/wondertone/core"

// Void is deep-space indigo darkness — not black, but the richness of
// an infinite dark sky. Visible purple depth at high contrast.
var Void = tone.New(
	tone.Light(18),
	tone.Vibrancy(48),
	tone.Hue(268),
	tone.Energy(0.85),
	tone.Named("Void"),
	tone.Moody("deep"),
)
