# palette

```go
import "github.com/leraniode/wondertone/palette"
```

An ordered, named collection of Tones. Every Tone in a Palette must have
a unique name. Order of insertion is preserved.

---

## Building

```go
p, err := palette.New("Brand").
    Description("Our core palette").
    Author("leraniode").
    Mood("focused").
    Version("1.0.0").
    Add(primary).
    Add(secondary).
    Add(accent).
    Build()
```

`Build()` returns an error if any Tone is unnamed or if names are duplicated.

```go
// Panic variant — safe for package-level variables
p := palette.New("Brand").Add(a).Add(b).MustBuild()
```

---

## Reading

```go
p.Name()         // "Brand"
p.Description()  // string
p.Mood()         // string
p.Author()       // string
p.Version()      // string
p.Len()          // int — number of Tones

p.All()          // []tone.Tone — copy, insertion order

t, ok := p.Get("Primary")    // lookup by name
t      := p.MustGet("Primary") // panics if not found
t      := p.At(0)             // by position (0-based), panics if out of range
ok     := p.Has("Primary")   // existence check
```

---

## Operations

All operations return a new Palette. The original is never modified.

### Fork

Creates a Builder pre-populated with all Tones from this Palette.
Use when you want to modify or extend under a new name.

```go
dark, err := p.Fork("Brand Dark").
    Add(extraDarkTone).
    Build()
```

### Extend

Add Tones to a copy without modifying the original.
Cannot override existing names — use Fork for that.

```go
extended, err := p.Extend("Brand Extended", extraTone, anotherTone)
```

### Replace

Swap a named Tone for a different one.

```go
updated, err := p.Replace("Primary", newPrimary)
```

### WithEnergy

Apply the same Energy to every Tone. The "mood dial" for a whole palette.

```go
hushed := p.WithEnergy(0.3)  // same colours, much quieter
vivid  := p.WithEnergy(1.0)  // full presence
```

---

## Validation

```go
report := p.Validate()
fmt.Println(report)      // "✓ Brand — all checks passed"
report.Passed            // bool
report.Issues            // []string — descriptions of each problem
```

Checks performed:
- At least 2 Tones
- Warning if more than 16 Tones (suggests splitting)
- All Tones in sRGB gamut
- Adjacent Tones differ by at least ΔLight=5 (perceptually distinguishable)
