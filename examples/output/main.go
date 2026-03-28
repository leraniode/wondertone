// Package main demonstrates wondertone — the full developer experience.
//
// Run with:
//
//	go run examples/output/main.go
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

	fmt.Printf("%-20s  %-12s  hex=%-9s  Light=%.0f  Vibrancy=%.0f  Hue=%.0f  Energy=%.1f  Mood=%s\n",
		spark.Name(), render.Swatch(spark, profile, 2), spark.Hex(), spark.Light(), spark.Vibrancy(), spark.Hue(), spark.Energy(), spark.Mood())

	// Power user: raw OKLCH
	raw := tone.FromOKLCH(0.68, 0.18, 142)
	fmt.Printf("%-20s  %-12s  hex=%-9s  (from raw OKLCH)\n", "raw OKLCH", render.Swatch(raw, profile, 2), raw.Hex())

	// From hex
	legacy, _ := tone.FromHex("#e94560")
	fmt.Printf("%-20s  %-12s  hex=%-9s  Light=%.0f  Vibrancy=%.0f\n",
		"from #e94560", render.Swatch(legacy, profile, 2), legacy.Hex(), legacy.Light(), legacy.Vibrancy())

	fmt.Println()

	// ── 2. Leraniode named tones ─────────────────────────────────────────────
	fmt.Println("── 2. Wondertone colour collection ──────────────────────")
	for _, t := range colour.All() {
		swatch := render.Swatch(t, profile, 2)
		fmt.Printf("%s  %-12s  %-9s  mood=%-10s  temp=%s\n",
			swatch, t.Name(), t.Hex(), t.Mood(), t.Temperature())
	}
	fmt.Println()

	// ── 3. Tone Scale ────────────────────────────────────────────────────────
	fmt.Println("── 3. Tone scale (Crimson) ───────────────────────────")
	scale := colour.Crimson.Scale()
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
	fmt.Print("\n")

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
	custom, _ := palette.New("WonderPalette").
		Description("A custom palette built with wondertone").
		Author("DominionDev").
		Mood("vibrant").
		Add(colour.Unix).
		Add(colour.Starlight).
		Add(colour.Bloom).
		Add(colour.Ember).
		Add(colour.Crimson).
		Add(colour.Dawn).
		Build()

	report := custom.Validate()
	fmt.Printf("  Palette: %s (%d tones) ", custom.Name(), custom.Len())
	for _, t := range custom.All() {
		fmt.Print(render.Swatch(t, profile, 1) + " ")
	}
	fmt.Println()
	fmt.Printf("  Validation: %s\n", report.String())
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

	// ── 11. WonderMath — Temperature ─────────────────────────────────────────
	fmt.Println("── 11. WonderMath — Temperature ─────────────────────────")
	temperatures := []tone.Tone{
		tone.New(tone.Light(60), tone.Vibrancy(70), tone.Hue(25),  tone.Named("Ember")),
		tone.New(tone.Light(60), tone.Vibrancy(70), tone.Hue(60),  tone.Named("Gold")),
		tone.New(tone.Light(60), tone.Vibrancy(70), tone.Hue(142), tone.Named("Unix")),
		tone.New(tone.Light(60), tone.Vibrancy(70), tone.Hue(196), tone.Named("Glacier")),
		tone.New(tone.Light(60), tone.Vibrancy(70), tone.Hue(240), tone.Named("Starlight")),
	}
	for _, t := range temperatures {
		fmt.Printf("  %-12s  %-10s  hex=%-9s  temp=%-8s  scalar=%+.2f\n",
			t.Name(), render.Swatch(t, profile, 2), t.Hex(), t.Temperature(), t.TemperatureScalar())
	}
	fmt.Println()

	// ── 12. WonderMath — Mood as math ────────────────────────────────────────
	fmt.Println("── 12. WonderMath — Derived Mood ────────────────────────")
	moodTones := []tone.Tone{
		tone.New(tone.Light(75), tone.Vibrancy(95), tone.Hue(40),  tone.Energy(1.0), tone.Named("Sunrise")),
		tone.New(tone.Light(30), tone.Vibrancy(40), tone.Hue(250), tone.Energy(0.3), tone.Named("Midnight")),
		tone.New(tone.Light(85), tone.Vibrancy(20), tone.Hue(50),  tone.Energy(0.6), tone.Named("Linen")),
		tone.New(tone.Light(55), tone.Vibrancy(85), tone.Hue(320), tone.Energy(0.9), tone.Named("Bloom")),
		tone.New(tone.Light(50), tone.Vibrancy(65), tone.Hue(270), tone.Energy(0.5), tone.Named("Dusk")),
	}
	for _, t := range moodTones {
		fmt.Printf("  %-12s %-10s  hex=%-9s  mood=%-10s  valence=%+.2f  arousal=%+.2f\n",
			t.Name(), render.Swatch(t, profile, 3), t.Hex(), t.DerivedMoodValue(), t.ValenceValue(), t.ArousalValue())
	}
	fmt.Println()

	// ── 13. WonderMath — Energy Stevens law ──────────────────────────────────
	fmt.Println("── 13. WonderMath — Energy (Stevens' power law) ─────────")
	base := colour.Bloom
	fmt.Printf("  base            hex=%s  effectiveC=%.4f\n", base.Hex(), base.EffectiveC())
	for _, e := range []float64{0.8, 0.6, 0.4, 0.2, 0.0} {
		t := base.WithEnergy(e)
		linearC := base.EffectiveC() * e
		fmt.Printf("  energy=%.1f  %s  hex=%s  effectiveC=%.4f  (linear would be %.4f)\n",
			e, render.Swatch(t, profile, 3), t.Hex(), t.EffectiveC(), linearC)
	}
	fmt.Println()

	// ── 14. WonderMath — Perceptual vibrancy equality ────────────────────────
	fmt.Println("── 14. WonderMath — Perceptual Vibrancy Equality ────────")
	hues := []struct{ name string; hue float64 }{
		{"Red",     0},
		{"Orange",  30},
		{"Yellow",  60},
		{"Green",   142},
		{"Cyan",    180},
		{"Blue",    240},
		{"Magenta", 300},
	}
	for _, h := range hues {
		t := tone.New(tone.Light(65), tone.Vibrancy(80), tone.Hue(h.hue))
		fmt.Printf("  %-8s  H=%3.0f  %s  hex=%s\n",
			h.name, h.hue, render.Swatch(t, profile, 3), t.Hex())
	}
	fmt.Println("  Each swatch is Vibrancy=80 — perceptually equalised across hues")
	fmt.Println()
}
