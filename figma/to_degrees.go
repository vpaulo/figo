package figma

import "math"

func ToDegrees(angle float64) int {
	degree := math.Pi / 180

	return int(angle / degree)
}
