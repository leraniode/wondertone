package colour_test

import (
	"testing"

	"github.com/leraniode/wondertone/colour"
	"github.com/leraniode/wondertone/internal/testutil"
)

func TestAllTonesLoad(t *testing.T) {
	all := colour.All()
	testutil.Equal(t, 12, len(all))
}

func TestAllTonesHaveNames(t *testing.T) {
	for _, tone := range colour.All() {
		testutil.True(t, len(tone.Name()) > 0, "every tone must have a name")
	}
}

func TestAllTonesHaveMoods(t *testing.T) {
	for _, tone := range colour.All() {
		testutil.True(t, len(tone.Mood()) > 0, "%s must have a mood", tone.Name())
	}
}

func TestAllTonesInGamut(t *testing.T) {
	for _, tone := range colour.All() {
		t.Run(tone.Name(), func(t *testing.T) {
			r, g, b := tone.RGBFloat()
			testutil.GreaterOrEqual(t, r, 0.0, "R below gamut")
			testutil.LessOrEqual(t, r, 1.0, "R above gamut")
			testutil.GreaterOrEqual(t, g, 0.0, "G below gamut")
			testutil.LessOrEqual(t, g, 1.0, "G above gamut")
			testutil.GreaterOrEqual(t, b, 0.0, "B below gamut")
			testutil.LessOrEqual(t, b, 1.0, "B above gamut")
		})
	}
}

func TestAllTonesHaveValidEnergy(t *testing.T) {
	for _, tone := range colour.All() {
		e := tone.Energy()
		testutil.True(t, e > 0 && e <= 1.0, "%s energy should be in (0,1], got %f", tone.Name(), e)
	}
}

func TestNamedTonesCorrect(t *testing.T) {
	testutil.Equal(t, "Unix", colour.Unix.Name())
	testutil.Equal(t, "Starlight", colour.Starlight.Name())
	testutil.Equal(t, "Ember", colour.Ember.Name())
	testutil.Equal(t, "Glacier", colour.Glacier.Name())
	testutil.Equal(t, "Crimson", colour.Crimson.Name())
	testutil.Equal(t, "Void", colour.Void.Name())
	testutil.Equal(t, "Dawn", colour.Dawn.Name())
	testutil.Equal(t, "Bloom", colour.Bloom.Name())
	testutil.Equal(t, "Slate", colour.Slate.Name())
	testutil.Equal(t, "Signal", colour.Signal.Name())
	testutil.Equal(t, "Ink", colour.Ink.Name())
	testutil.Equal(t, "Paper", colour.Paper.Name())
}

func TestUnixIsGreen(t *testing.T) {
	// Unix is a green tone — hue should be in green territory
	testutil.True(t, colour.Unix.Hue() > 100 && colour.Unix.Hue() < 180,
		"Unix hue should be in green range, got %f", colour.Unix.Hue())
}

func TestVoidIsDark(t *testing.T) {
	testutil.True(t, colour.Void.IsDark(), "Void should be dark")
}

func TestPaperIsLight(t *testing.T) {
	testutil.True(t, colour.Paper.IsLight(), "Paper should be light")
}

func TestCrimsonIsWarm(t *testing.T) {
	testutil.Equal(t, "warm", colour.Crimson.Temperature())
}

func TestGlacierIsCool(t *testing.T) {
	testutil.Equal(t, "cool", colour.Glacier.Temperature())
}
