package figma

import (
	"testing"
)

func TestToDegrees(t *testing.T) {
	var ans int
	var want int

	ans = ToDegrees(-1.5707964)
	want = -90
	if ans != want {
		t.Errorf("ToDegrees = %v; want %v", ans, want)
	}

	ans = ToDegrees(1.5707964)
	want = 90
	if ans != want {
		t.Errorf("ToDegrees = %v; want %v", ans, want)
	}

	ans = ToDegrees(-0.7853982)
	want = -45
	if ans != want {
		t.Errorf("ToDegrees = %v; want %v", ans, want)
	}

	ans = ToDegrees(0.7853982)
	want = 45
	if ans != want {
		t.Errorf("ToDegrees = %v; want %v", ans, want)
	}

	ans = ToDegrees(-5.551115e-17) // This number is very close to 0, so we will assume that it is 0
	want = 0
	if ans != want {
		t.Errorf("ToDegrees = %v; want %v", ans, want)
	}
}
