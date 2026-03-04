package colour

import tone "github.com/leraniode/wondertone/core"

// Crimson is Leraniode's signal red — direct, urgent, impossible to ignore.
// For errors that deserve attention, not apology.
var Crimson = tone.New(
	tone.Light(55),
	tone.Vibrancy(88),
	tone.Hue(14),
	tone.Energy(1.0),
	tone.Named("Crimson"),
	tone.Moody("urgent"),
)
