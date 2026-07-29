package contrast_test

import (
	"testing"

	"github.com/leraniode/wondertone/contrast"
	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/tone"
	"github.com/leraniode/x/testutil"
)

func TestContrastPair(t *testing.T) {
	black := tone.New(tone.Light(2), tone.Vibrancy(0), tone.Hue(0), tone.Named("Black"))
	white := tone.New(tone.Light(98), tone.Vibrancy(0), tone.Hue(0), tone.Named("White"))
	p, _ := palette.New("bw").Add(black).Add(white).Build()

	ratio, err := contrast.ContrastPair(p, "Black", "White")
	testutil.NoError(t, err)
	testutil.Greater(t, ratio, 15.0, "black on white contrast should be high")
}

func TestContrastPairMissing(t *testing.T) {
	p, _ := palette.New("test").Add(red()).Build()
	_, err := contrast.ContrastPair(p, "Red", "NotHere")
	testutil.Error(t, err)
}

func TestFindReadablePairs(t *testing.T) {
	black := tone.New(tone.Light(2), tone.Vibrancy(0), tone.Hue(0), tone.Named("Black"))
	white := tone.New(tone.Light(98), tone.Vibrancy(0), tone.Hue(0), tone.Named("White"))
	p, _ := palette.New("bw").Add(black).Add(white).Build()

	pairs := contrast.FindReadablePairs(p, "AA")
	testutil.True(t, len(pairs) > 0, "should find at least one readable pair")
	for _, pair := range pairs {
		testutil.True(t, pair.PassesAA, "all found pairs should pass AA")
	}
}

func red() tone.Tone {
	return tone.New(tone.Light(50), tone.Vibrancy(80), tone.Hue(25), tone.Named("Red"))
}
