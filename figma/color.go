package figma

import (
	"fmt"
	"math"
)

func (c *Color) IsTransparent() bool {
	return c.Alpha == 0.0
}

func (c *Color) Rgba() string {
	red := c.Red * 255.0
	green := c.Green * 255.0
	blue := c.Blue * 255.0
	alpha := c.Alpha

	return fmt.Sprintf("rgba(%v,%v,%v,%v)", int(red), int(green), int(blue), alpha)
}

func (c *Color) Hsl() string {
	rf := c.Red
	gf := c.Green
	bf := c.Blue

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	h, s, l := 0.0, 0.0, (max+min)/2

	if max != min {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case rf:
			h = (gf - bf) / d
			if gf < bf {
				h += 6
			}
		case gf:
			h = (bf-rf)/d + 2
		case bf:
			h = (rf-gf)/d + 4
		}
		h /= 6
	}

	return fmt.Sprintf("hsl(%v,%v%%,%v%%)", math.Round(h*360), math.Round(s*100), math.Round(l*100))
}
