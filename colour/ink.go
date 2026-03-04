package colour

import tone "github.com/leraniode/wondertone/core"

// Ink is Leraniode's near-black text tone — not pure black, which reads as
// harsh and flat. Ink has just enough blue to feel intentional on any background.
var Ink = tone.New(
	tone.Light(15),
	tone.Vibrancy(12),
	tone.Hue(230),
	tone.Energy(0.75),
	tone.Named("Ink"),
	tone.Moody("deep"),
)
