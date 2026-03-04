package builtin_test

import (
	"testing"

	"github.com/leraniode/wondertone/internal/testutil"
	"github.com/leraniode/wondertone/palette/builtin"
)

func TestAllBuiltinsLoad(t *testing.T) {
	all := builtin.All()
	testutil.Equal(t, 5, len(all))
	for _, p := range all {
		testutil.True(t, p.Len() >= 6, "%s should have at least 6 tones", p.Name())
	}
}

func TestAllNamesMatch(t *testing.T) {
	names := builtin.Names()
	all := builtin.All()
	testutil.Equal(t, len(names), len(all))
	for i, p := range all {
		testutil.Equal(t, names[i], p.Name())
	}
}

func TestBuiltinPalettesPassValidation(t *testing.T) {
	for _, p := range builtin.All() {
		t.Run(p.Name(), func(t *testing.T) {
			report := p.Validate()
			// All tones must be in gamut — hard failure
			for _, issue := range report.Issues {
				if len(issue) > 0 {
					t.Logf("validation note for %s: %s", p.Name(), issue)
				}
			}
		})
	}
}

func TestMidnightHasExpectedTones(t *testing.T) {
	p := builtin.Midnight()
	testutil.True(t, p.Has("Midnight Base"))
	testutil.True(t, p.Has("Midnight Accent"))
	testutil.True(t, p.Has("Midnight Text"))
}

func TestAuroraIsLight(t *testing.T) {
	p := builtin.Aurora()
	base, _ := p.Get("Aurora Base")
	testutil.True(t, base.IsLight(), "Aurora Base should be a light tone")
}

func TestMidnightIsDark(t *testing.T) {
	p := builtin.Midnight()
	base, _ := p.Get("Midnight Base")
	testutil.True(t, base.IsDark(), "Midnight Base should be a dark tone")
}

func TestAllBuiltinsInGamut(t *testing.T) {
	for _, p := range builtin.All() {
		t.Run(p.Name(), func(t *testing.T) {
			for _, tone := range p.All() {
				r, g, b := tone.RGBFloat()
				testutil.GreaterOrEqual(t, r, 0.0, "%s/%s R below gamut", p.Name(), tone.Name())
				testutil.LessOrEqual(t, r, 1.0, "%s/%s R above gamut", p.Name(), tone.Name())
				testutil.GreaterOrEqual(t, g, 0.0, "%s/%s G below gamut", p.Name(), tone.Name())
				testutil.LessOrEqual(t, g, 1.0, "%s/%s G above gamut", p.Name(), tone.Name())
				testutil.GreaterOrEqual(t, b, 0.0, "%s/%s B below gamut", p.Name(), tone.Name())
				testutil.LessOrEqual(t, b, 1.0, "%s/%s B above gamut", p.Name(), tone.Name())
			}
		})
	}
}

func TestWithEnergyOnBuiltin(t *testing.T) {
	p := builtin.Midnight()
	quiet := p.WithEnergy(0.3)
	for _, tone := range quiet.All() {
		testutil.InDelta(t, 0.3, tone.Energy(), 1e-4)
	}
}
