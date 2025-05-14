package figma

import "fmt"

func (e *Effect) Value() string {
	value := ""
	switch e.Type {
	case EffectTypeInnerShadow:
		value = fmt.Sprintf("inset %v", e.Shadow())
	case EffectTypeDropShadow:
		value = e.Shadow()
	case EffectTypeLayerBlur:
		value = ""
	case EffectTypeBackgroundBlur:
		value = ""
	}
	return value
}

func (e *Effect) Shadow() string {
	x := int(e.Offset.X)
	y := int(e.Offset.Y)
	radius := int(e.Radius)
	spread := int(e.Spread)
	color := e.Color.Rgba()

	return fmt.Sprintf("%vpx %vpx %vpx %vpx %v", x, y, radius, spread, color)
}
