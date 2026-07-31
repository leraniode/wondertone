package tone

import "github.com/leraniode/wondertone/space"

// IsLight reports whether the Tone is perceptually light (Light > 50).
func (t Tone) IsLight() bool { return t.light > 50 }

// IsDark reports whether the Tone is perceptually dark (Light <= 50).
func (t Tone) IsDark() bool { return t.light <= 50 }

// Temperature returns "warm", "cool", or "neutral".
// Upgraded in v0.2: uses WonderMath continuous formula instead of hue-range
// lookup — chroma and lightness now modulate the reading.
func (t Tone) Temperature() string {
	tv := space.TemperatureValue(t.hue, t.c, t.l)
	return space.TemperatureLabel(tv)
}

// TemperatureScalar returns the continuous warm↔cool value T ∈ [-1, +1].
// +1 = maximally warm, -1 = maximally cool, 0 = neutral.
// More precise than Temperature() which returns a label.
func (t Tone) TemperatureScalar() float64 {
	return space.TemperatureValue(t.hue, t.c, t.l)
}

// DerivedMoodValue returns the mathematically derived mood string.
// Computed from valence and arousal — independent of the stored Mood() tag.
// Use Mood() for the display label (manual override takes precedence).
func (t Tone) DerivedMoodValue() string {
	s := space.NormalizedSaturation(t.c, t.l, t.hue)
	tv := space.TemperatureValue(t.hue, t.c, t.l)
	val := space.Valence(tv, t.l, s)
	aro := space.Arousal(s, t.energy, tv)
	return space.DerivedMood(val, aro, tv)
}

// Valence returns the emotional valence of this tone ∈ [-1, +1].
// +1 = positive (bright, warm, vivid). -1 = negative (dark, cool, muted).
func (t Tone) ValenceValue() float64 {
	s := space.NormalizedSaturation(t.c, t.l, t.hue)
	tv := space.TemperatureValue(t.hue, t.c, t.l)
	return space.Valence(tv, t.l, s)
}

// ArousalValue returns the emotional arousal of this tone ∈ [-1, +1].
// +1 = activated (vivid, energetic). -1 = calm (muted, quiet).
func (t Tone) ArousalValue() float64 {
	s := space.NormalizedSaturation(t.c, t.l, t.hue)
	tv := space.TemperatureValue(t.hue, t.c, t.l)
	return space.Arousal(s, t.energy, tv)
}
