package figma

import (
	"testing"
)

func TestRoundToDecimals(t *testing.T) {
	var ans float64
	var want float64

	ans = RoundToDecimals(0.100023030345340303, 5)
	want = 0.10002
	if ans != want {
		t.Errorf("RoundToDecimals = %v; want %v", ans, want)
	}
}
