# Changelog

All notable changes to wondertone are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

---

## [0.1.0] — 2026-03-04

Initial release. The full wondertone Go library.

### Added

**`core/` — the Tone type**
- `Tone` struct with Light, Vibrancy, Hue, Energy vocabulary
- Constructors: `New()` (option pattern), `FromHex()`, `FromOKLCH()`, `FromOKLCHString()`
- Full immutable manipulation API: `WithLight`, `Lighten`, `Darken`, `Rotate`, `Complement`, `Saturate`, `Desaturate`, `WithEnergy`, etc.
- OKLCH pipeline, public-domain math — zero external dependencies in core
- Iterative gamut mapping — hue never drifts
- `ToneScale` — 12-step perceptual scale with semantic accessors
- OKLab mixing: `Mix`, `Gradient`, `Blend`, `Harmonize`, `Shift`
- WCAG accessibility: `ContrastWith`, `PassesAA`, `PassesAAA`, `EnsureContrast`
- Intelligence: `IsLight`, `IsDark`, `Temperature`

**`colour/` — Leraniode named tones**
- Unix, Starlight, Ember, Glacier, Crimson, Void, Dawn, Bloom, Slate, Signal, Ink, Paper
- `All()` collection helper

**`palette/` — Palette management**
- `Palette` with builder pattern: `New().Add().Build()`
- Immutable operations: `Fork`, `Extend`, `Replace`, `WithEnergy`
- `Validate()` with `ValidationReport`
- Harmony generators: `Complementary`, `Triadic`, `Analogous`, `SplitComplementary`, `Tetradic`, `Monochrome`, `Rainbow`
- Contrast tools: `ContrastPair`, `EnsurePairContrast`, `ContrastMatrix`, `FindReadablePairs`

**`palette/builtin/` — Built-in palettes**
- Midnight (dark navy), Aurora (light), Ember (warm dark), Glacier (cool dark), Rosewood (rose dark)
- `All()` and `Names()` helpers

**`render/` — Terminal output**
- Profile detection: `Detect()` reads `NO_COLOR`, `COLORTERM`, `TERM`, `TERM_PROGRAM`
- `FG()`, `BG()`, `Colorize()`, `ColorizeOnBG()`, `Swatch()`
- Perceptual ANSI256/16 downsampling via OKLab ΔE nearest-neighbor
- `LipglossColor()` adapter for charmbracelet/lipgloss

**`wtone/` — File format**
- Load/save `.wtone` (TOML-based palette files)
- `[[colors]]` array format — order preserved
- Shorthand `oklch = "L C H"` and explicit `l/c/h` fields
- Mood inheritance from palette to tones

**`internal/testutil/`**
- Zero-dependency test helpers — no gopkg.in/yaml.v3 transitive issues

**`adapters/`**
- `lipgloss/` — Adapter for wondertone interop with charmbracelet/lipgloss
- `go-colorful/` — Adapter for interop with github.com/lucasb-eyer/go-colorful
- Released as seperate modules each with `go.mod` and `go.sum` to avoid transitive dependency issues

**CI**
- GitHub Actions: test matrix (ubuntu/macos/windows × go1.22/1.23)
- golangci-lint
- Release workflow on `v*.*.*` tags

[0.1.0]: https://github.com/leraniode/wondertone/releases/tag/v0.1.0
