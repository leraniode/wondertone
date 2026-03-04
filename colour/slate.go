package colour

import tone "github.com/leraniode/wondertone/core"

// Slate is Leraniode's neutral mid-tone — the color of good documentation,
// quiet comments, and things that don't need to shout.
var Slate = tone.New(
	tone.Light(55),
	tone.Vibrancy(14),
	tone.Hue(225),
	tone.Energy(0.70),
	tone.Named("Slate"),
	tone.Moody("calm"),
)
