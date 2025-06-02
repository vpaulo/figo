package figma

import (
	"fmt"
	"strings"
)

func (n *Node) IsComponentSet() bool {
	return n.Type == NodeTypeComponentSet
}

func (n *Node) IsComponent() bool {
	return n.Type == NodeTypeComponent
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

func (n *Node) IsAutoLayout() bool {
	return n.LayoutMode == LayoutModeHorizontal || n.LayoutMode == LayoutModeVertical
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
		if *effect.Visible {
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

func (n *Node) Css(parent Node) map[string]string {
	rules := make(map[string]string)

	if !*n.Visible {
		rules["display"] = "none"
	}

	if n.ClipsContent {
		rules["overflow"] = "hidden"
	}

	for key, value := range n.Sizes(parent) {
		rules[key] = value
	}

	if n.IsAutoLayout() {
		if *n.Visible {
			rules["display"] = "flex"
		}
		if n.LayoutWrap == LayoutWrapWrap {
			rules["flex-wrap"] = "wrap"
		}
		if n.LayoutMode == LayoutModeVertical {
			rules["flex-direction"] = "column"
		}

		for key, value := range n.Alignment() {
			rules[key] = value
		}

		if n.ItemSpacing != 0.0 {
			rules["gap"] = fmt.Sprintf("%vpx", int(n.ItemSpacing))
		}

		if n.Padding() != "" {
			rules["padding"] = n.Padding()
		}
	}

	// Rotation only works well for 90 * n degrees, for other values like 45deg figma changes the sizes of width and height.
	if n.Rotation != 0.0 {
		rules["transform"] = fmt.Sprintf("rotate(%vdeg)", ToDegrees(n.Rotation))
	}

	if n.BorderRadius() != "" {
		rules["border-radius"] = n.BorderRadius()
	}

	for key, value := range n.Border() {
		rules[key] = value
	}

	if n.Background() != "" {
		background := n.Background()

		// TODO: Get token for background if exixts
		// if n.Styles["fill"] != "" {
		// 	background = fmt.Sprintf("var(%v)", tokens[n.Styles["fill"]])
		// }

		rules["background"] = background
	}

	if n.BoxShadow() != "" {
		boxShadow := n.BoxShadow()

		// TODO: Get token for box shadow if exixts
		// if n.Styles["effect"] != "" {
		// 	boxShadow = fmt.Sprintf("var(%v)", tokens[n.Styles["effect"]])
		// }

		rules["box-shadow"] = boxShadow
	}

	if n.Blur() != "" {
		rules["filter"] = n.Blur()
	}

	if n.BackgroundBlur() != "" {
		rules["backdrop-filter"] = n.BackgroundBlur()
	}

	return rules
}

func (n *Node) Sizes(parent Node) map[string]string {
	rules := make(map[string]string)

	if n.MinWidth != 0.0 {
		rules["min-width"] = fmt.Sprintf("%vpx", int(n.MinWidth))
	}
	if n.MaxWidth != 0.0 {
		rules["max-width"] = fmt.Sprintf("%vpx", int(n.MinWidth))
	}
	if n.MinHeight != 0.0 {
		rules["min-height"] = fmt.Sprintf("%vpx", int(n.MinWidth))
	}
	if n.MaxHeight != 0.0 {
		rules["max-height"] = fmt.Sprintf("%vpx", int(n.MinWidth))
	}

	if n.LayoutMode == LayoutModeNone {
		if parent.IsAutoLayout() {
			if n.LayoutSizingHorizontal == LayoutSizingFixed && n.AbsoluteBoundingBox.Width != 0.0 {
				rules["width"] = fmt.Sprintf("%vpx", int(n.AbsoluteBoundingBox.Width))
			}
			if n.LayoutGrow == 0.0 {
				rules["flex-shrink"] = fmt.Sprintf("%v", 0)
			}
			if n.LayoutSizingHorizontal == LayoutSizingFill {
				if n.LayoutAlign == LayoutAlignStretch {
					rules["align-self"] = "stretch"
				} else {
					rules["flex"] = "1 0 0"
				}
			}
			if n.LayoutSizingVertical == LayoutSizingFixed && n.AbsoluteBoundingBox.Height != 0.0 {
				rules["height"] = fmt.Sprintf("%vpx", int(n.AbsoluteBoundingBox.Height))
			}
			if n.LayoutSizingVertical == LayoutSizingFill {
				if n.LayoutGrow == 1.0 {
					rules["flex"] = "1 0 0"
				} else {
					rules["align-self"] = "stretch"
				}
			}
		} else {
			if n.AbsoluteBoundingBox.Width != 0.0 {
				rules["width"] = fmt.Sprintf("%vpx", int(n.AbsoluteBoundingBox.Width))
			}
			if n.AbsoluteBoundingBox.Height != 0.0 {
				rules["height"] = fmt.Sprintf("%vpx", int(n.AbsoluteBoundingBox.Height))
			}
		}
	} else if n.IsAutoLayout() {
		if n.LayoutSizingHorizontal == LayoutSizingHug {
			rules["width"] = "fit-content"
		}
		if n.LayoutSizingHorizontal == LayoutSizingFixed {
			rules["width"] = fmt.Sprintf("%vpx", int(n.AbsoluteBoundingBox.Width))
		}
		if n.LayoutSizingHorizontal == LayoutSizingFill {
			rules["width"] = "100%"
		}

		if n.LayoutSizingVertical == LayoutSizingHug {
			rules["height"] = "fit-content"
		}
		if n.LayoutSizingVertical == LayoutSizingFixed {
			rules["height"] = fmt.Sprintf("%vpx", int(n.AbsoluteBoundingBox.Height))
		}
		if n.LayoutSizingVertical == LayoutSizingFill {
			rules["height"] = "100%"
		}
	}

	return rules
}

func (n *Node) Alignment() map[string]string {
	rules := make(map[string]string)

	switch n.CounterAxisAlignItems {
	case AlignItemsCenter:
		rules["align-items"] = "center"
	case AlignItemsMax:
		rules["align-items"] = "flex-end"
	case AlignItemsSpaceBetween:
		// align items does not have space between
	case AlignItemsBaseline:
		rules["align-items"] = "baseline"
	default:
		rules["align-items"] = "flex-start" // Default AlignItemsMin
	}

	switch n.PrimaryAxisAlignItems {
	case AlignItemsCenter:
		rules["justify-content"] = "center"
	case AlignItemsMax:
		rules["justify-content"] = "flex-end"
	case AlignItemsSpaceBetween:
		rules["justify-content"] = "space-between"
	default:
		rules["justify-content"] = "flex-start" // Default AlignItemsMin
	}

	if n.CounterAxisAlignContent == AlignContentSpaceBetween && n.LayoutWrap == LayoutWrapWrap {
		rules["align-content"] = "space-between"
	}

	return rules
}

func (n *Node) Padding() string {
	// TODO: remove px for when value is 0
	top := n.PaddingTop
	right := n.PaddingRight
	bottom := n.PaddingBottom
	left := n.PaddingLeft

	if top == 0.0 && right == 0.0 && bottom == 0.0 && left == 0.0 {
		return ""
	}

	if top == bottom && right == left && top == right {
		return fmt.Sprintf("%vpx", int(top))
	} else if top == bottom && right == left {
		return fmt.Sprintf("%vpx %vpx", int(top), int(right))
	} else if right == left {
		return fmt.Sprintf("%vpx %vpx %vpx", int(top), int(right), int(bottom))
	} else {
		return fmt.Sprintf("%vpx %vpx %vpx %vpx", int(top), int(right), int(bottom), int(left))
	}
}

func (n *Node) BorderRadius() string {
	if n.CornerRadius != 0.0 {
		return fmt.Sprintf("%vpx", int(n.CornerRadius))
	}

	if len(n.RectangleCornerRadii) > 0 {
		topLeft := n.RectangleCornerRadii[0]
		topRight := n.RectangleCornerRadii[1]
		bottomRight := n.RectangleCornerRadii[2]
		bottomLeft := n.RectangleCornerRadii[3]

		if topLeft == bottomRight && topRight == bottomLeft {
			return fmt.Sprintf("%vpx %vpx", int(topLeft), int(topRight))
		} else if topRight == bottomLeft {
			return fmt.Sprintf("%vpx %vpx %vpx", int(topLeft), int(topRight), int(bottomRight))
		} else {
			return fmt.Sprintf("%vpx %vpx %vpx %vpx", int(topLeft), int(topRight), int(bottomRight), int(bottomLeft))
		}
	}

	return ""
}

func (n *Node) Border() map[string]string {
	// TODO: when multiple colours and sizes convert into "border-width", "border-color" and "border-style"
	rules := make(map[string]string)

	style := n.BorderStyle()
	color := n.BorderColor()
	width := ""

	if n.StrokeWeight != 0.0 {
		width = fmt.Sprintf("%vpx", int(n.StrokeWeight))
	}

	if width != "" && color != "" {
		rules["border"] = fmt.Sprintf("%v %v %v", width, style, color)
	} else if color != "" && n.IndividualStrokeWeights != (StrokeWeights{}) {
		if n.IndividualStrokeWeights.Top > 0.0 {
			rules["border-top"] = fmt.Sprintf("%vpx %v %v", int(n.IndividualStrokeWeights.Top), style, color)
		}
		if n.IndividualStrokeWeights.Right > 0.0 {
			rules["border-right"] = fmt.Sprintf("%vpx %v %v", int(n.IndividualStrokeWeights.Right), style, color)
		}
		if n.IndividualStrokeWeights.Bottom > 0.0 {
			rules["border-bottom"] = fmt.Sprintf("%vpx %v %v", int(n.IndividualStrokeWeights.Bottom), style, color)
		}
		if n.IndividualStrokeWeights.Left > 0.0 {
			rules["border-left"] = fmt.Sprintf("%vpx %v %v", int(n.IndividualStrokeWeights.Left), style, color)
		}
	}

	return rules
}

func (n *Node) BorderStyle() string {
	if len(n.StrokeDashes) > 0 {
		return "dashed"
	}
	return "solid"
}

func (n *Node) Blur() string {
	blur := ""
	for _, effect := range n.Effects {
		if *effect.Visible && effect.Type == EffectTypeLayerBlur {
			blur = effect.Value()
		}
	}
	return blur
}

func (n *Node) BackgroundBlur() string {
	blur := ""
	for _, effect := range n.Effects {
		if *effect.Visible && effect.Type == EffectTypeBackgroundBlur {
			blur = effect.Value()
		}
	}
	return blur
}

func (n *Node) TextCss() map[string]string {
	rules := make(map[string]string)

	if n.Style.FontFamily != "" {
		rules["font-family"] = n.Style.FontFamily
	}

	if n.Style.FontSize != 0.0 {
		rules["font-size"] = fmt.Sprintf("%vpx", int(n.Style.FontSize))
	}

	if n.Style.FontWeight != 0.0 {
		rules["font-weight"] = fmt.Sprintf("%v", int(n.Style.FontWeight))
	}

	if n.Style.LineHeightPx != 0.0 {
		rules["line-height"] = fmt.Sprintf("%vpx", int(n.Style.LineHeightPx))
	}

	if n.Style.LetterSpacing != 0.0 {
		rules["letter-spacing"] = fmt.Sprintf("%vpx", int(n.Style.LetterSpacing))
	}

	// if !self.text_colour().is_empty() {
	//            rules.insert("color".to_string(), self.text_colour());
	//        }

	//        if !style.font_family.is_empty() {
	//            rules.insert("font-family".to_string(), style.font_family.to_string());
	//        }

	//        if style.font_size != 0.0 {
	//            rules.insert("font-size".to_string(), format!("{:.0}px", style.font_size));
	//        }

	//        if style.font_weight != 0.0 {
	//            rules.insert(
	//                "font-weight".to_string(),
	//                format!("{:.0}", style.font_weight),
	//            );
	//        }

	//        if style.line_height() > 0.0 {
	//            rules.insert(
	//                "line-height".to_string(),
	//                format!("{}", style.line_height()),
	//            );
	//        }

	//        if style.letter_spacing != 0.0 {
	//            rules.insert(
	//                "letter-spacing".to_string(),
	//                format!("{:.0}px", style.letter_spacing),
	//            );
	//        }

	//        if !self.sizes().is_empty() {
	//            for (key, value) in self.sizes().iter() {
	//                rules.insert(key.to_string(), value.to_string());
	//            }
	//        }

	//        if !style.text_align().is_empty() {
	//            rules.insert("text-align".to_string(), format!("{}", style.text_align()));
	//        }

	//        if !style.text_decoration().is_empty() {
	//            rules.insert(
	//                "text-decoration-line".to_string(),
	//                format!("{}", style.text_decoration()),
	//            );
	//        }

	//        if !style.text_transform().is_empty() {
	//            rules.insert(
	//                "text-transform".to_string(),
	//                format!("{}", style.text_transform()),
	//            );
	//        }

	//        if !style.font_variant().is_empty() {
	//            rules.insert(
	//                "font-variant".to_string(),
	//                format!("{}", style.font_variant()),
	//            );
	//        }

	//        if style.text_truncation == TextTruncation::Ending {
	//            rules.insert("text-overflow".to_string(), "ellipsis".to_string());

	//            if let Some(max) = style.max_lines {
	//                rules.insert("-webkit-box-orient".to_string(), "vertical".to_string());
	//                rules.insert("-webkit-line-clamp".to_string(), format!("{:.0}", max));
	//            }
	//        }
	return rules
}
