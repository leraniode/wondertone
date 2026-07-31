# contrast

```go
import "github.com/leraniode/wondertone/contrast"
```

WCAG contrast ratios and accessibility tools operating on Palettes.
All ratio calculations use the WCAG 2.1 relative luminance formula.

---

## Ratio and Level Check

```go
ratio := contrast.APCARatio(text, bg)  // float64 — WCAG contrast ratio [1–21]
ok    := contrast.PassesAA(text, bg)   // ratio >= 4.5:1
ok    := contrast.PassesAAA(text, bg)  // ratio >= 7.0:1
```

---

## Palette-level tools

### ContrastPair

Check the contrast ratio between two named Tones in a Palette.

```go
ratio, err := contrast.ContrastPair(p, "Text", "Background")
```

### EnsurePairContrast

Returns a new Palette where the foreground Tone has been adjusted to meet
the given level against the background. Only lightness is adjusted — hue,
vibrancy, and energy are preserved.

```go
fixed, err := contrast.EnsurePairContrast(p, "Text", "Background", "AA")
fixed, err := contrast.EnsurePairContrast(p, "Text", "Background", "AAA")
```

### ContrastMatrix

Contrast ratio between every Tone pair in the Palette.
Returns a map of `"fg/bg"` → ratio. Self-pairs are excluded.

```go
matrix := contrast.ContrastMatrix(p)
for pair, ratio := range matrix {
    fmt.Printf("%s: %.2f:1\n", pair, ratio)
}
```

### FindReadablePairs

All (fg, bg) Tone pairs that meet the given WCAG level.

```go
pairs := contrast.FindReadablePairs(p, "AA")
pairs := contrast.FindReadablePairs(p, "AAA")

for _, pair := range pairs {
    fmt.Println(pair) // "Text              on Background          14.00:1  ✓ AAA"
    pair.FG           // tone.Tone
    pair.BG           // tone.Tone
    pair.Ratio        // float64
    pair.PassesAA     // bool
    pair.PassesAAA    // bool
}
```

---

## Direct Tone contrast

For per-Tone contrast without importing this package:

```go
t.ContrastWith(other)   // WCAG ratio
t.PassesAA(bg)          // bool
t.PassesAAA(bg)         // bool
t.EnsureContrast(bg, "AA") // adjusted Tone
```
