package figma

import "testing"

func TestEffectValue(t *testing.T) {
	var effect Effect
	var ans string
	var want string

	isVisible := true

	effect = Effect{
		Type:    EffectTypeDropShadow,
		Visible: &isVisible,
		Radius:  4.0,
		Color: Color{
			Red:   0.1,
			Green: 0.2,
			Blue:  0.3,
			Alpha: 0.5,
		},
		BlendMode: BlendModeNormal,
		Offset: Vector{
			X: 0.0,
			Y: 4.0,
		},
		Spread: 0.0,
	}

	ans = effect.Value()
	want = "0px 4px 4px 0px rgba(25,51,76,0.5)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "DropSahdow", ans, want)
	}

	effect.Type = EffectTypeInnerShadow
	ans = effect.Value()
	want = "inset 0px 4px 4px 0px rgba(25,51,76,0.5)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "InnerSahdow", ans, want)
	}

	effect.Type = EffectTypeLayerBlur
	ans = effect.Value()
	want = "blur(4px)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "LayerBlur", ans, want)
	}

	effect.Type = EffectTypeBackgroundBlur
	ans = effect.Value()
	want = "blur(4px)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "BackgroundBlur", ans, want)
	}
}
