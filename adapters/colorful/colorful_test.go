package colorful_test

import (
	"testing"

	gocolorful "github.com/lucasb-eyer/go-colorful"
	tone "github.com/leraniode/wondertone/core"
	wcolorful "github.com/leraniode/wondertone/adapters/colorful"
	"github.com/leraniode/wondertone/internal/testutil"
)

func TestToColorful(t *testing.T) {
	wt := tone.MustFromHex("#e94560")
	cf := wcolorful.ToColorful(wt)

	// Should roundtrip back to the same hex
	testutil.Equal(t, wt.Hex(), cf.Hex())
}

func TestFromColorful(t *testing.T) {
	cf, _ := gocolorful.Hex("#1a1a2e")
	wt := wcolorful.FromColorful(cf)
	testutil.Equal(t, "#1a1a2e", wt.Hex())
}

func TestToColorfulOKLCH(t *testing.T) {
	wt := tone.New(tone.Light(65), tone.Vibrancy(75), tone.Hue(142))
	cf := wcolorful.ToColorfulOKLCH(wt)

	l, c, h := cf.OkLch()
	wl, wc, wh := wt.OKLCH()

	testutil.InDelta(t, wl, l, 1e-4, "L should match")
	testutil.InDelta(t, wc, c, 1e-4, "C should match")
	testutil.InDelta(t, wh, h, 0.05, "H should match (XYZ path variance)")
}

func TestFromColorfulOKLCH(t *testing.T) {
	// go-colorful OkLch → wondertone
	cf := gocolorful.OkLch(0.68, 0.18, 142.0)
	wt := wcolorful.FromColorfulOKLCH(cf)

	l, c, h := wt.OKLCH()
	testutil.InDelta(t, 0.68, l, 1e-4)
	testutil.InDelta(t, 0.18, c, 1e-4)
	testutil.InDelta(t, 142.0, h, 0.05)
}

func TestBlendColorfulVsNative(t *testing.T) {
	// go-colorful blend and native wondertone Mix should produce nearly identical results
	a := tone.New(tone.Light(30), tone.Vibrancy(60), tone.Hue(30))
	b := tone.New(tone.Light(70), tone.Vibrancy(60), tone.Hue(200))

	native := tone.Mix(a, b, 0.5)
	via    := wcolorful.BlendColorful(a, b, 0.5)

	// Same lightness to within rounding
	testutil.InDelta(t, native.Light(), via.Light(), 2.0,
		"native Mix and go-colorful blend should agree on lightness")
}

func TestRoundtripOKLCH(t *testing.T) {
	// wondertone → go-colorful (OKLCH) → wondertone
	original := tone.New(tone.Light(55), tone.Vibrancy(80), tone.Hue(262))
	cf := wcolorful.ToColorfulOKLCH(original)
	recovered := wcolorful.FromColorfulOKLCH(cf)

	ol, oc, oh := original.OKLCH()
	rl, rc, rh := recovered.OKLCH()

	testutil.InDelta(t, ol, rl, 1e-4, "L should roundtrip")
	testutil.InDelta(t, oc, rc, 1e-4, "C should roundtrip")
	testutil.InDelta(t, oh, rh, 0.05, "H should roundtrip (XYZ path variance)")
}
