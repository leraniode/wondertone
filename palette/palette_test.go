package palette_test

import (
	"testing"

	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/internal/testutil"
	"github.com/leraniode/wondertone/palette"
)

// helpers
func red() tone.Tone {
	return tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(25), tone.Named("Red"))
}
func green() tone.Tone {
	return tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(142), tone.Named("Green"))
}
func blue() tone.Tone {
	return tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(250), tone.Named("Blue"))
}

// --- Construction ---

func TestBuildBasic(t *testing.T) {
	p, err := palette.New("test").Add(red()).Add(green()).Add(blue()).Build()
	testutil.NoError(t, err)
	testutil.Equal(t, "test", p.Name())
	testutil.Equal(t, 3, p.Len())
}

func TestBuildEmptyFails(t *testing.T) {
	_, err := palette.New("empty").Build()
	testutil.Error(t, err)
}

func TestBuildDuplicateNameFails(t *testing.T) {
	_, err := palette.New("dup").Add(red()).Add(red()).Build()
	testutil.Error(t, err)
}

func TestBuildUnnamedToneFails(t *testing.T) {
	unnamed := tone.New(tone.Light(50), tone.Hue(30)) // no Name
	_, err := palette.New("test").Add(unnamed).Build()
	testutil.Error(t, err)
}

func TestMustBuildPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustBuild should panic on empty palette")
		}
	}()
	palette.New("panic").MustBuild()
}

// --- Accessors ---

func TestGetExists(t *testing.T) {
	p, _ := palette.New("test").Add(red()).Add(green()).Build()
	t1, ok := p.Get("Red")
	testutil.True(t, ok)
	testutil.Equal(t, "Red", t1.Name())
}

func TestGetMissing(t *testing.T) {
	p, _ := palette.New("test").Add(red()).Build()
	_, ok := p.Get("NotHere")
	testutil.False(t, ok)
}

func TestMustGetPanics(t *testing.T) {
	p, _ := palette.New("test").Add(red()).Build()
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet should panic on missing tone")
		}
	}()
	p.MustGet("NotHere")
}

func TestAt(t *testing.T) {
	p, _ := palette.New("test").Add(red()).Add(green()).Add(blue()).Build()
	testutil.Equal(t, "Red", p.At(0).Name())
	testutil.Equal(t, "Green", p.At(1).Name())
	testutil.Equal(t, "Blue", p.At(2).Name())
}

func TestAll(t *testing.T) {
	p, _ := palette.New("test").Add(red()).Add(green()).Add(blue()).Build()
	all := p.All()
	testutil.Equal(t, 3, len(all))
}

func TestHas(t *testing.T) {
	p, _ := palette.New("test").Add(red()).Build()
	testutil.True(t, p.Has("Red"))
	testutil.False(t, p.Has("Blue"))
}

// --- Operations ---

func TestFork(t *testing.T) {
	p, _ := palette.New("original").Add(red()).Add(green()).Build()
	fork, err := p.Fork("fork").Build()
	testutil.NoError(t, err)
	testutil.Equal(t, "fork", fork.Name())
	testutil.Equal(t, 2, fork.Len())
	testutil.True(t, fork.Has("Red"))
	testutil.True(t, fork.Has("Green"))
}

func TestForkIsIndependent(t *testing.T) {
	original, _ := palette.New("original").Add(red()).Add(green()).Build()
	fork, _ := original.Fork("fork").Add(blue()).Build()
	// original should not have Blue
	testutil.Equal(t, 2, original.Len())
	testutil.Equal(t, 3, fork.Len())
}

func TestExtend(t *testing.T) {
	p, _ := palette.New("base").Add(red()).Add(green()).Build()
	extended, err := p.Extend("extended", blue())
	testutil.NoError(t, err)
	testutil.Equal(t, 3, extended.Len())
	testutil.True(t, extended.Has("Blue"))
}

func TestExtendCannotOverride(t *testing.T) {
	p, _ := palette.New("base").Add(red()).Build()
	_, err := p.Extend("override", red()) // Red already exists
	testutil.Error(t, err)
}

func TestReplace(t *testing.T) {
	p, _ := palette.New("base").Add(red()).Add(green()).Build()
	newRed := tone.New(tone.Light(60), tone.Hue(25), tone.Named("Red"))
	updated, err := p.Replace("Red", newRed)
	testutil.NoError(t, err)
	testutil.InDelta(t, 60.0, updated.MustGet("Red").Light(), 1e-4)
	// Green must still be there
	testutil.True(t, updated.Has("Green"))
}

func TestReplaceMissing(t *testing.T) {
	p, _ := palette.New("base").Add(red()).Build()
	_, err := p.Replace("NotHere", green())
	testutil.Error(t, err)
}

func TestWithEnergy(t *testing.T) {
	p, _ := palette.New("base").Add(red()).Add(green()).Build()
	quiet := p.WithEnergy(0.3)
	for _, t2 := range quiet.All() {
		testutil.InDelta(t, 0.3, t2.Energy(), 1e-4, "all tones should have energy 0.3")
	}
}

// --- Contrast ---

func TestContrastPair(t *testing.T) {
	black := tone.New(tone.Light(2), tone.Vibrancy(0), tone.Hue(0), tone.Named("Black"))
	white := tone.New(tone.Light(98), tone.Vibrancy(0), tone.Hue(0), tone.Named("White"))
	p, _ := palette.New("bw").Add(black).Add(white).Build()

	ratio, err := palette.ContrastPair(p, "Black", "White")
	testutil.NoError(t, err)
	testutil.Greater(t, ratio, 15.0, "black on white contrast should be high")
}

func TestContrastPairMissing(t *testing.T) {
	p, _ := palette.New("test").Add(red()).Build()
	_, err := palette.ContrastPair(p, "Red", "NotHere")
	testutil.Error(t, err)
}

func TestFindReadablePairs(t *testing.T) {
	black := tone.New(tone.Light(2), tone.Vibrancy(0), tone.Hue(0), tone.Named("Black"))
	white := tone.New(tone.Light(98), tone.Vibrancy(0), tone.Hue(0), tone.Named("White"))
	p, _ := palette.New("bw").Add(black).Add(white).Build()

	pairs := palette.FindReadablePairs(p, "AA")
	testutil.True(t, len(pairs) > 0, "should find at least one readable pair")
	for _, pair := range pairs {
		testutil.True(t, pair.PassesAA, "all found pairs should pass AA")
	}
}

// --- Validation ---

func TestValidatePasses(t *testing.T) {
	dark := tone.New(tone.Light(15), tone.Vibrancy(60), tone.Hue(25), tone.Named("Dark"))
	mid := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(142), tone.Named("Mid"))
	light := tone.New(tone.Light(85), tone.Vibrancy(40), tone.Hue(250), tone.Named("Light"))
	p, _ := palette.New("valid").Add(dark).Add(mid).Add(light).Build()
	report := p.Validate()
	testutil.True(t, report.Passed, "well-spaced palette should pass: %s", report)
}

// --- Harmony ---

func TestComplementary(t *testing.T) {
	p, err := palette.Complementary(red())
	testutil.NoError(t, err)
	testutil.Equal(t, 2, p.Len())
}

func TestTriadic(t *testing.T) {
	p, err := palette.Triadic(red())
	testutil.NoError(t, err)
	testutil.Equal(t, 3, p.Len())
}

func TestAnalogous(t *testing.T) {
	p, err := palette.Analogous(red(), 5, 30)
	testutil.NoError(t, err)
	testutil.Equal(t, 5, p.Len())
}

func TestAnalogousTooFew(t *testing.T) {
	_, err := palette.Analogous(red(), 1, 30)
	testutil.Error(t, err)
}

func TestTetradic(t *testing.T) {
	p, err := palette.Tetradic(red())
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
	p, err := palette.Monochrome(red(), 6)
	testutil.NoError(t, err)
	testutil.Equal(t, 6, p.Len())
	// All should share the same hue
	for _, t2 := range p.All() {
		testutil.InDelta(t, red().Hue(), t2.Hue(), 1e-4)
	}
}

func TestRainbow(t *testing.T) {
	p, err := palette.Rainbow(red(), 8)
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
