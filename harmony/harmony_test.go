package harmony_test

import (
	"testing"

	"github.com/leraniode/wondertone/harmony"
	"github.com/leraniode/wondertone/tone"
	"github.com/leraniode/x/testutil"
)

func TestComplementary(t *testing.T) {
	p, err := harmony.Complementary(red())
	testutil.NoError(t, err)
	testutil.Equal(t, 2, p.Len())
}

func TestTriadic(t *testing.T) {
	p, err := harmony.Triadic(red())
	testutil.NoError(t, err)
	testutil.Equal(t, 3, p.Len())
}

func TestAnalogous(t *testing.T) {
	p, err := harmony.Analogous(red(), 5, 30)
	testutil.NoError(t, err)
	testutil.Equal(t, 5, p.Len())
}

func TestAnalogousTooFew(t *testing.T) {
	_, err := harmony.Analogous(red(), 1, 30)
	testutil.Error(t, err)
}

func TestTetradic(t *testing.T) {
	p, err := harmony.Tetradic(red())
	testutil.NoError(t, err)
	testutil.Equal(t, 4, p.Len())
	// Each tone should be ~90° apart
	all := p.All()
	for i := 1; i < len(all); i++ {
		diff := all[i].Hue() - all[0].Hue()
		if diff < 0 {
			diff += 360
		}
		testutil.InDelta(t, float64(i)*90.0, diff, 1e-3, "tetradic step %d should be %d° from base", i, i*90)
	}
}

func TestMonochrome(t *testing.T) {
	p, err := harmony.Monochrome(red(), 6)
	testutil.NoError(t, err)
	testutil.Equal(t, 6, p.Len())
	// All should share the same hue
	for _, t2 := range p.All() {
		testutil.InDelta(t, red().Hue(), t2.Hue(), 1e-4)
	}
}

func TestRainbow(t *testing.T) {
	p, err := harmony.Rainbow(red(), 8)
	testutil.NoError(t, err)
	testutil.Equal(t, 8, p.Len())
	// Hues should be 45° apart (360/8)
	all := p.All()
	for i := 1; i < len(all); i++ {
		expected := red().Hue() + 45.0*float64(i)
		for expected >= 360 {
			expected -= 360
		}
		testutil.InDelta(t, expected, all[i].Hue(), 1e-3)
	}
}

func red() tone.Tone {
	return tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(25), tone.Named("Red"))
}
