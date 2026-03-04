package colour

import tone "github.com/leraniode/wondertone/core"

// Ember is Leraniode's warm amber — the last light of something burning well.
// Used for warnings that feel like wisdom, not alarm.
var Ember = tone.New(
	tone.Light(72),
	tone.Vibrancy(85),
	tone.Hue(38),
	tone.Energy(0.90),
	tone.Named("Ember"),
	tone.Moody("warm"),
)
