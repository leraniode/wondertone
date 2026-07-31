# harmony

```go
import "github.com/leraniode/wondertone/harmony"
```

Palette-returning harmony generators based on hue relationships in OKLCH.
All generators preserve Lightness, Vibrancy, and Energy from the base Tone.
Only Hue varies. Every returned Tone is gamut-safe.

---

## Generators

### Complementary

Two Tones — base and its direct opposite (+180°).

```go
p, err := harmony.Complementary(base)
// p.Len() == 2
```

### Triadic

Three Tones evenly spaced 120° apart.

```go
p, err := harmony.Triadic(base)
// hues: base, base+120°, base+240°
```

### Analogous

Adjacent Tones centered on the base hue.

```go
p, err := harmony.Analogous(base, 5, 30)
// 5 tones, 30° apart, centered on base
// hues: base-60°, base-30°, base, base+30°, base+60°
```

Minimum 2 tones. `spreadDeg` is the gap between adjacent tones.

### SplitComplementary

Base plus two Tones flanking its complement.

```go
p, err := harmony.SplitComplementary(base, 30)
// hues: base, base+150°, base+210°
```

`splitDeg` must be in (0, 90). Smaller values = tighter flanking.

### Tetradic

Four Tones evenly spaced 90° apart.

```go
p, err := harmony.Tetradic(base)
// hues: base, base+90°, base+180°, base+270°
```

### Monochrome

`count` Tones from the same hue, spread across the lightness scale.
Lightness runs from 15 to 92 — avoids pure black and pure white.

```go
p, err := harmony.Monochrome(base, 6)
// 6 tones, same hue, evenly distributed from L=15 to L=92
```

Minimum 2 tones.

### Rainbow

`count` Tones evenly distributed around the full hue wheel.
Same lightness and vibrancy as the base. First Tone starts at base hue.

```go
p, err := harmony.Rainbow(base, 12)
// 12 tones, 30° apart, full hue circle
```

Minimum 2 tones.

---

## vs. mix.Harmonize

`mix.Harmonize` returns raw `[]tone.Tone` slices for complement, triadic,
analogous, split, and tetradic schemes. Use it when you want just the Tones
without building a Palette.

`harmony` returns named, ordered Palettes with automatically generated
Tone names. Use it when you want a ready-to-use Palette.
