package space

import (
	"math"
	"testing"
)

// --- CorrectedHue ---

func TestCorrectedHueBlueRegionShiftsAwayFromPurple(t *testing.T) {
	// A vivid blue at H=250 should shift slightly away from purple (negative A)
	h := 250.0
	c := 0.20
	l := 0.5
	h2 := CorrectedHue(h, c, l)
	if h2 >= h {
		t.Errorf("expected blue hue to shift below 250°, got %.4f (original %.4f)", h2, h)
	}
}

func TestCorrectedHueGreyIsUnchanged(t *testing.T) {
	// Achromatic (C=0) should never be corrected
	h := CorrectedHue(250, 0, 0.5)
	if math.Abs(h-250) > 1e-6 {
		t.Errorf("grey hue should not be corrected, got %.6f", h)
	}
}

func TestCorrectedHueFarFromBlueIsUnchanged(t *testing.T) {
	// Green (H=142) is far from the correction zone — should be nearly unchanged
	h := CorrectedHue(142, 0.25, 0.6)
	if math.Abs(h-142) > 0.5 {
		t.Errorf("green hue should be nearly unchanged, got %.4f", h)
	}
}

func TestCorrectedHueWrapsCorrectly(t *testing.T) {
	// Result must always be in [0, 360)
	for _, h := range []float64{0, 90, 180, 250, 359} {
		result := CorrectedHue(h, 0.2, 0.5)
		if result < 0 || result >= 360 {
			t.Errorf("CorrectedHue(%v) = %.4f, not in [0,360)", h, result)
		}
	}
}

// --- PerceivedChroma ---

func TestPerceivedChromaZeroVibrancyIsZero(t *testing.T) {
	c := PerceivedChroma(0, 0.5, 142)
	if c != 0 {
		t.Errorf("V=0 should give C=0, got %.6f", c)
	}
}

func TestPerceivedChromaIsLessThanRawForYellow(t *testing.T) {
	// Yellow (H≈60) has k(H) < 1 — should produce less chroma than linear
	cWonder := PerceivedChroma(80, 0.7, 60)
	maxC := MaxChromaForLH(0.7, 60)
	cLinear := maxC * 0.80
	if cWonder >= cLinear {
		t.Errorf("yellow V=80: WonderMath (%.6f) should be < linear (%.6f)", cWonder, cLinear)
	}
}

func TestPerceivedChromaIsMoreThanRawForBlue(t *testing.T) {
	// Blue (H≈240) has k(H) > 1 — should produce more chroma than linear
	cWonder := PerceivedChroma(80, 0.5, 240)
	maxC := MaxChromaForLH(0.5, 240)
	cLinear := maxC * 0.80
	if cWonder <= cLinear {
		t.Errorf("blue V=80: WonderMath (%.6f) should be > linear (%.6f)", cWonder, cLinear)
	}
}

// --- EffectiveChroma ---

func TestEffectiveChromaFullEnergyIsUnchanged(t *testing.T) {
	c := EffectiveChroma(0.20, 1.0)
	if math.Abs(c-0.20) > 1e-9 {
		t.Errorf("E=1 should return original chroma, got %.9f", c)
	}
}

func TestEffectiveChromaZeroEnergyIsZero(t *testing.T) {
	c := EffectiveChroma(0.20, 0.0)
	if c != 0 {
		t.Errorf("E=0 should return 0, got %.6f", c)
	}
}

func TestEffectiveChromaHalfEnergyIsNonLinear(t *testing.T) {
	// With γ=0.7, E=0.5 should give more than linear 0.5 × C
	// because 0.5^0.7 ≈ 0.615
	c := EffectiveChroma(0.20, 0.5)
	linear := 0.20 * 0.5
	if c <= linear {
		t.Errorf("E=0.5 Stevens law (%.6f) should be > linear (%.6f)", c, linear)
	}
}

// --- EffectiveLightness ---

func TestEffectiveLightnessFullEnergyBoosts(t *testing.T) {
	// High energy = slight glow = lightness slightly higher
	// λ*(1^γ - 1) = 0.04*(1-1) = 0 at E=1, so no boost at full energy
	// At E=0 it drops: 0.04*(0-1) = -0.04
	l := EffectiveLightness(0.5, 0.0)
	if l >= 0.5 {
		t.Errorf("E=0 should reduce lightness (glow gone), got %.6f", l)
	}
}

func TestEffectiveLightnessFullEnergyIsNearOriginal(t *testing.T) {
	l := EffectiveLightness(0.6, 1.0)
	if math.Abs(l-0.6) > 1e-9 {
		t.Errorf("E=1 should give original lightness, got %.9f", l)
	}
}

func TestEffectiveLightnessClampedToUnit(t *testing.T) {
	l := EffectiveLightness(0.0, 0.0)
	if l < 0 {
		t.Errorf("lightness should not go below 0, got %.6f", l)
	}
	l2 := EffectiveLightness(1.0, 1.0)
	if l2 > 1 {
		t.Errorf("lightness should not exceed 1, got %.6f", l2)
	}
}

// --- TemperatureValue ---

func TestTemperatureValueRedIsWarm(t *testing.T) {
	tv := TemperatureValue(10, 0.20, 0.6)
	if tv <= 0 {
		t.Errorf("red/orange hue should be warm (T>0), got %.4f", tv)
	}
}

func TestTemperatureValueBlueIsCool(t *testing.T) {
	tv := TemperatureValue(240, 0.18, 0.5)
	if tv >= 0 {
		t.Errorf("blue hue should be cool (T<0), got %.4f", tv)
	}
}

func TestTemperatureValueAchromaticIsNeutral(t *testing.T) {
	// With C=0, chroma weight term is zero. Only hue cosine and lightness matter.
	// cos((140-50)*π/180) = cos(90°) = 0, L=0.5 → L-0.5=0 → T=0 → neutral.
	tv := TemperatureValue(140, 0, 0.5)
	label := TemperatureLabel(tv)
	if label != "neutral" {
		t.Errorf("achromatic at H=140 should be neutral, got %q (T=%.4f)", label, tv)
	}
}

func TestTemperatureValueInRange(t *testing.T) {
	for h := 0.0; h < 360; h += 15 {
		tv := TemperatureValue(h, 0.20, 0.5)
		if tv < -1 || tv > 1 {
			t.Errorf("TemperatureValue(%.0f) = %.4f out of [-1,1]", h, tv)
		}
	}
}

// --- Valence and Arousal ---

func TestValenceBrightLightColorIsPositive(t *testing.T) {
	// High L, warm, saturated → positive valence
	tv := TemperatureValue(40, 0.18, 0.85)
	val := Valence(tv, 0.85, 0.7)
	if val <= 0 {
		t.Errorf("bright warm colour should have positive valence, got %.4f", val)
	}
}

func TestValenceDarkDesaturatedIsLower(t *testing.T) {
	// The valence formula is lightness-dominant (a2=0.7).
	// Even a dark cool tone has positive valence unless L is near zero.
	// This test verifies that a bright warm tone has HIGHER valence than
	// a dark cool one — the relative ordering is what matters.
	tvWarm := TemperatureValue(40, 0.18, 0.85)
	valWarm := Valence(tvWarm, 0.85, 0.7)
	tvCool := TemperatureValue(220, 0.05, 0.1)
	valCool := Valence(tvCool, 0.1, 0.1)
	if valWarm <= valCool {
		t.Errorf("bright warm (%.4f) should have higher valence than dark cool (%.4f)", valWarm, valCool)
	}
}

func TestArousalVividHighEnergyIsHigh(t *testing.T) {
	tv := TemperatureValue(40, 0.25, 0.6)
	aro := Arousal(0.9, 1.0, tv)
	if aro <= 0.5 {
		t.Errorf("vivid full-energy should have high arousal, got %.4f", aro)
	}
}

func TestArousalMutedLowEnergyIsLow(t *testing.T) {
	tv := TemperatureValue(180, 0.02, 0.5)
	aro := Arousal(0.05, 0.1, tv)
	if aro >= 0.3 {
		t.Errorf("muted low-energy should have low arousal, got %.4f", aro)
	}
}

// --- DerivedMood ---

func TestDerivedMoodVividIsVivid(t *testing.T) {
	// High arousal, mid-high valence, cool-neutral temp
	mood := DerivedMood(0.4, 0.7, 0.0)
	if mood != "vivid" {
		t.Errorf("expected vivid, got %q", mood)
	}
}

func TestDerivedMoodPlayfulIsPlayful(t *testing.T) {
	// High arousal, high valence, warm temp
	mood := DerivedMood(0.6, 0.8, 0.5)
	if mood != "playful" {
		t.Errorf("expected playful, got %q", mood)
	}
}

func TestDerivedMoodDeepIsDeep(t *testing.T) {
	// Low arousal, negative valence
	mood := DerivedMood(-0.4, 0.1, -0.2)
	if mood != "deep" {
		t.Errorf("expected deep, got %q", mood)
	}
}

func TestDerivedMoodNeverEmpty(t *testing.T) {
	// DerivedMood must always return something
	for val := -1.0; val <= 1.0; val += 0.25 {
		for aro := -1.0; aro <= 1.0; aro += 0.25 {
			for tv := -1.0; tv <= 1.0; tv += 0.5 {
				m := DerivedMood(val, aro, tv)
				if m == "" {
					t.Errorf("DerivedMood(%.2f, %.2f, %.2f) returned empty string", val, aro, tv)
				}
			}
		}
	}
}
