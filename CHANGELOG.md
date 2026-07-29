# Changelog

All notable changes to wondertone are documented here.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning: [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

---

## [0.3.0] — 2026-07-29

A decision-point release. Wondertone is a pure Go colour library — nothing
more, nothing less. Everything that didn't serve that identity is removed.

### Changed

**Package structure — packages expanded into focused packages**

| Old                                  | New                   |
| ------------------------------------ | --------------------- |
| `core/` (Tone type, scale)           | `tone/`               |
| `core/` (WonderMath)                 | `space/`              |
| `core/` (OKLab mixing)               | `mix/`                |
| `palette/` (harmony generators)      | `harmony/`            |
| `palette/` (contrast, accessibility) | `contrast/`           |
| `palette/` (Palette type)            | `palette/`            |
| `render/`                            | `render/` (unchanged) |

Import paths change accordingly:

```go
// before
import tone "github.com/leraniode/wondertone/core"
// after
import "github.com/leraniode/wondertone/tone"
import "github.com/leraniode/wondertone/space"
import "github.com/leraniode/wondertone/mix"
```

**Raw OKLCH accessors added to Tone**

`RawL()`, `RawC()`, `RawH()` expose the internal OKLCH coordinates for
downstream packages and power users who need to work below the vocabulary level.

### Removed

- `adapters/` — gone entirely. `render/lipgloss.go` has two lightweight functions, no new dep introduced.
- `wtone/` — removed. Wondertone is pure Go. but can still be found at `leraniode/x/wtone`.
- `colour/`, `palette/builtin/` - planned for separate repo
- `internal/testutil/` — moved to `leraniode/x` as `x/testutil`.

### Updated

- `go.mod` — removed `x/wtone` dep, added `x/testutil v0.1.0`.
- CI workflow — removed dead adapter jobs, cleaned up cache paths.

### Breaking

All `core/` import paths changed, some expected APIs have been moved to new packages.
`colour/` and `palette/builtin/` packages no longer exist in this module.

---

## [0.2.0] — 2026-03-06

WonderMath — Wondertone's perceptual colour science layer above OKLCH.

### Added

**`core/wondermath.go` — six perceptual formulas**

- **Corrected Hue** — Gaussian blue-drift correction (H₀=250, w=30°, A=-3°)
- **Perceived Chroma** — V^α power law + k(H) per-hue weight table (α=0.9)
- **Effective Chroma** — Stevens' power law Energy scaling (γ=0.7)
- **Effective Lightness** — subtle glow at high Energy (λ=0.04)
- **Temperature Value** — continuous warm↔cool scalar [-1, +1]
- **Valence + Arousal + DerivedMood** — Mood derived from colour math

**New Tone accessors**

```go
t.TemperatureScalar() float64
t.DerivedMoodValue()  string
t.ValenceValue()      float64
t.ArousalValue()      float64
```

### Breaking

- Hex output shifts — blues corrected, Energy non-linear, high-energy glow
- `Temperature()` behaviour changed — continuous formula, not hue-range lookup
- `EffectiveC()` non-linear — returns `C × E^γ` not `C × E`
- Vibrancy → Chroma mapping changed — V^α + k(H) applied at construction

---

## [0.1.1] — 2026-03-05

### Fixed

`adapters/lipgloss` — `SetColorProfile` sync for non-TTY rendering.
lipgloss's default renderer detected "no TTY" at startup and stripped
all colour. Fixed by calling `lipgloss.SetColorProfile()` in both
`init()` and `SetProfile()`.

---

## [0.1.0] — 2026-03-04

Initial release.

### Added

- `Tone` type with Light/Vibrancy/Hue/Energy vocabulary
- Constructors: `New()`, `FromHex()`, `FromOKLCH()`, `FromOKLCHString()`
- Full immutable manipulation API
- OKLCH pipeline, iterative gamut safety (hue never drifts)
- `ToneScale` — 12-step perceptual scale
- OKLab mixing: `Mix`, `Gradient`, `Blend`
- WCAG accessibility: `ContrastWith`, `PassesAA`, `PassesAAA`
- Harmony generators: Complementary, Triadic, Analogous, etc.
- Profile detection, terminal rendering, ANSI256/16 downsampling
- GitHub Actions CI (ubuntu / macos / windows)

---

[Unreleased]: https://github.com/leraniode/wondertone/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/leraniode/wondertone/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/leraniode/wondertone/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/leraniode/wondertone/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/leraniode/wondertone/releases/tag/v0.1.0
