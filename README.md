<p align="left">
    <a href="https://github.com/leraniode/wondertone">
        <img src="https://raw.githubusercontent.com/leraniode/.github/main/images/wondertone.svg"/>
    </a>
</p>

# Wondertone 🎨

A perceptual color intelligence library for Go — OKLCH under the hood, a human vocabulary on the surface.

[![CI](https://github.com/leraniode/wondertone/actions/workflows/ci.yml/badge.svg)](https://github.com/leraniode/wondertone/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/leraniode/wondertone.svg)](https://pkg.go.dev/github.com/leraniode/wondertone)
[![License](https://img.shields.io/github/license/leraniode/wondertone)](https://github.com/leraniode/wondertone/blob/main/LICENSE)
[![Version](https://img.shields.io/github/v/tag/leraniode/wondertone)](https://github.com/leraniode/wondertone/releases)
[![Go Modules](https://img.shields.io/github/go-mod/go-version/leraniode/wondertone)](https://github.com/leraniode/wondertone/blob/main/go.mod)

> [!IMPORTANT]
> **Pre-v1 software. No API or output stability is guaranteed.**
> Minor versions may change hex output, method behaviour, or import paths
> without deprecation notices.

```go
import "github.com/leraniode/wondertone/tone"

spark := tone.New(
    tone.Light(75),
    tone.Vibrancy(80),
    tone.Hue(30),
    tone.Energy(0.9),
)

fmt.Println(spark.Hex())               // gamut-safe, always
fmt.Println(spark.TemperatureScalar()) // +0.74
fmt.Println(spark.DerivedMoodValue())  // "playful"
```

---

## What is Wondertone?

Most colour libraries are coordinate converters. Wondertone has a model.

It knows what _alive_ means, what _warm_ means, what emotional character
a colour carries — derived from colour science, not guesswork.

| Term            | Range  | Meaning                                                     |
| --------------- | ------ | ----------------------------------------------------------- |
| **Light**       | 0–100  | Perceptual lightness                                        |
| **Vibrancy**    | 0–100  | Colorfulness as % of gamut max                              |
| **Hue**         | 0–360  | Colour angle                                                |
| **Energy**      | 0–1    | Aliveness — how much a colour presses into its environment  |
| **Temperature** | -1–+1  | Continuous warm↔cool, derived from hue + chroma + lightness |
| **Mood**        | string | Emotional character, derived from Valence + Arousal math    |

Under the hood: OKLCH, WonderMath perceptual corrections, OKLab mixing,
gamut-safe rendering. You never have to know any of that.

---

## Install

```bash
go get github.com/leraniode/wondertone
```

---

## Packages

```
wondertone/tone      Tone type, constructors, manipulation API
wondertone/space     WonderMath — pure perceptual colour science
wondertone/mix       OKLab mixing, gradients, blending
wondertone/palette   Palette type and builder
wondertone/harmony   Harmony generators (triadic, analogous, etc.)
wondertone/contrast  APCA contrast, accessibility, readable pairs
wondertone/render    Terminal output, profile detection
```

---

## Tones

### Creating tones

```go
// From the wondertone vocabulary
t := tone.New(
    tone.Light(68),
    tone.Vibrancy(72),
    tone.Hue(142),
    tone.Energy(0.95),
    tone.Named("Forest"),
)

// From an existing hex — full WonderMath analysis applied
t, err := tone.FromHex("#42a939")

// From raw OKLCH
t := tone.FromOKLCH(0.65, 0.18, 142)
```

### Reading tones

```go
t.Hex()                // "#42a939"
t.Light()              // 68.0
t.Vibrancy()           // 72.0
t.Hue()                // 142.0
t.Energy()             // 0.95
t.Temperature()        // "warm" / "cool" / "neutral"
t.TemperatureScalar()  // continuous [-1, +1]
t.DerivedMoodValue()   // "focused", "playful", "deep", etc.
t.ValenceValue()       // [-1, +1]
t.ArousalValue()       // [-1, +1]
```

### Manipulating tones

Every method returns a new Tone. Tones are immutable.

```go
t.Lighten(10)
t.Darken(15)
t.Rotate(120)
t.Saturate(20)
t.Desaturate(30)
t.WithEnergy(0.3)
t.WithMood("focused")
```

### Energy

```go
// Same tone, different aliveness
full  := t.WithEnergy(1.0)
half  := t.WithEnergy(0.5)  // genuinely half as alive (Stevens' power law)
quiet := t.WithEnergy(0.2)
```

---

## Mixing

```go
import "github.com/leraniode/wondertone/mix"

mid   := mix.Mix(a, b, 0.5)
steps, _ := mix.Gradient(a, b, 7)
blend, _ := mix.Blend([]tone.Tone{a, b, c}, []float64{0.5, 0.3, 0.2})
```

---

## Palette

```go
import "github.com/leraniode/wondertone/palette"

p, err := palette.New("Brand").
    Author("leraniode").
    Mood("focused").
    Add(primary, secondary, accent).
    Build()

quieter  := p.WithEnergy(0.4)
extended := p.Extend(extra)
report   := p.Validate()
```

---

## Harmony

```go
import "github.com/leraniode/wondertone/harmony"

comp,   _ := harmony.Complementary(base)
triad,  _ := harmony.Triadic(base)
analog, _ := harmony.Analogous(base, 3, 30)
```

---

## Contrast

```go
import "github.com/leraniode/wondertone/contrast"

ratio := contrast.APCARatio(text, bg)
ok    := contrast.PassesAA(text, bg)
pairs := contrast.FindReadablePairs(p, "AA")
```

---

## Terminal Output

```go
import "github.com/leraniode/wondertone/render"

profile := render.Detect()
fmt.Println(render.FG(t, profile, "hello"))
fmt.Println(render.Swatch(t, profile, 3))

// lipgloss
color := render.LipglossColor(t, profile)
```

---

## License

MIT — Leraniode

---

Part of [Leraniode](https://github.com/leraniode).

<p align="center">
    <br/>
    <a href="https://github.com/leraniode">
        <img src="https://raw.githubusercontent.com/leraniode/.github/main/assets/footer.svg" width="1024" />
    </a>
</p>
