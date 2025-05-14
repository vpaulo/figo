package figma

import (
	"fmt"
	"strings"
)

func (n *Node) IsComponentSet() bool {
	return n.Type == NodeTypeComponentSet
}

func (n *Node) IsComponentOrSet() bool {
	return n.Type == NodeTypeComponentSet || n.Type == NodeTypeComponent
}

func (n *Node) IsInstance() bool {
	return n.Type == NodeTypeInstance
}

func (n *Node) IsText() bool {
	return n.Type == NodeTypeText
}

func (n *Node) IsFrame() bool {
	return n.Type == NodeTypeComponentSet ||
		n.Type == NodeTypeComponent ||
		n.Type == NodeTypeInstance ||
		n.Type == NodeTypeFrame ||
		n.Type == NodeTypeGroup
}

func (n *Node) Background() string {
	color := ""
	for _, fill := range n.Fills {
		switch fill.Type {
		case PaintTypeSolid:
			color = fill.Color.Rgba()
		case PaintTypeGradientLinear:
			color = ""
		case PaintTypeGradientAngular:
			color = ""
		case PaintTypeGradientRadial:
			color = ""
		case PaintTypeGradientDiamond:
			color = ""
		}
	}
	return color
}

func (n *Node) BorderColor() string {
	color := ""
	for _, fill := range n.Strokes {
		switch fill.Type {
		case PaintTypeSolid:
			color = fill.Color.Rgba()
			// case PaintTypeGradientLinear:
			// 	background = ""
			// case PaintTypeGradientAngular:
			// 	background = ""
			// case PaintTypeGradientRadial:
			// 	background = ""
			// case PaintTypeGradientDiamond:
			// 	background = ""
		}
	}
	return color
}

func (n *Node) BoxShadow() string {
	var value []string

	for _, effect := range n.Effects {
		if effect.Visible {
			value = append(value, effect.Value())
		}
	}

	return strings.Join(value, ", ")
}

func (n *Node) Font() string {
	var value []string

	if n.Style.FontFamily != "" {
		value = append(value, fmt.Sprintf("font-family: %v;", n.Style.FontFamily))
	}

	if n.Style.FontSize != 0.0 {
		value = append(value, fmt.Sprintf("font-size: %vpx;", int(n.Style.FontSize)))
	}

	if n.Style.FontWeight != 0.0 {
		value = append(value, fmt.Sprintf("font-weight: %v;", int(n.Style.FontWeight)))
	}

	if n.Style.LineHeightPx != 0.0 {
		value = append(value, fmt.Sprintf("line-height: %vpx;", int(n.Style.LineHeightPx)))
	}

	if n.Style.LetterSpacing != 0.0 {
		value = append(value, fmt.Sprintf("letter-spacing: %vpx;", int(n.Style.LetterSpacing)))
	}
	return strings.Join(value, "|")
}
