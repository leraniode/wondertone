package tone_test

import (
	"math"
	"testing"

	"github.com/leraniode/wondertone/tone"
	"github.com/leraniode/x/testutil"
)

// --- Construction ---

func TestNewDefaults(t *testing.T) {
	c := tone.New()
	testutil.InDelta(t, 50.0, c.Light(), 1e-4, "default Light should be 50")
	testutil.InDelta(t, 1.0, c.Energy(), 1e-4, "default Energy should be 1.0")
	testutil.InDelta(t, 1.0, c.AlphaValue(), 1e-4, "default Alpha should be 1.0")
}

func TestNewOptions(t *testing.T) {
	c := tone.New(
		tone.Light(75),
		tone.Vibrancy(80),
		tone.Hue(30),
		tone.Energy(0.9),
		tone.Named("Spark"),
		tone.Moody("vibrant"),
	)
	testutil.InDelta(t, 75.0, c.Light(), 1e-4)
	testutil.InDelta(t, 80.0, c.Vibrancy(), 1e-4)
	testutil.InDelta(t, 30.0, c.Hue(), 1e-4)
	testutil.InDelta(t, 0.9, c.Energy(), 1e-4)
	testutil.Equal(t, "Spark", c.Name())
	testutil.Equal(t, "vibrant", c.Mood())
}

func TestFromHex(t *testing.T) {
	cases := []string{"#e94560", "#1a1a2e", "#ffffff", "#000000", "#8b8fa8"}
	for _, hex := range cases {
		t.Run(hex, func(t *testing.T) {
			c, err := tone.FromHex(hex)
			testutil.NoError(t, err)
			testutil.Equal(t, hex, c.Hex(), "hex roundtrip failed for %s", hex)
		})
	}
}

func TestFromHexInvalid(t *testing.T) {
	_, err := tone.FromHex("notacolor")
	testutil.Error(t, err)
	_, err = tone.FromHex("#gg0000")
	testutil.Error(t, err)
	_, err = tone.FromHex("#12345")
	testutil.Error(t, err)
}

func TestMustFromHexPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustFromHex with bad input should panic")
		}
	}()
	tone.MustFromHex("not-a-color")
}

func TestFromOKLCH(t *testing.T) {
	c := tone.FromOKLCH(0.55, 0.19, 14.2)
	l, ch, h := c.OKLCH()
	testutil.InDelta(t, 0.55, l, 1e-4)
	testutil.InDelta(t, 0.19, ch, 1e-4)
	testutil.InDelta(t, 14.2, h, 1e-4)
}

func TestFromOKLCHString(t *testing.T) {
	c, err := tone.FromOKLCHString("0.55 0.19 14.2")
	testutil.NoError(t, err)
	l, ch, h := c.OKLCH()
	testutil.InDelta(t, 0.55, l, 1e-4)
	testutil.InDelta(t, 0.19, ch, 1e-4)
	testutil.InDelta(t, 14.2, h, 1e-4)
}

func TestFromOKLCHStringWithAlpha(t *testing.T) {
	c, err := tone.FromOKLCHString("0.55 0.19 14.2 / 0.8")
	testutil.NoError(t, err)
	testutil.InDelta(t, 0.8, c.AlphaValue(), 1e-4)
}

// --- Developer vocabulary ---

func TestLightVibrancyHueRanges(t *testing.T) {
	// Light clamped to [0,100]
	tooLight := tone.New(tone.Light(150))
	testutil.InDelta(t, 100.0, tooLight.Light(), 1e-4, "Light should clamp to 100")

	tooDark := tone.New(tone.Light(-10))
	testutil.InDelta(t, 0.0, tooDark.Light(), 1e-4, "Light should clamp to 0")

	// Vibrancy clamped to [0,100]
	tooVibrant := tone.New(tone.Vibrancy(200))
	testutil.InDelta(t, 100.0, tooVibrant.Vibrancy(), 1e-4)

	// Energy clamped to [0,1]
	tooEnergetic := tone.New(tone.Energy(5))
	testutil.InDelta(t, 1.0, tooEnergetic.Energy(), 1e-4)
}

func TestVibrancy0IsGrey(t *testing.T) {
	grey := tone.New(tone.Light(50), tone.Vibrancy(0), tone.Hue(30))
	_, c, _ := grey.OKLCH()
	testutil.InDelta(t, 0.0, c, 1e-6, "Vibrancy 0 should produce zero chroma")
}

func TestHueNormalization(t *testing.T) {
	c := tone.New(tone.Hue(400))
	testutil.InDelta(t, 40.0, c.Hue(), 1e-4, "Hue 400 should normalize to 40")

	c2 := tone.New(tone.Hue(-30))
	testutil.InDelta(t, 330.0, c2.Hue(), 1e-4, "Hue -30 should normalize to 330")
}

// --- Immutability ---

func TestImmutability(t *testing.T) {
	original, _ := tone.FromHex("#e94560")
	hex := original.Hex()
	_ = original.Lighten(30)
	_ = original.Rotate(90)
	_ = original.Desaturate(20)
	_ = original.WithEnergy(0.3)
	testutil.Equal(t, hex, original.Hex(), "original must not change after manipulation")
}

// --- Manipulation ---

func TestLightenDarken(t *testing.T) {
	base := tone.New(tone.Light(50), tone.Hue(30))
	lighter := base.Lighten(20)
	darker := base.Darken(20)
	testutil.Greater(t, lighter.Light(), base.Light(), "Lighten should increase Light")
	testutil.Less(t, darker.Light(), base.Light(), "Darken should decrease Light")
}

func TestHuePreservedOnLightenDarken(t *testing.T) {
	base := tone.New(tone.Light(50), tone.Hue(120))
	testutil.InDelta(t, 120.0, base.Lighten(20).Hue(), 1e-4, "hue must not change on Lighten")
	testutil.InDelta(t, 120.0, base.Darken(20).Hue(), 1e-4, "hue must not change on Darken")
}

func TestRotate(t *testing.T) {
	c := tone.New(tone.Hue(30))
	testutil.InDelta(t, 120.0, c.Rotate(90).Hue(), 1e-4)
}

func TestRotateWraps(t *testing.T) {
	c := tone.New(tone.Hue(340))
	testutil.InDelta(t, 20.0, c.Rotate(40).Hue(), 1e-4)
}

func TestComplement(t *testing.T) {
	c := tone.New(tone.Hue(30))
	comp := c.Complement()
	expected := math.Mod(30+180, 360)
	testutil.InDelta(t, expected, comp.Hue(), 1e-3)
}

func TestSaturateDesaturate(t *testing.T) {
	base := tone.New(tone.Light(50), tone.Vibrancy(50), tone.Hue(120))
	more := base.Saturate(20)
	less := base.Desaturate(20)
	testutil.Greater(t, more.Vibrancy(), base.Vibrancy(), "Saturate should increase Vibrancy")
	testutil.Less(t, less.Vibrancy(), base.Vibrancy(), "Desaturate should decrease Vibrancy")
}

// --- Energy ---

func TestEnergy(t *testing.T) {
	full := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(30), tone.Energy(1.0))
	quiet := full.WithEnergy(0.3)

	fullL, fullC, _ := full.OKLCH()
	quietL, quietC, _ := quiet.OKLCH()

	// Stored L and C are unchanged — only rendering is affected
	testutil.InDelta(t, fullL, quietL, 1e-4, "L should not change on WithEnergy")
	testutil.InDelta(t, fullC, quietC, 1e-4, "stored C should not change on WithEnergy")

	// But effective chroma differs
	testutil.Greater(t, full.EffectiveC(), quiet.EffectiveC(), "EffectiveC should be lower with lower Energy")
}

// --- Accessibility ---

func TestContrastBlackWhite(t *testing.T) {
	black, _ := tone.FromHex("#000000")
	white, _ := tone.FromHex("#ffffff")
	ratio := black.ContrastWith(white)
	testutil.InDelta(t, 21.0, ratio, 0.1, "black on white should be ~21:1")
}

func TestPassesAA(t *testing.T) {
	black, _ := tone.FromHex("#000000")
	white, _ := tone.FromHex("#ffffff")
	testutil.True(t, black.PassesAA(white), "black on white should pass AA")
	testutil.True(t, black.PassesAAA(white), "black on white should pass AAA")
}

func TestEnsureContrast(t *testing.T) {
	fg := tone.New(tone.Light(55), tone.Vibrancy(0), tone.Hue(0)) // mid grey
	bg, _ := tone.FromHex("#ffffff")
	if fg.PassesAA(bg) {
		t.Skip("test tone already passes AA — update the test")
	}
	fixed := fg.EnsureContrast(bg, "AA")
	testutil.True(t, fixed.PassesAA(bg), "EnsureContrast result should pass AA")
}

// --- Intelligence ---

func TestIsLightIsDark(t *testing.T) {
	light := tone.New(tone.Light(75))
	dark := tone.New(tone.Light(25))
	testutil.True(t, light.IsLight())
	testutil.True(t, dark.IsDark())
	testutil.False(t, light.IsDark())
	testutil.False(t, dark.IsLight())
}

func TestTemperature(t *testing.T) {
	warm := tone.New(tone.Light(50), tone.Vibrancy(60), tone.Hue(30))
	cool := tone.New(tone.Light(50), tone.Vibrancy(60), tone.Hue(210))
	// v0.2: Temperature uses WonderMath continuous formula.
	// With C=0, only hue cosine matters. cos((140-50)*π/180)=cos(90°)=0 → neutral.
	trueGrey := tone.New(tone.Light(50), tone.Vibrancy(0), tone.Hue(140))
	testutil.Equal(t, "warm", warm.Temperature())
	testutil.Equal(t, "cool", cool.Temperature())
	testutil.Equal(t, "neutral", trueGrey.Temperature())
}

// --- Gamut safety ---

func TestAllOutputsInGamut(t *testing.T) {
	cases := []struct {
		name string
		tone tone.Tone
	}{
		{"vivid red", tone.New(tone.Light(50), tone.Vibrancy(100), tone.Hue(25))},
		{"vivid green", tone.New(tone.Light(50), tone.Vibrancy(100), tone.Hue(142))},
		{"vivid blue", tone.New(tone.Light(50), tone.Vibrancy(100), tone.Hue(250))},
		{"light tone", tone.New(tone.Light(95), tone.Vibrancy(80), tone.Hue(30))},
		{"dark tone", tone.New(tone.Light(10), tone.Vibrancy(80), tone.Hue(200))},
		{"from hex", tone.MustFromHex("#e94560")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, g, b := tc.tone.RGBFloat()
			testutil.GreaterOrEqual(t, r, 0.0, "R out of gamut (below 0)")
			testutil.LessOrEqual(t, r, 1.0, "R out of gamut (above 1)")
			testutil.GreaterOrEqual(t, g, 0.0, "G out of gamut (below 0)")
			testutil.LessOrEqual(t, g, 1.0, "G out of gamut (above 1)")
			testutil.GreaterOrEqual(t, b, 0.0, "B out of gamut (below 0)")
			testutil.LessOrEqual(t, b, 1.0, "B out of gamut (above 1)")
		})
	}
}

// --- ToneScale ---

func TestScaleLength(t *testing.T) {
	c := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(30))
	testutil.Equal(t, 12, len(c.Scale().All()))
}

func TestScaleOrdering(t *testing.T) {
	c := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(30))
	all := c.Scale().All()
	for i := 1; i < len(all); i++ {
		testutil.Greater(t, all[i-1].Light(), all[i].Light(),
			"tone %d should be lighter than tone %d", i, i+1)
	}
}

func TestScaleHueLocked(t *testing.T) {
	c := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(142))
	for i, step := range c.Scale().All() {
		testutil.InDelta(t, 142.0, step.Hue(), 1e-4, "tone %d hue must match base hue", i+1)
	}
}

func TestScaleInGamut(t *testing.T) {
	c := tone.New(tone.Light(50), tone.Vibrancy(100), tone.Hue(25))
	for i, step := range c.Scale().All() {
		r, g, b := step.RGBFloat()
		testutil.GreaterOrEqual(t, r, 0.0, "tone %d R below gamut", i+1)
		testutil.LessOrEqual(t, r, 1.0, "tone %d R above gamut", i+1)
		testutil.GreaterOrEqual(t, g, 0.0, "tone %d G below gamut", i+1)
		testutil.LessOrEqual(t, g, 1.0, "tone %d G above gamut", i+1)
		testutil.GreaterOrEqual(t, b, 0.0, "tone %d B below gamut", i+1)
		testutil.LessOrEqual(t, b, 1.0, "tone %d B above gamut", i+1)
	}
}

func TestScaleSemanticAccessors(t *testing.T) {
	c := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(30))
	scale := c.Scale()
	testutil.Greater(t, scale.Background().Light(), scale.Solid().Light(),
		"Background should be lighter than Solid")
	testutil.Less(t, scale.Text().Light(), scale.Solid().Light(),
		"Text should be darker than Solid")
}

func TestScaleStepClamping(t *testing.T) {
	scale := tone.New(tone.Light(50), tone.Hue(30)).Scale()
	testutil.Equal(t, scale.Step(1).Hex(), scale.Step(0).Hex(), "Step(0) should clamp to Step(1)")
	testutil.Equal(t, scale.Step(12).Hex(), scale.Step(99).Hex(), "Step(99) should clamp to Step(12)")
}

// --- Tone method integration ---
func TestToneTemperatureScalar(t *testing.T) {
	warm := tone.New(tone.Light(65), tone.Vibrancy(70), tone.Hue(30))
	cool := tone.New(tone.Light(65), tone.Vibrancy(70), tone.Hue(210))
	if warm.TemperatureScalar() <= cool.TemperatureScalar() {
		t.Errorf("warm tone (%.4f) should have higher T than cool (%.4f)",
			warm.TemperatureScalar(), cool.TemperatureScalar())
	}
}

func TestToneDerivedMoodIsNonEmpty(t *testing.T) {
	tones := []tone.Tone{
		tone.New(tone.Light(80), tone.Vibrancy(90), tone.Hue(40), tone.Energy(1.0)),
		tone.New(tone.Light(20), tone.Vibrancy(30), tone.Hue(220), tone.Energy(0.3)),
		tone.New(tone.Light(50), tone.Vibrancy(5), tone.Hue(0)),
	}
	for _, tone := range tones {
		if tone.DerivedMoodValue() == "" {
			t.Errorf("DerivedMoodValue should never be empty, got empty for %s", tone.Hex())
		}
	}
}

func TestToneValenceAndArousalInRange(t *testing.T) {
	tn := tone.New(tone.Light(70), tone.Vibrancy(80), tone.Hue(30), tone.Energy(0.9))
	v := tn.ValenceValue()
	a := tn.ArousalValue()
	if v < -1 || v > 1 {
		t.Errorf("ValenceValue %.4f out of [-1,1]", v)
	}
	if a < -1 || a > 1 {
		t.Errorf("ArousalValue %.4f out of [-1,1]", a)
	}
}

func TestToneEffectiveCUsesStevensLaw(t *testing.T) {
	base := tone.New(tone.Light(60), tone.Vibrancy(80), tone.Hue(142))
	half := base.WithEnergy(0.5)
	// With γ=0.7: E=0.5 → 0.5^0.7 ≈ 0.615 of base C
	// Must be more than linear 0.5× but less than full
	if half.EffectiveC() <= base.EffectiveC()*0.5 {
		t.Errorf("Stevens law: E=0.5 effective C should be > linear half")
	}
	if half.EffectiveC() >= base.EffectiveC() {
		t.Errorf("E=0.5 should reduce chroma vs E=1.0")
	}
}
