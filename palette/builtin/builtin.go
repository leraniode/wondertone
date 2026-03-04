package builtin

import "github.com/leraniode/wondertone/palette"

// All returns every built-in palette.
// Useful for listing, previewing, or choosing a palette at runtime.
func All() []*palette.Palette {
	return []*palette.Palette{
		Midnight(),
		Aurora(),
		Ember(),
		Glacier(),
		Rosewood(),
	}
}

// Names returns the name of every built-in palette, in order.
func Names() []string {
	return []string{"Midnight", "Aurora", "Ember", "Glacier", "Rosewood"}
}
