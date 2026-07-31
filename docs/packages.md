# Package Reference

---

## `tone/`

```go
import "github.com/leraniode/wondertone/tone"
```

The Tone type, constructors, options, and the full manipulation API.
The base of the dependency graph — everything else imports tone, tone
imports nothing from within wondertone except space.

---

## `space/`

```go
import "github.com/leraniode/wondertone/space"
```

WonderMath — the mathematical layer of Wondertone. Pure perceptual colour
science: hue correction, perceived chroma, energy scaling, temperature,
valence, arousal, mood derivation. All functions are pure float64. No Tone
type imported. tone imports space, not the other way around.

---

## `mix/`

```go
import "github.com/leraniode/wondertone/mix"
```

OKLab-space colour mixing, gradients, weighted blending, and harmonic
schemes. All operations take and return Tones.

---

## `palette/`

```go
import "github.com/leraniode/wondertone/palette"
```

The Palette type and builder. An ordered, named collection of Tones with
immutable operations, validation, and energy/mood application across the set.

---

## `harmony/`

```go
import "github.com/leraniode/wondertone/harmony"
```

Palette-returning harmony generators based on hue relationships in OKLCH:
Analogous, Complementary, Triadic, SplitComplementary, Tetradic, Monochrome.

---

## `contrast/`

```go
import "github.com/leraniode/wondertone/contrast"
```

APCA contrast ratios, AA/AAA accessibility validation, contrast matrices,
and readable pair finding.

---

## `render/`

```go
import "github.com/leraniode/wondertone/render"
```

Terminal colour output. Profile detection from environment variables
(`NO_COLOR`, `COLORTERM`, `TERM`, `TERM_PROGRAM`). FG/BG/Swatch rendering.
Perceptual ANSI256/16 downsampling. `LipglossColor()` for lipgloss users.
