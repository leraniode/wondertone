# mix

```go
import "github.com/leraniode/wondertone/mix"
```

OKLab-space colour mixing, gradients, blending, and harmony schemes.
All mixing happens in OKLab — the perceptually uniform mixing space.
No grey mudpoint. No unexpected hue shifts at the midpoint.

---

## Mix

Blend two Tones at a ratio `t` [0–1]. `t=0` returns `a` exactly.
`t=1` returns `b` exactly. `t=0.5` is the perceptual midpoint.

```go
mid := mix.Mix(a, b, 0.5)
```

---

## Gradient

Produce `n` perceptually uniform steps from `start` to `end`.
Minimum 2 steps. First step is exactly `start`, last is exactly `end`.

```go
steps, err := mix.Gradient(a, b, 7)
// steps[0] == a
// steps[6] == b
// steps[1..5] — perceptual midpoints
```

---

## Blend

Mix any number of Tones with per-tone weights. Weights are normalised
automatically — they do not need to sum to 1.

```go
result, err := mix.Blend(
    []tone.Tone{red, green, blue},
    []float64{0.5, 0.3, 0.2},
)
```

Errors: empty slice, mismatched lengths, negative weights, all-zero weights.

---

## Harmonize

Returns a slice of Tones related by a harmonic scheme. Hue rotates in OKLCH.
Lightness, Vibrancy, and Energy are preserved from the base Tone.

```go
tones, err := mix.Harmonize(base, "complement")   // [base, base+180°]
tones, err := mix.Harmonize(base, "triadic")       // [base, base+120°, base+240°]
tones, err := mix.Harmonize(base, "analogous")     // [base-30°, base, base+30°]
tones, err := mix.Harmonize(base, "split")         // [base, base+150°, base+210°]
tones, err := mix.Harmonize(base, "tetradic")      // [base, base+90°, base+180°, base+270°]
```

`mix.Harmonize` returns raw Tone slices. For Palette-returning generators with
named tones, use the `harmony` package instead.
