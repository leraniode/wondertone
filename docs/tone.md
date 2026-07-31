# tone

```go
import "github.com/leraniode/wondertone/tone"
```

The Tone type — the fundamental unit of wondertone. A Tone is an immutable,
named colour with a human vocabulary. Every other package works with Tones.

---

## Creating Tones

### From the vocabulary

```go
t := tone.New(
    tone.Light(75),      // lightness [0–100]
    tone.Vibrancy(80),   // colorfulness [0–100]
    tone.Hue(30),        // colour angle [0–360]
    tone.Energy(0.9),    // aliveness [0–1]
    tone.Named("Ember"),
    tone.Moody("warm"),
    tone.Alpha(1.0),     // opacity [0–1], default 1.0
)
```

Defaults when options are omitted: Light=50, Vibrancy=100, Hue=0, Energy=1.0.

### From hex

```go
t, err := tone.FromHex("#e88566")   // #rgb, #rrggbb, #rrggbbaa
t       := tone.MustFromHex("#e88566") // panics on error — safe for constants
```

### From OKLCH

```go
t := tone.FromOKLCH(0.65, 0.18, 30) // L [0–1], C [0–~0.37], H [0–360)

t, err := tone.FromOKLCHString("0.65 0.18 30")
t, err := tone.FromOKLCHString("0.65 0.18 30 / 0.9") // with alpha
```

---

## Reading a Tone

### Vocabulary

```go
t.Light()      // float64 — [0–100]
t.Vibrancy()   // float64 — [0–100]
t.Hue()        // float64 — [0–360)
t.Energy()     // float64 — [0–1]
t.Name()       // string
t.Mood()       // string — manual tag
t.AlphaValue() // float64 — [0–1]
```

### Colour output

```go
t.Hex()        // "#e88566" — gamut-safe CSS hex
t.RGB()        // (r, g, b uint8) — [0–255]
t.RGBFloat()   // (r, g, b float64) — [0–1]
t.String()     // same as Hex()
```

### OKLCH

```go
l, c, h := t.OKLCH()     // raw internal values
s        := t.OKLCHString() // "0.650000 0.180000 30.000000"
t.RawL()  // float64 — L [0–1]
t.RawC()  // float64 — C chroma
```

### WonderMath derived

```go
t.Temperature()        // "warm" / "cool" / "neutral"
t.TemperatureScalar()  // float64 — continuous [-1, +1]
t.DerivedMoodValue()   // "playful" / "calm" / "deep" / etc.
t.ValenceValue()       // float64 — emotional positivity [-1, +1]
t.ArousalValue()       // float64 — activation level [-1, +1]
t.EffectiveC()         // float64 — chroma after Stevens' power law
```

---

## Manipulating Tones

All methods return a new Tone. The original is never modified.

### Direct setters

```go
t.WithLight(60)
t.WithVibrancy(90)
t.WithHue(200)
t.WithEnergy(0.5)
t.WithName("New Name")
t.WithMood("serene")
t.WithAlpha(0.8)
```

### Relative adjustments

```go
t.Lighten(10)       // Light + 10
t.Darken(15)        // Light - 15
t.Saturate(20)      // Vibrancy + 20
t.Desaturate(30)    // Vibrancy - 30
t.Rotate(120)       // Hue + 120°, wraps at 360
t.Complement()      // Rotate(180), clears name and mood
```

### Inspection

```go
t.IsLight()         // Light > 50
t.IsDark()          // Light <= 50
t.Equal(other)      // perceptual equality within floating-point tolerance
```

---

## Accessibility

```go
t.Luminance()           // WCAG 2.1 relative luminance [0–1]
t.ContrastWith(other)   // WCAG 2.1 contrast ratio [1–21]
t.PassesAA(bg)          // contrast >= 4.5:1
t.PassesAAA(bg)         // contrast >= 7.0:1
t.EnsureContrast(bg, "AA")   // returns adjusted Tone that passes AA
t.EnsureContrast(bg, "AAA")  // returns adjusted Tone that passes AAA
```

`EnsureContrast` adjusts lightness only. Hue, vibrancy, and energy are preserved.

---

## Tone Scale

A 12-step perceptual ladder from a single base Tone. Same hue across all steps.
Every step is gamut-safe.

```go
scale := t.Scale()         // ToneScale — all 12 steps
step  := t.Step(9)         // single step [1–12], 1=lightest, 12=darkest

// Numeric
scale.Step(1)   // lightest
scale.Step(12)  // darkest

// Semantic (UI roles)
scale.Background()        // step 1  — page background
scale.SubtleBackground()  // step 2  — alternating rows, stripes
scale.ElementBackground() // step 3  — cards, inputs
scale.HoveredBackground() // step 4  — hover state
scale.ActiveBackground()  // step 5  — selected state
scale.SubtleBorder()      // step 6  — separators
scale.Border()            // step 7  — input borders
scale.StrongBorder()      // step 8  — focus rings
scale.Solid()             // step 9  — buttons, badges — base tone lives here
scale.HoveredSolid()      // step 10 — button hover
scale.Text()              // step 11 — body text
scale.HighContrastText()  // step 12 — headings, emphasis

scale.All()               // []Tone, lightest to darkest
```
