package tone

import "github.com/leraniode/wondertone/space"

// WithLight returns a new Tone with the given lightness [0–100].
func (t Tone) WithLight(v float64) Tone {
	return New(
		Light(v), Vibrancy(t.vibrancy), Hue(t.hue), Energy(t.energy),
		Named(t.name), Moody(t.mood), Alpha(t.a),
	)
}

// WithVibrancy returns a new Tone with the given vibrancy [0–100].
func (t Tone) WithVibrancy(v float64) Tone {
	return New(
		Light(t.light), Vibrancy(v), Hue(t.hue), Energy(t.energy),
		Named(t.name), Moody(t.mood), Alpha(t.a),
	)
}

// WithHue returns a new Tone with the given hue [0–360).
func (t Tone) WithHue(h float64) Tone {
	return New(
		Light(t.light), Vibrancy(t.vibrancy), Hue(h), Energy(t.energy),
		Named(t.name), Moody(t.mood), Alpha(t.a),
	)
}

// WithEnergy returns a new Tone with the given energy [0–1].
func (t Tone) WithEnergy(e float64) Tone {
	nt := t
	nt.energy = space.Clamp(e, 0, 1)
	return nt
}

// WithName returns a new Tone with the given name.
func (t Tone) WithName(name string) Tone {
	nt := t
	nt.name = name
	return nt
}

// WithMood returns a new Tone with the given mood.
func (t Tone) WithMood(mood string) Tone {
	nt := t
	nt.mood = mood
	return nt
}

// WithAlpha returns a new Tone with the given alpha [0–1].
func (t Tone) WithAlpha(a float64) Tone {
	nt := t
	nt.a = space.Clamp(a, 0, 1)
	return nt
}

// Lighten returns a new Tone with lightness increased by amount [0–100].
func (t Tone) Lighten(amount float64) Tone {
	return t.WithLight(t.light + amount)
}

// Darken returns a new Tone with lightness decreased by amount [0–100].
func (t Tone) Darken(amount float64) Tone {
	return t.WithLight(t.light - amount)
}

// Saturate returns a new Tone with vibrancy increased by amount [0–100].
func (t Tone) Saturate(amount float64) Tone {
	return t.WithVibrancy(t.vibrancy + amount)
}

// Desaturate returns a new Tone with vibrancy decreased by amount [0–100].
func (t Tone) Desaturate(amount float64) Tone {
	return t.WithVibrancy(t.vibrancy - amount)
}

// Rotate returns a new Tone with hue rotated by degrees.
// Accepts negative values (counter-clockwise).
func (t Tone) Rotate(degrees float64) Tone {
	return t.WithHue(space.NormalizeHue(t.hue + degrees))
}

// Complement returns the Tone directly opposite on the hue wheel (+180°).
func (t Tone) Complement() Tone {
	return t.Rotate(180).WithName("").WithMood("")
}
