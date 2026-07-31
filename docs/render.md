# render

```go
import "github.com/leraniode/wondertone/render"
```

Terminal colour output. Detects capability from environment variables and
outputs the best possible ANSI escape sequence. Downsampling to ANSI256
and ANSI16 is perceptual — OKLab ΔE nearest-neighbour, not RGB distance.

---

## Profile Detection

```go
profile := render.Detect()
```

Reads environment variables in order:

| Variable | Condition | Profile |
|---|---|---|
| `NO_COLOR` | set (any value) | `NoColor` |
| `TERM` | `"dumb"` | `NoColor` |
| `COLORTERM` | `"truecolor"` or `"24bit"` | `TrueColor` |
| `TERM_PROGRAM` | iTerm.app, Hyper, vscode, WezTerm, ghostty | `TrueColor` |
| `TERM` | ends with `-256color` | `ANSI256` |
| `TERM` | contains `color` or starts with `xterm` | `TrueColor` |
| `CI` | set | `ANSI256` |
| — | fallback | `ANSI16` |

```go
// Force a specific profile (useful in tests)
profile := render.Force(render.TrueColor)
profile := render.Force(render.NoColor)
```

Profiles: `render.TrueColor`, `render.ANSI256`, `render.ANSI16`, `render.NoColor`.

---

## Output Functions

### FG / BG

Raw ANSI escape sequences — use when building your own formatted strings.

```go
fmt.Print(render.FG(t, profile), "coloured text", render.Reset)
fmt.Print(render.BG(t, profile), "coloured background", render.Reset)
```

### Colorize

Wraps text with FG colour and a reset. The reset is included.

```go
s := render.Colorize(t, profile, "hello")
fmt.Println(s)
```

### ColorizeOnBG

Both foreground and background coloured. Reset included.

```go
s := render.ColorizeOnBG(fg, bg, profile, "label")
```

### Swatch

A solid block of colour, `width` characters wide. Useful for palette previews.

```go
fmt.Println(render.Swatch(t, profile, 3))  // ███
```

---

## ANSI Constants

```go
render.Reset     // "\x1b[0m"
render.Bold      // "\x1b[1m"
render.Dim       // "\x1b[2m"
render.Italic    // "\x1b[3m"
render.Underline // "\x1b[4m"
```

---

## lipgloss

For lipgloss users — returns a colour value compatible with
`lipgloss.NewStyle().Foreground()` and `.Background()`.

```go
color := render.LipglossColor(t, profile)
// returns an ANSI colour string matched to the profile

hex := render.LipglossColorHex(t)
// always returns the full hex — use when lipgloss handles profile itself
```
