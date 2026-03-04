// Package colour contains Leraniode's named tones — the wondertone colour collection.
//
// Every tone lives in its own file. Each is a carefully considered color with
// a name, a mood, and a reason to exist. These are not random swatches.
//
// Usage:
//
//	import "github.com/leraniode/wondertone/colour"
//
//	fmt.Println(colour.Unix.Hex())       // #...
//	fmt.Println(colour.Starlight.Name()) // "Starlight"
//
// The full collection:
//
//	Unix      — terminal green, focused
//	Starlight — deep indigo, mystical
//	Ember     — warm amber, comfortable
//	Glacier   — cool cyan, serene
//	Crimson   — signal red, urgent
//	Void      — near-black, deep
//	Dawn      — soft pink-orange, hopeful
//	Bloom     — vivid magenta, joyful
//	Slate     — neutral blue-grey, calm
//	Signal    — success green, earned
//	Ink       — near-black text, intentional
//	Paper     — warm off-white, easy
package colour

import tone "github.com/leraniode/wondertone/core"

// All returns every Leraniode named tone in the collection.
func All() []tone.Tone {
	return []tone.Tone{
		Unix,
		Starlight,
		Ember,
		Glacier,
		Crimson,
		Void,
		Dawn,
		Bloom,
		Slate,
		Signal,
		Ink,
		Paper,
	}
}
