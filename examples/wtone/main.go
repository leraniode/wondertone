// package main demostrates the use of the `.wtone` file
//
// Run with
//
// cd examples/wtone && go run main.go

package main

import (
	"fmt"

	"github.com/leraniode/wondertone/render"
	wtone "github.com/leraniode/wondertone/wtone"
)

func main() {
    // set profile
	profile := render.Detect()

	// Load assets/tone.wtone and panic on error
	p, err := wtone.LoadWTone("assets/tone.wtone")
	if err != nil {
		panic(err)
	}

	// render the `Starlight` tone in the `wtone` file and panic on error
    Starlight := p.MustGet("Starlight")
	renderstar := render.Swatch(Starlight, profile, 2)
	fmt.Printf("Starlight: %-20s\n", renderstar)

	// render the `Unix` tone in the `wtone` file and panic on error
	Unix := p.MustGet("Unix")
	renderUnix := render.Swatch(Unix, profile, 2)
	fmt.Printf("Unix: %-20s\n", renderUnix)

	// render the `Bloom` tone
	Bloom := p.MustGet("Bloom")
	renderBloom := render.Swatch(Bloom, profile, 2)
	fmt.Printf("Bloom: %-20s\n", renderBloom)

	// render the `Ember` tone
	Ember := p.MustGet("Ember")
	renderEmber := render.Swatch(Ember, profile, 2)
	fmt.Printf("Ember: %-20s\n", renderEmber)
}