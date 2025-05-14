package figma

import "testing"

func TestTransparency(t *testing.T) {
	var color Color
	var got bool

	color = Color{
		Red:   0.0,
		Green: 0.0,
		Blue:  0.0,
		Alpha: 0.0,
	}

	got = color.IsTransparent()
	if !got {
		t.Errorf("%+v = %v; want true", color, got)
	}

	color.Alpha = 0.1

	got = color.IsTransparent()
	if got {
		t.Errorf("%+v = %v; want false", color, got)
	}
}

func TestRgba(t *testing.T) {
	var color Color
	var ans string
	var want string

	color = Color{
		Red:   0.0,
		Green: 0.0,
		Blue:  0.0,
		Alpha: 1,
	}

	ans = color.Rgba()
	want = "rgba(0,0,0,1)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", color, ans, want)
	}

	color = Color{
		Red:   0.1,
		Green: 0.2,
		Blue:  0.3,
		Alpha: 0.5,
	}

	ans = color.Rgba()
	want = "rgba(25,51,76,0.5)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", color, ans, want)
	}
}

func TestHsl(t *testing.T) {
	var color Color
	var ans string
	var want string

	color = Color{
		Red:   0.0,
		Green: 0.0,
		Blue:  0.0,
		Alpha: 1,
	}

	ans = color.Hsl()
	want = "hsl(0,0%,0%)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", color, ans, want)
	}

	color = Color{
		Red:   0.1,
		Green: 0.2,
		Blue:  0.3,
		Alpha: 0.5,
	}

	ans = color.Hsl()
	want = "hsl(210,50%,20%)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", color, ans, want)
	}
}
