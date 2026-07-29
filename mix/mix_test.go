package mix_test

import (
	"testing"

	"github.com/leraniode/wondertone/mix"
	"github.com/leraniode/wondertone/tone"
	"github.com/leraniode/x/testutil"
)

func TestMixMidpoint(t *testing.T) {
	a := tone.New(tone.Light(30), tone.Vibrancy(50), tone.Hue(30))
	b := tone.New(tone.Light(70), tone.Vibrancy(50), tone.Hue(90))
	mid := mix.Mix(a, b, 0.5)
	testutil.InDelta(t, 50.0, mid.Light(), 3.0, "midpoint Light should be ~50")
}

func TestMixBoundaries(t *testing.T) {
	a := tone.New(tone.Light(30), tone.Hue(30))
	b := tone.New(tone.Light(70), tone.Hue(90))
	testutil.Equal(t, a.Hex(), mix.Mix(a, b, 0).Hex(), "Mix(a,b,0) should equal a")
	testutil.Equal(t, b.Hex(), mix.Mix(a, b, 1).Hex(), "Mix(a,b,1) should equal b")
}

func TestGradient(t *testing.T) {
	a := tone.New(tone.Light(20), tone.Hue(30))
	b := tone.New(tone.Light(80), tone.Hue(200))
	grad, err := mix.Gradient(a, b, 10)
	testutil.NoError(t, err)
	testutil.Equal(t, 10, len(grad))
	testutil.Equal(t, a.Hex(), grad[0].Hex())
	testutil.Equal(t, b.Hex(), grad[9].Hex())
}

func TestGradientTooFewSteps(t *testing.T) {
	a := tone.New(tone.Light(20), tone.Hue(30))
	b := tone.New(tone.Light(80), tone.Hue(200))
	_, err := mix.Gradient(a, b, 1)
	testutil.Error(t, err)
}

func TestHarmonize(t *testing.T) {
	base := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(30))

	comp, err := mix.Harmonize(base, "complement")
	testutil.NoError(t, err)
	testutil.Equal(t, 2, len(comp))
	testutil.InDelta(t, 210.0, comp[1].Hue(), 1e-3)

	triadic, err := mix.Harmonize(base, "triadic")
	testutil.NoError(t, err)
	testutil.Equal(t, 3, len(triadic))

	_, err = mix.Harmonize(base, "unknown")
	testutil.Error(t, err)
}

func TestEqual(t *testing.T) {
	a := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(30))
	b := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(30))
	c := tone.New(tone.Light(60), tone.Vibrancy(80), tone.Hue(30))
	testutil.True(t, a.Equal(b), "same tones should be equal")
	testutil.False(t, a.Equal(c), "different tones should not be equal")
}
