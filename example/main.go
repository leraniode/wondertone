// Package main demonstrates wondertone — the full developer experience.
//
// Run with:
//
//	cd example && go run main.go
package main

import (
	"fmt"

	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/colour"
	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/palette/builtin"
	"github.com/leraniode/wondertone/render"
	"github.com/leraniode/wondertone/wtone"
)

func main() {
	profile := render.Detect()
	fmt.Printf("Terminal profile: %s\n\n", profile)

	// ── 1. Create a tone with the Wondertone vocabulary ─────────────────────
	fmt.Println("── 1. Creating tones ───────────────────────────────────")

	spark := tone.New(
		tone.Light(75),
		tone.Vibrancy(80),
		tone.Hue(30),
		tone.Energy(0.9),
		tone.Named("Primary Spark"),
		tone.Moody("vibrant"),
	)

	fmt.Printf("%-20s  hex=%-9s  Light=%.0f  Vibrancy=%.0f  Hue=%.0f  Energy=%.1f  Mood=%s\n",
		spark.Name(), spark.Hex(), spark.Light(), spark.Vibrancy(), spark.Hue(), spark.Energy(), spark.Mood())

	// Power user: raw OKLCH
	raw := tone.FromOKLCH(0.68, 0.18, 142)
	fmt.Printf("%-20s  hex=%-9s  (from raw OKLCH)\n", "raw OKLCH", raw.Hex())

	// From hex
	legacy, _ := tone.FromHex("#e94560")
	fmt.Printf("%-20s  hex=%-9s  Light=%.0f  Vibrancy=%.0f\n",
		"from #e94560", legacy.Hex(), legacy.Light(), legacy.Vibrancy())

	fmt.Println()

	// ── 2. Leraniode named tones ─────────────────────────────────────────────
	fmt.Println("── 2. Leraniode colour collection ──────────────────────")
	for _, t := range colour.All() {
		swatch := render.Swatch(t, profile, 2)
		fmt.Printf("%s  %-12s  %-9s  mood=%-10s  temp=%s\n",
			swatch, t.Name(), t.Hex(), t.Mood(), t.Temperature())
	}
	fmt.Println()

	// ── 3. Tone Scale ────────────────────────────────────────────────────────
	fmt.Println("── 3. Tone scale (Unix green) ───────────────────────────")
	scale := colour.Unix.Scale()
	roles := []string{"Background", "SubtleBackground", "ElementBackground", "HoveredBackground",
		"ActiveBackground", "SubtleBorder", "Border", "StrongBorder",
		"Solid", "HoveredSolid", "Text", "HighContrastText"}
	for i, t := range scale.All() {
		swatch := render.Swatch(t, profile, 2)
		fmt.Printf("  %2d  %s  %-22s  Light=%.0f\n", i+1, swatch, roles[i], t.Light())
	}
	fmt.Println()

	// ── 4. Mixing ────────────────────────────────────────────────────────────
	fmt.Println("── 4. Mixing in OKLab space ─────────────────────────────")
	a := colour.Unix
	b := colour.Starlight
	grad, _ := tone.Gradient(a, b, 7)
	fmt.Printf("  %s  --------→  gradient  →  %s\n  ", render.Swatch(a, profile, 1), render.Swatch(b, profile, 1))
	for _, step := range grad {
		fmt.Print(render.Swatch(step, profile, 1) + " ")
	}
	fmt.Println("\n")

	// ── 5. Harmony ───────────────────────────────────────────────────────────
	fmt.Println("── 5. Harmony (triadic from Crimson) ────────────────────")
	triad, _ := tone.Harmonize(colour.Crimson, "triadic")
	for _, t := range triad {
		fmt.Printf("  %s  hue=%.0f\n", render.Swatch(t, profile, 2), t.Hue())
	}
	fmt.Println()

	// ── 6. Built-in palettes ─────────────────────────────────────────────────
	fmt.Println("── 6. Built-in palettes ─────────────────────────────────")
	for _, p := range builtin.All() {
		fmt.Printf("  %-12s  %d tones  mood=%-10s  ", p.Name(), p.Len(), p.Mood())
		for _, t := range p.All()[:5] {
			fmt.Print(render.Swatch(t, profile, 1) + " ")
		}
		fmt.Println()
	}
	fmt.Println()

	// ── 7. Build your own palette ────────────────────────────────────────────
	fmt.Println("── 7. Build a custom palette ────────────────────────────")
	custom, _ := palette.New("My Palette").
		Description("A custom palette built with wondertone").
		Mood("vibrant").
		Add(colour.Unix).
		Add(colour.Starlight).
		Add(colour.Bloom).
		Add(colour.Ember).
		Build()

	report := custom.Validate()
	fmt.Printf("  Palette: %s (%d tones)\n", custom.Name(), custom.Len())
	fmt.Printf("  Validation: %s", report.String())
	fmt.Println()

	// ── 8. Accessibility ────────────────────────────────────────────────────
	fmt.Println("── 8. Accessibility ─────────────────────────────────────")
	midnight := builtin.Midnight()
	fg := midnight.MustGet("Midnight Text")
	bg := midnight.MustGet("Midnight Base")
	ratio := fg.ContrastWith(bg)
	passAA := fg.PassesAA(bg)
	fmt.Printf("  Text on Base: %.2f:1  AA=%v\n", ratio, passAA)

	// Fix contrast if needed
	if !passAA {
		fixed := fg.EnsureContrast(bg, "AA")
		fmt.Printf("  Fixed:        %.2f:1  AA=%v\n", fixed.ContrastWith(bg), fixed.PassesAA(bg))
	}

	// Readable pairs
	pairs := palette.FindReadablePairs(midnight, "AA")
	fmt.Printf("  Readable AA pairs in Midnight: %d\n", len(pairs))
	fmt.Println()

	// ── 9. Energy ────────────────────────────────────────────────────────────
	fmt.Println("── 9. Energy — same tones, different aliveness ──────────")
	energyLevels := []float64{1.0, 0.7, 0.4, 0.1}
	for _, e := range energyLevels {
		quiet := colour.Bloom.WithEnergy(e)
		fmt.Printf("  energy=%.1f  %s  hex=%s\n", e, render.Swatch(quiet, profile, 3), quiet.Hex())
	}
	fmt.Println()

	// ── 10. .wtone file ──────────────────────────────────────────────────────
	fmt.Println("── 10. .wtone file ──────────────────────────────────────")
	wtone.SaveWTone("assets/tone.wtone", custom)
	fmt.Println("file created: ./assets/tone.wtone")
}
