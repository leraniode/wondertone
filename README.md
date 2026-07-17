<p align="left">
    <a href="https://github.com/leraniode/wondertone">
        <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/images/wondertone.svg" width="1024" />
    </a>
</p>

# ✨ wondertone 🌈🎨

A perceptual color intelligence library for Go — OKLCH under the hood, a human vocabulary on the surface.

[![CI](https://github.com/leraniode/wondertone/actions/workflows/ci.yml/badge.svg)](https://github.com/leraniode/wondertone/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/leraniode/wondertone.svg)](https://pkg.go.dev/github.com/leraniode/wondertone)
[![License](https://img.shields.io/github/license/leraniode/wondertone)](https://github.com/leraniode/wondertone/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/leraniode/wondertone)](https://goreportcard.com/report/github.com/leraniode/wondertone)
[![Version](https://img.shields.io/github/v/tag/leraniode/wondertone)](https://github.com/leraniode/wondertone/releases)
[![lipgloss adapter version](https://pkg.go.dev/github.com/leraniode/wondertone/adapters/lipgloss?tab=versions&color=blue)](https://pkg.go.dev/github.com/leraniode/wondertone/adapters/lipgloss)
[![go-colourful adapter version](https://pkg.go.dev/github.com/leraniode/wondertone/adapters/colorful?tab=versions&color=blue)](https://pkg.go.dev/github.com/leraniode/wondertone/adapters/colorful)
[![Go Modules](https://img.shields.io/github/go-mod/go-version/leraniode/wondertone)](https://github.com/leraniode/wondertone/blob/main/go.mod)

> [!CAUTION]
> **Breaking changes in v0.2.0**
>
> Wondertone is still in Development. The API and colour output are not yet stable.
>
> - **Hex output shifts.** `Tone.Hex()` and all render output will produce
>   slightly different values than v0.1.x. The WonderMath pipeline (corrected
>   hue, perceptual chroma, energy glow) changes what every tone renders to.
>   Update any hardcoded hex values or golden test fixtures.
> - **`Temperature()` behaviour changed.** The result is now computed from a
>   continuous formula that factors in chroma and lightness, not just hue ranges.
>   A near-achromatic tone that previously returned `"neutral"` may now return
>   `"warm"` or `"cool"` depending on its hue angle.
> - **`EffectiveC()` is no longer linear.** It previously returned `C × Energy`.
>   It now returns `C × Energy^γ` (γ≈0.7, Stevens' power law). The same
>   Energy value produces a different chroma output.
> - **Vibrancy → chroma mapping changed.** `PerceivedChroma` now applies a
>   power-law exponent (V^α, α=0.9) and a per-hue weight k(H). Blues get
>   slightly more chroma at the same Vibrancy value; yellows get slightly less.

```go
import tone "github.com/leraniode/wondertone/core"

spark := tone.New(
    tone.Light(75),
    tone.Vibrancy(80),
    tone.Hue(30),
    tone.Energy(0.9),
    tone.Named("Primary Spark"),
    tone.Moody("vibrant"),
)

fmt.Println(spark.Hex()) // #f5a04a (gamut-safe, always)
```

---

## What is wondertone?

Wondertone speaks a new vocabulary. Instead of OKLCH's `L`, `C`, `H` — which require color-science knowledge — wondertone gives you:

| Term         | Range | Meaning                                                          |
| ------------ | ----- | ---------------------------------------------------------------- |
| **Light**    | 0–100 | Perceptual lightness. 0=black, 100=white.                        |
| **Vibrancy** | 0–100 | Colorfulness as % of gamut max. 0=grey, 100=most vivid possible. |
| **Hue**      | 0–360 | Color angle. 0=red, 120=green, 240=blue.                         |
| **Energy**   | 0–1   | Aliveness multiplier. 1=full, 0=muted.                           |

Under the hood: OKLCH, OKLab mixing, binary-search gamut safety, perceptual tone scales. You never have to know any of that.

---

## Install

```bash
go get github.com/leraniode/wondertone
```

Import alias convention:

```go
import (
    tone    "github.com/leraniode/wondertone/core"
    palette "github.com/leraniode/wondertone/palette"
    builtin "github.com/leraniode/wondertone/palette/builtin"
    colour  "github.com/leraniode/wondertone/colour"
    render  "github.com/leraniode/wondertone/render"
    "github.com/leraniode/wondertone/wtone"
)
```

---

## Core — Tones

### Creating tones

```go
// Wondertone vocabulary (recommended)
t := tone.New(
    tone.Light(68),
    tone.Vibrancy(72),
    tone.Hue(142),
    tone.Energy(0.95),
    tone.Named("Unix"),
    tone.Moody("focused"),
)

// From hex
t, err := tone.FromHex("#e94560")
t    := tone.MustFromHex("#e94560")

// From raw OKLCH (power users)
t := tone.FromOKLCH(0.68, 0.18, 142.0)

// From .wtone string
t, err := tone.FromOKLCHString("0.68 0.18 142")
```

### Reading tones

```go
t.Light()      // float64 [0–100]
t.Vibrancy()   // float64 [0–100]
t.Hue()        // float64 [0–360)
t.Energy()     // float64 [0–1]
t.Name()       // string
t.Mood()       // string
t.Hex()        // "#rrggbb"
t.RGB()        // uint8, uint8, uint8
t.RGBFloat()   // float64, float64, float64
t.OKLCH()      // l, c, h float64  (raw values — power users)
t.IsLight()    // bool
t.IsDark()     // bool
t.Temperature() // "warm" | "cool" | "neutral"
t.EffectiveC()  // C × Energy — what actually renders
```

### Manipulating tones

All methods return a new Tone — originals are never modified.

```go
t.WithLight(80)
t.WithVibrancy(50)
t.WithHue(200)
t.WithEnergy(0.5)
t.WithName("New Name")
t.WithMood("serene")
t.WithAlpha(0.8)

t.Lighten(15)      // Light += 15
t.Darken(15)       // Light -= 15
t.Saturate(20)     // Vibrancy += 20
t.Desaturate(20)   // Vibrancy -= 20
t.Rotate(90)       // Hue += 90
t.Complement()     // Hue += 180
```

### Energy — the aliveness dial

Energy scales chroma at render time without altering the stored tone:

```go
vivid := colour.Bloom                 // Energy=1.0
quiet := colour.Bloom.WithEnergy(0.3) // same color, quieter

vivid.Hex() == quiet.Hex()            // false — different render
vivid.OKLCH()                         // same L, C, H
quiet.OKLCH()                         // same L, C, H — stored truth unchanged
```

### Tone Scale

12-step perceptual ladder from any tone. Hue never drifts. Every step guaranteed in-gamut.

```go
scale := colour.Unix.Scale()

scale.Background()       // step 1  — page background
scale.SubtleBackground() // step 2  — stripes, alternating rows
scale.ElementBackground()// step 3  — cards, inputs
scale.HoveredBackground()// step 4  — hover state
scale.ActiveBackground() // step 5  — selected
scale.SubtleBorder()     // step 6  — separators
scale.Border()           // step 7  — input borders
scale.StrongBorder()     // step 8  — focus rings
scale.Solid()            // step 9  — buttons, badges
scale.HoveredSolid()     // step 10 — button hover
scale.Text()             // step 11 — body text
scale.HighContrastText() // step 12 — headings
scale.Step(n)            // step n [1–12]
scale.All()              // []Tone, all 12
```

### Mixing

```go
// Mix two tones in OKLab space (no grey midpoint artifacts)
mid := tone.Mix(a, b, 0.5)

// Gradient — n steps from a to b
steps, err := tone.Gradient(a, b, 10)

// Weighted blend of multiple tones
result, err := tone.Blend([]tone.Tone{a, b, c}, []float64{1, 2, 1})

// Harmony
tones, err := tone.Harmonize(base, "complement")   // 2 tones
tones, err := tone.Harmonize(base, "triadic")      // 3 tones
tones, err := tone.Harmonize(base, "analogous")    // 3 tones
tones, err := tone.Harmonize(base, "split")        // 3 tones
tones, err := tone.Harmonize(base, "tetradic")     // 4 tones

// Temperature shift
warmer, err := tone.Shift(t, "warmer", 0.3)
cooler, err  := tone.Shift(t, "cooler", 0.5)
```

### Accessibility

```go
ratio := fg.ContrastWith(bg)          // float64 [1–21]
fg.PassesAA(bg)                       // bool — 4.5:1
fg.PassesAAA(bg)                      // bool — 7.0:1
fixed := fg.EnsureContrast(bg, "AA")  // adjusts lightness only
```

---

## WonderMath — Perceptual Colour Science

v0.2 introduces WonderMath: a layer of perceptual corrections and new
dimensions built above OKLCH. Every tone now passes through this pipeline
at render time.

### Corrected Hue — blue drift fix

OKLCH blues at high chroma drift toward purple. WonderMath applies a
chroma-weighted Gaussian correction. Grey tones are never affected.

```go
// Before v0.2: vivid blue at H=250 could look slightly purple
// After v0.2: corrected back to true blue automatically
blue := tone.New(tone.Light(50), tone.Vibrancy(95), tone.Hue(250))
fmt.Println(blue.Hex()) // perceptually accurate blue
```

### Perceptual Vibrancy — equal vividness across hues

Yellow naturally appears more vivid than blue at the same raw chroma.
WonderMath applies a per-hue weight `k(H)` so `Vibrancy(80)` feels
equally vivid at every hue:

```go
yellow := tone.New(tone.Light(70), tone.Vibrancy(80), tone.Hue(60))
blue   := tone.New(tone.Light(70), tone.Vibrancy(80), tone.Hue(240))
// These now feel equally vivid — yellow reduced, blue boosted
```

### Energy — now perceptually linear

Energy uses Stevens' power law (`E^γ`, γ=0.7).
`Energy=0.5` now genuinely feels half as alive, not 60%.

```go
full  := colour.Bloom                  // Energy=1.0 — full aliveness
half  := colour.Bloom.WithEnergy(0.5)  // feels half as alive
quiet := colour.Bloom.WithEnergy(0.2)  // barely a whisper

// Whole palette — quieten everything at once
hushed := builtin.Midnight().WithEnergy(0.4)
```

### Temperature — continuous warm↔cool

`Temperature()` now returns a label driven by a real formula that factors
in hue, chroma, and lightness. `TemperatureScalar()` gives you the raw value:

```go
ember   := tone.New(tone.Light(60), tone.Vibrancy(70), tone.Hue(25))
glacier := tone.New(tone.Light(60), tone.Vibrancy(70), tone.Hue(196))

ember.Temperature()        // "warm"
ember.TemperatureScalar()  // 0.74

glacier.Temperature()        // "cool"
glacier.TemperatureScalar()  // -0.61
```

### Mood — derived from colour math

Mood is now computed from Valence and Arousal — mathematical properties
of the tone. `DerivedMoodValue()` gives the derived mood. `Mood()` returns
your manual override (or the derived value if none is set).

```go
vivid := tone.New(tone.Light(65), tone.Vibrancy(95), tone.Hue(40), tone.Energy(1.0))
muted := tone.New(tone.Light(35), tone.Vibrancy(15), tone.Hue(220), tone.Energy(0.2))

vivid.DerivedMoodValue() // "playful"
vivid.ValenceValue()     // 0.79   — positive, warm
vivid.ArousalValue()     // 0.97   — activated

muted.DerivedMoodValue() // "deep"
muted.ValenceValue()     // 0.18   — low positive
muted.ArousalValue()     // 0.14   — calm

// Manual override still works
named := vivid.WithMood("sunrise") // custom label, math still runs
named.Mood()             // "sunrise"
named.DerivedMoodValue() // "playful" — math unchanged
```

## Colour — Leraniode's named tones

```go
import "github.com/leraniode/wondertone/colour"

colour.Unix      // terminal green, focused
colour.Starlight // deep indigo, mystical
colour.Ember     // warm amber, comfortable
colour.Glacier   // cool cyan, serene
colour.Crimson   // signal red, urgent
colour.Void      // near-black, deep
colour.Dawn      // soft pink-orange, hopeful
colour.Bloom     // vivid magenta, joyful
colour.Slate     // neutral blue-grey, calm
colour.Signal    // success green, earned
colour.Ink       // near-black text, intentional
colour.Paper     // warm off-white, easy

colour.All()     // []Tone — all 12
```

> [!IMPORTANT]
> The colour package will be removed and will be used
> for an entirely different purpose on the next release

---

## Palette

```go
p, err := palette.New("My Palette").
    Description("A beautiful collection").
    Mood("vibrant").
    Author("leraniode").
    Add(colour.Unix).
    Add(colour.Starlight).
    Add(colour.Crimson).
    Build()

p.Get("Unix")          // (Tone, bool)
p.MustGet("Unix")      // Tone — panics if missing
p.At(0)                // Tone — by index
p.Has("Unix")          // bool
p.All()                // []Tone
p.Len()                // int

// Modify
fork    := p.Fork("My Fork").Add(newTone).Build()
ext, _  := p.Extend("Extended", extraTone)
rep, _  := p.Replace("Unix", differentTone)
quiet   := p.WithEnergy(0.4) // quieten every tone at once

// Validate
report := p.Validate()
fmt.Println(report) // ✓ My Palette — all checks passed
```

### Harmony generators

```go
palette.Complementary(base)            // 2 tones
palette.Triadic(base)                  // 3 tones
palette.Analogous(base, 5, 30)         // 5 tones, 30° spread
palette.SplitComplementary(base, 30)   // 3 tones
palette.Tetradic(base)                 // 4 tones
palette.Monochrome(base, 8)            // 8 tones, same hue
palette.Rainbow(base, 12)              // 12 tones, full wheel
```

### Contrast tools

```go
ratio, err := palette.ContrastPair(p, "Midnight Text", "Midnight Base")
fixed, err  := palette.EnsurePairContrast(p, "Text", "Background", "AA")
matrix      := palette.ContrastMatrix(p)          // every pair
pairs        := palette.FindReadablePairs(p, "AA") // passing pairs only
```

### Built-in palettes

```go
import "github.com/leraniode/wondertone/palette/builtin"

builtin.Midnight()  // deep dark navy
builtin.Aurora()    // bright airy light
builtin.Ember()     // warm amber dark
builtin.Glacier()   // cool icy dark
builtin.Rosewood()  // rich rose dark

builtin.All()       // []*Palette — all five
builtin.Names()     // []string
```

> [!IMPORTANT]
> The palette/builtin package will be removed entirely
> with no ready-made palette available on the next release

---

## Render — terminal output

```go
import render "github.com/leraniode/wondertone/render"

profile := render.Detect() // TrueColor | ANSI256 | ANSI16 | NoColor

render.FG(t, profile)                    // foreground escape sequence
render.BG(t, profile)                    // background escape sequence
render.Colorize(t, profile, "text")      // wrap text with FG + reset
render.ColorizeOnBG(fg, bg, profile, "text")
render.Swatch(t, profile, 2)             // colored block preview

// Lipgloss integration
color := render.LipglossColor(t, profile) // string for lipgloss.Color()
```

Downsampling for ANSI256/ANSI16 uses **OKLab ΔE** nearest-neighbor — perceptual accuracy, not RGB distance.

---

## .wtone files

The `.wtone` format is wondertone's native, toml-based palette file.

```toml
name        = "Leraniode Starlight"
description = "bright starful leraniode"
mood        = "joyful"
version     = "1.0.0"
author      = "leraniode"

[[colors]]
name   = "Primary Spark"
l      = 0.75
c      = 0.15
h      = 30.0
energy = 0.85
mood   = "vibrant"

[[colors]]
name  = "Accent Glow"
oklch = "0.60 0.20 200"   # shorthand — L C H
energy = 0.70
```

```go
import "github.com/leraniode/wondertone/wtone"

// Load
p, err := wtone.LoadWTone("my-palette.wtone")

// Parse from embedded bytes (go:embed)
//go:embed my-palette.wtone
var paletteFile []byte
p, err := wtone.ParseWTone(paletteFile)

// Save
err := wtone.SaveWTone("output.wtone", p)

// Marshal to bytes
data, err := wtone.MarshalWTone(p)
```

---

## Adapters

Wondertone has Adapters for two libraries:

- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [go-colourful](https://github.com/lucasb-eyer/go-colourful)

The adapters convert between wondertone's OKLCH and the target library's color model, so you can use wondertone's colour intelligence in your existing projects without a full rewrite.

They are located in the `adapters/` package, as a seperate module

### Installing the adapters

#### Lipgloss adapter

```bash
go get github.com/leraniode/wondertone/adapters/lipgloss
```

#### go-colourful adapter

```bash
go get github.com/leraniode/wondertone/adapters/colourful
```

Each adapter has its respective dependency and `wondertone` as a dependency.

### Usage

#### Lipgloss adapter

Import the adapter as wtlip to avoid conflict with lipgloss package API:

```go
package main


import (
   tone  "github.com/leraniode/wondertone/core"
   wtlip "github.com/leraniode/wondertone/adapters/lipgloss"
)

wtlip.FG(colour.Unix).Bold(true).Render("hello")

wtlip.Style(colour.Unix).
   Background(colour.Void).
   Padding(0, 1).
   Render("hello")

wtlip.PaletteStyles(builtin.Midnight()) // map[name]lipgloss.Style
```

#### go-colourful adapter

Import the adapter as wcolourful to avoid conflict with go-colourful package API:

```go
package main

import (
    tone      "github.com/leraniode/wondertone/core"
    wcolorful "github.com/leraniode/wondertone/adapters/colorful"
)

// wondertone → go-colorful
cf := wcolorful.ToColorful(myTone)
lab, _ := cf.Lab()

// go-colorful → wondertone
t := wcolorful.FromColorful(cf)
fmt.Println(t.Hex())
```

## License

MIT — Leraniode

---

Part of [Leraniode](https://github.com/leraniode) – Building Tools that feel alive 🌱.

<p align="left">
    <a href="https://github.com/leraniode">
        <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/footer.png" width="600" />
    </a>
</p>
