package figma

import (
	"fmt"
	"slices"
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

func (n *Node) IsVector() bool {
	return n.Type == NodeTypeVector
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

	switch n.Style.TextAlignHorizontal {
	case TextAlignRight:
		rules["text-align"] = "right"
	case TextAlignCenter:
		rules["text-align"] = "center"
	case TextAlignJustified:
		rules["text-align"] = "justify"
	}

	switch n.Style.TextDecoration {
	case TextDecorationStrikethrough:
		rules["text-decoration-line"] = "strikethrough"
	case TextDecorationUnderline:
		rules["text-decoration-line"] = "underline"
	}

	switch n.Style.TextCase {
	case TextCaseUpper:
		rules["text-transform"] = "uppercase"
	case TextCaseLower:
		rules["text-transform"] = "lowercase"
	case TextCaseTitle:
		rules["text-transform"] = "capitalize"
	case TextCaseSmallCaps:
		rules["font-variant"] = "small-caps"
	case TextCaseSmallCapsForced:
		rules["font-variant"] = "all-small-caps"
	}

	if n.Style.TextTruncation != "" {
		rules["text-overflow"] = "ellipsis"

		if n.Style.MaxLines != 0.0 {
			rules["-webkit-box-orient"] = "vertical"
			rules["-webkit-line-clamp"] = fmt.Sprintf("%v", int(n.Style.MaxLines))
		}
	}

	if n.Background() != "" {
		rules["color"] = n.Background()
	}

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

	if n.LayoutSizingHorizontal == LayoutSizingFixed {
		rules["width"] = fmt.Sprintf("%vpx", int(n.AbsoluteBoundingBox.Width))
	}

	if n.LayoutSizingHorizontal == LayoutSizingHug {
		rules["width"] = "fit-content"
	}

	if n.LayoutSizingHorizontal == LayoutSizingFill {
		rules["width"] = "100%"
	}

	if n.LayoutSizingVertical == LayoutSizingFixed {
		rules["height"] = fmt.Sprintf("%vpx", int(n.AbsoluteBoundingBox.Width))
	}

	if n.LayoutSizingVertical == LayoutSizingHug {
		rules["height"] = "fit-content"
	}

	if n.LayoutSizingVertical == LayoutSizingFill {
		rules["height"] = "100%"
	}

	return rules
}

func (n *Node) Variants() []Variant {
	var variants []Variant

	for key, def := range n.ComponentPropertyDefinitions {
		if def.Type == ComponentPropertyTypeVariant {
			variant := Variant{
				Name:    key,
				Value:   def.DefaultValue,
				Options: def.VariantOptions,
			}

			variants = append(variants, variant)
		}
	}

	return variants
}

func (n *Node) Classes() string {
	var classes []string

	if strings.Contains(n.Name, ",") {
		variants := strings.Split(n.Name, ",")
		var pseudoArr []string

		for i, variant := range variants {
			variants[i] = strings.TrimSpace(variant)

			attribute := getVariantClasses(variants[i])
			if attribute != "" {
				classes = append(classes, attribute)
			}

			pseudo := getPseudoClasses(variants[i])
			if pseudo != "" {
				pseudoArr = append(pseudoArr, pseudo)
			}
		}
		classes = append(classes, pseudoArr...)
	} else if strings.Contains(n.Name, "=") {
		attribute := getVariantClasses(n.Name)
		pseudo := getPseudoClasses(n.Name)
		if attribute != "" {
			classes = append(classes, attribute)
		}
		if pseudo != "" {
			classes = append(classes, pseudo)
		}
	} else {
		classes = []string{fmt.Sprintf(".%v", ToKebabCase(n.Name))}
	}

	return strings.Join(classes, "")
}

var pseudoClasses = []string{
	"hover",
	"active",
	"focus",
	"disabled",
	"focus-visible",
	"focus-within",
}

func getVariantClasses(variant string) string {
	parts := strings.Split(variant, "=")
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		attribute := ToKebabCase(parts[0])
		value := ToKebabCase(parts[1])

		if slices.Contains(pseudoClasses, attribute) {
			return ""
		}

		valueParts := strings.Split(value, ";")
		if len(valueParts) == 2 && valueParts[0] != "" && valueParts[1] != "" {
			val := ToKebabCase(valueParts[0])
			second := ToKebabCase(valueParts[1])
			cl := ""

			if second != "default" && !slices.Contains(pseudoClasses, second) {
				cl = fmt.Sprintf(".%v", second)
			}

			if val == "default" {
				return fmt.Sprintf("%v", cl)
			}

			return fmt.Sprintf("[%v=\"%v\"]%v", attribute, val, cl)
		}

		if value != "default" && !slices.Contains(pseudoClasses, value) {
			return fmt.Sprintf("[%v=\"%v\"]", attribute, value)
		}
	}

	return ""
}

func getPseudoClasses(variant string) string {
	parts := strings.Split(variant, "=")
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		attribute := ToKebabCase(parts[0])
		value := ToKebabCase(parts[1])

		valueParts := strings.Split(value, ";")
		if len(valueParts) == 2 && valueParts[1] != "" {
			pseudo := ToKebabCase(valueParts[1])

			if pseudo != "default" && slices.Contains(pseudoClasses, pseudo) {
				return fmt.Sprintf(":%v", pseudo)
			}

			return ""
		}

		if value != "default" && slices.Contains(pseudoClasses, value) {
			return fmt.Sprintf(":%v", value)
		}

		if slices.Contains(pseudoClasses, attribute) && value == "true" {
			return fmt.Sprintf(":%v", attribute)
		}
	}

	return ""
}
