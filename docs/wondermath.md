# WonderMath

WonderMath is the mathematical layer of Wondertone. It is the implementation
of WonderSpace formulas and logic — what transforms wondertone from a plain
OKLCH library into a perceptual colour intelligence system.

WonderMath lives in `space/`. All functions are pure — float64 in, float64
out, no state, no side effects. The constants are named and tunable without
touching logic.

---

## What WonderMath Does

Every Tone passes through WonderMath at render time. The stored vocabulary
values (Light, Vibrancy, Hue, Energy, Mood) never change. Only the rendered
output is corrected and enriched.

---

## The Pipeline

Applied in this order during `Hex()` and all colour output:

### 1. Corrected Hue

OKLCH has a known residual non-linearity in the blue region (~H=220–280):
vivid blues drift toward purple. WonderMath applies a chroma-weighted
Gaussian correction that nudges them back. Greys (C≈0) are never touched.

```
H' = H + A × exp(-(H - H₀)² / (2w²)) × (C / C_max)
H₀ = 250°,  w = 30°,  A = -3°
```

### 2. Perceived Chroma

At the same raw Vibrancy, yellow looks more vivid than blue — not a
formula error but how the eye works. WonderMath applies a per-hue
weight `k(H)` and power-law shaping so that Vibrancy=80 feels equally
vivid at every hue.

```
C = C_max(L, H) × V^α × k(H)
α = 0.9
k(H) examples: Yellow (H≈60) = 0.85,  Blue (H≈240) = 1.10
```

### 3. Effective Chroma — Energy (Stevens' power law)

Linear Energy scaling feels perceptually wrong. Energy=0.5 on a linear
scale still feels ~65% alive. Stevens' power law corrects this:

```
C_effective = C_base × E^γ    γ = 0.7
```

Energy=0.5 now genuinely feels half as alive.

### 4. Effective Lightness — Energy glow

High-Energy tones receive a subtle lightness boost. Stored Light is never
modified — this is expression only:

```
L_effective = L + λ × (E^γ - 1)    λ = 0.04
```

### 5. Temperature Value

Not a hue-range lookup. A continuous formula factoring hue, chroma, and
lightness simultaneously. Range [-1, +1]:

```
T = clamp(w_h × cos((H - 50°) × π/180) + w_c × (C/C_max) + w_l × (L - 0.5), -1, 1)
w_h = 0.7,  w_c = 0.2,  w_l = 0.1
```

### 6. Valence + Arousal + Derived Mood

Mood is not a tag. It is derived from two real properties of the colour:

```
valence = clamp(a₁×T + a₂×L + a₃×S, -1, 1)
arousal = clamp(b₁×S + b₂×E + b₃×T, -1, 1)
```

Nine mood regions: vivid, playful, urgent, warm, mystical, focused, airy,
deep, calm.

**Coherence** — the third mood axis — is defined by WonderSpace research.
Formula derivation is the current open research target.

---

## Tuning

All constants sit at the top of `space/wondermath.go`, clearly named.
Adjust a constant — the entire pipeline shifts. Logic and values are
intentionally separated.
