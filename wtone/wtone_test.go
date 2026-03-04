package wtone_test

import (
	"os"
	"path/filepath"
	"testing"

	tone "github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/internal/testutil"
	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/wtone"
)

func buildTestPalette() *palette.Palette {
	unix := tone.New(
		tone.Light(68), tone.Vibrancy(72), tone.Hue(142),
		tone.Energy(0.95), tone.Named("Unix"), tone.Moody("focused"),
	)
	ember := tone.New(
		tone.Light(72), tone.Vibrancy(85), tone.Hue(38),
		tone.Energy(0.9), tone.Named("Ember"), tone.Moody("warm"),
	)
	glacier := tone.New(
		tone.Light(74), tone.Vibrancy(58), tone.Hue(196),
		tone.Energy(0.82), tone.Named("Glacier"),
	)
	return palette.New("Leraniode").
		Description("Leraniode brand tones").
		Mood("focused").
		Author("leraniode").
		Version("1.0.0").
		Add(unix).Add(ember).Add(glacier).
		MustBuild()
}

func TestMarshalRoundtrip(t *testing.T) {
	original := buildTestPalette()

	data, err := wtone.MarshalWTone(original)
	testutil.NoError(t, err)
	testutil.True(t, len(data) > 0, "marshalled data should not be empty")

	loaded, err := wtone.ParseWTone(data)
	testutil.NoError(t, err)

	testutil.Equal(t, original.Name(), loaded.Name())
	testutil.Equal(t, original.Description(), loaded.Description())
	testutil.Equal(t, original.Author(), loaded.Author())
	testutil.Equal(t, original.Len(), loaded.Len())
}

func TestRoundtripPreservesValues(t *testing.T) {
	original := buildTestPalette()
	data, _ := wtone.MarshalWTone(original)
	loaded, err := wtone.ParseWTone(data)
	testutil.NoError(t, err)

	for _, orig := range original.All() {
		loaded_t, ok := loaded.Get(orig.Name())
		testutil.True(t, ok, "tone %q should exist after roundtrip", orig.Name())

		ol, oc, oh := orig.OKLCH()
		ll, lc, lh := loaded_t.OKLCH()

		testutil.InDelta(t, ol, ll, 1e-4, "L should roundtrip for %s", orig.Name())
		testutil.InDelta(t, oc, lc, 1e-4, "C should roundtrip for %s", orig.Name())
		testutil.InDelta(t, oh, lh, 1e-3, "H should roundtrip for %s", orig.Name())
		testutil.InDelta(t, orig.Energy(), loaded_t.Energy(), 1e-3, "Energy should roundtrip for %s", orig.Name())
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.wtone")

	original := buildTestPalette()
	err := wtone.SaveWTone(path, original)
	testutil.NoError(t, err)

	// File should exist
	_, statErr := os.Stat(path)
	testutil.NoError(t, statErr)

	loaded, err := wtone.LoadWTone(path)
	testutil.NoError(t, err)
	testutil.Equal(t, original.Name(), loaded.Name())
	testutil.Equal(t, original.Len(), loaded.Len())
}

func TestLoadMissingFile(t *testing.T) {
	_, err := wtone.LoadWTone("/nonexistent/path/palette.wtone")
	testutil.Error(t, err)
}

func TestParseOKLCHShorthand(t *testing.T) {
	// Shorthand via ParseWTone
	content := []byte(`
name = "shorthand"
[[colors]]
name  = "Test"
oklch = "0.68 0.18 142"
energy = 0.9
`)
	p, err := wtone.ParseWTone(content)
	testutil.NoError(t, err)

	t2, ok := p.Get("Test")
	testutil.True(t, ok)
	l, c, h := t2.OKLCH()
	testutil.InDelta(t, 0.68, l, 1e-4)
	testutil.InDelta(t, 0.18, c, 1e-4)
	testutil.InDelta(t, 142.0, h, 1e-3)
}

func TestParseMoodInheritance(t *testing.T) {
	content := []byte(`
name = "inherit"
mood = "serene"
[[colors]]
name = "One"
l = 0.5
c = 0.1
h = 200.0
[[colors]]
name = "Two"
l = 0.6
c = 0.12
h = 210.0
mood = "focused"
`)
	p, err := wtone.ParseWTone(content)
	testutil.NoError(t, err)

	one, _ := p.Get("One")
	two, _ := p.Get("Two")

	testutil.Equal(t, "serene", one.Mood(), "One should inherit palette mood")
	testutil.Equal(t, "focused", two.Mood(), "Two should keep its own mood")
}

func TestParseValidationErrors(t *testing.T) {
	// Missing name
	_, err := wtone.ParseWTone([]byte(`
[[colors]]
l = 0.5
c = 0.1
h = 30.0
`))
	testutil.Error(t, err)

	// No colors
	_, err = wtone.ParseWTone([]byte(`name = "empty"`))
	testutil.Error(t, err)

	// Color missing name
	_, err = wtone.ParseWTone([]byte(`
name = "test"
[[colors]]
l = 0.5
c = 0.1
h = 30.0
`))
	testutil.Error(t, err)
}
