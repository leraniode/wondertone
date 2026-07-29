package palette_test

import (
	"testing"

	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/tone"
	"github.com/leraniode/x/testutil"
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

// --- Validation ---

func TestValidatePasses(t *testing.T) {
	dark := tone.New(tone.Light(15), tone.Vibrancy(60), tone.Hue(25), tone.Named("Dark"))
	mid := tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(142), tone.Named("Mid"))
	light := tone.New(tone.Light(85), tone.Vibrancy(40), tone.Hue(250), tone.Named("Light"))
	p, _ := palette.New("valid").Add(dark).Add(mid).Add(light).Build()
	report := p.Validate()
	testutil.True(t, report.Passed, "well-spaced palette should pass: %s", report)
}
