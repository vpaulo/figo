package figma

import (
	"maps"
	"slices"
	"testing"
)

func TestNodeBackground(t *testing.T) {
	var node Node
	var ans string
	var want string

	node = Node{
		Type: NodeTypeFrame,
		Fills: []Paint{
			{
				Type: PaintTypeSolid,
				Color: Color{
					Red:   0.1,
					Green: 0.2,
					Blue:  0.3,
					Alpha: 0.5,
				},
			},
		},
	}

	ans = node.Background()
	want = "rgba(25,51,76,0.5)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Background", ans, want)
	}
}

func TestNodeBorderColor(t *testing.T) {
	var node Node
	var ans string
	var want string

	node = Node{
		Type: NodeTypeFrame,
		Strokes: []Paint{
			{
				Type: PaintTypeSolid,
				Color: Color{
					Red:   0.1,
					Green: 0.2,
					Blue:  0.3,
					Alpha: 0.5,
				},
			},
		},
	}

	ans = node.BorderColor()
	want = "rgba(25,51,76,0.5)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "BorderColor", ans, want)
	}
}

func TestNodeBoxShadow(t *testing.T) {
	var node Node
	var ans string
	var want string

	isVisible := true

	node = Node{
		Type: NodeTypeFrame,
		Effects: []Effect{
			{
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
			},
		},
	}

	ans = node.BoxShadow()
	want = "0px 4px 4px 0px rgba(25,51,76,0.5)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "BoxShadow", ans, want)
	}
}

func TestNodeFont(t *testing.T) {
	var node Node
	var ans string
	var want string

	node = Node{
		Type: NodeTypeText,
		Style: TypeStyle{
			FontFamily:    "Roboto",
			FontSize:      12.0,
			FontWeight:    400.0,
			LineHeightPx:  16.0,
			LetterSpacing: 1.0,
		},
	}

	ans = node.Font()
	want = "font-family: Roboto;|font-size: 12px;|font-weight: 400;|line-height: 16px;|letter-spacing: 1px;"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Font", ans, want)
	}
}

func TestNodeSizes(t *testing.T) {
	var node Node
	var ans map[string]string
	var want map[string]string

	parent := Node{
		Type: NodeTypeFrame,
	}

	node = Node{
		Type:       NodeTypeFrame,
		LayoutMode: LayoutModeNone,
		MinWidth:   1.0,
		MaxWidth:   1.0,
		MinHeight:  1.0,
		MaxHeight:  1.0,
		AbsoluteBoundingBox: Rectangle{
			Width:  10.0,
			Height: 10.0,
		},
	}

	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":  "1px",
		"max-width":  "1px",
		"min-height": "1px",
		"max-height": "1px",
		"width":      "10px",
		"height":     "10px",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}

	parent.LayoutMode = LayoutModeHorizontal
	node.LayoutSizingHorizontal = LayoutSizingFixed
	node.LayoutSizingVertical = LayoutSizingFixed
	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":   "1px",
		"max-width":   "1px",
		"min-height":  "1px",
		"max-height":  "1px",
		"flex-shrink": "0",
		"width":       "10px",
		"height":      "10px",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}

	node.LayoutSizingHorizontal = LayoutSizingFill
	node.LayoutSizingVertical = LayoutSizingFill
	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":   "1px",
		"max-width":   "1px",
		"min-height":  "1px",
		"max-height":  "1px",
		"flex-shrink": "0",
		"align-self":  "stretch",
		"flex":        "1 0 0",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}

	node.LayoutAlign = LayoutAlignStretch
	node.LayoutGrow = 1.0
	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":  "1px",
		"max-width":  "1px",
		"min-height": "1px",
		"max-height": "1px",
		"align-self": "stretch",
		"flex":       "1 0 0",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}

	node.LayoutMode = LayoutModeHorizontal
	node.LayoutSizingHorizontal = LayoutSizingHug
	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":  "1px",
		"max-width":  "1px",
		"min-height": "1px",
		"max-height": "1px",
		"width":      "fit-content",
		"height":     "100%",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}

	node.LayoutSizingHorizontal = LayoutSizingFixed
	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":  "1px",
		"max-width":  "1px",
		"min-height": "1px",
		"max-height": "1px",
		"width":      "10px",
		"height":     "100%",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}

	node.LayoutSizingHorizontal = LayoutSizingFill
	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":  "1px",
		"max-width":  "1px",
		"min-height": "1px",
		"max-height": "1px",
		"width":      "100%",
		"height":     "100%",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}

	node.LayoutSizingVertical = LayoutSizingHug
	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":  "1px",
		"max-width":  "1px",
		"min-height": "1px",
		"max-height": "1px",
		"width":      "100%",
		"height":     "fit-content",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}

	node.LayoutSizingVertical = LayoutSizingFixed
	ans = node.Sizes(parent)
	want = map[string]string{
		"min-width":  "1px",
		"max-width":  "1px",
		"min-height": "1px",
		"max-height": "1px",
		"width":      "100%",
		"height":     "10px",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Sizes", ans, want)
	}
}

func TestNodeAlignment(t *testing.T) {
	var node Node
	var ans map[string]string
	var want map[string]string

	node = Node{
		Type:                  NodeTypeFrame,
		CounterAxisAlignItems: AlignItemsCenter,
	}

	ans = node.Alignment()
	want = map[string]string{
		"align-items":     "center",
		"justify-content": "flex-start",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}

	node.CounterAxisAlignItems = AlignItemsMax
	ans = node.Alignment()
	want = map[string]string{
		"align-items":     "flex-end",
		"justify-content": "flex-start",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}

	node.CounterAxisAlignItems = AlignItemsSpaceBetween
	ans = node.Alignment()
	want = map[string]string{
		"justify-content": "flex-start",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}

	node.CounterAxisAlignItems = AlignItemsBaseline
	ans = node.Alignment()
	want = map[string]string{
		"align-items":     "baseline",
		"justify-content": "flex-start",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}

	node.CounterAxisAlignItems = AlignItemsMin
	ans = node.Alignment()
	want = map[string]string{
		"align-items":     "flex-start",
		"justify-content": "flex-start",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}

	node.PrimaryAxisAlignItems = AlignItemsCenter
	ans = node.Alignment()
	want = map[string]string{
		"align-items":     "flex-start",
		"justify-content": "center",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}

	node.PrimaryAxisAlignItems = AlignItemsMax
	ans = node.Alignment()
	want = map[string]string{
		"align-items":     "flex-start",
		"justify-content": "flex-end",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}

	node.PrimaryAxisAlignItems = AlignItemsSpaceBetween
	ans = node.Alignment()
	want = map[string]string{
		"align-items":     "flex-start",
		"justify-content": "space-between",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}

	node.CounterAxisAlignContent = AlignContentSpaceBetween
	node.LayoutWrap = LayoutWrapWrap
	ans = node.Alignment()
	want = map[string]string{
		"align-items":     "flex-start",
		"justify-content": "space-between",
		"align-content":   "space-between",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Alignment", ans, want)
	}
}

func TestNodePadding(t *testing.T) {
	var node Node
	var ans string
	var want string

	node = Node{
		Type:          NodeTypeFrame,
		PaddingTop:    1.0,
		PaddingRight:  1.0,
		PaddingBottom: 1.0,
		PaddingLeft:   1.0,
	}

	ans = node.Padding()
	want = "1px"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Padding", ans, want)
	}

	node.PaddingRight = 2.0
	node.PaddingLeft = 2.0
	ans = node.Padding()
	want = "1px 2px"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Padding", ans, want)
	}

	node.PaddingBottom = 3.0
	ans = node.Padding()
	want = "1px 2px 3px"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Padding", ans, want)
	}

	node.PaddingLeft = 4.0
	ans = node.Padding()
	want = "1px 2px 3px 4px"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Padding", ans, want)
	}
}

func TestNodeBorderRadius(t *testing.T) {
	var node Node
	var ans string
	var want string

	node = Node{
		Type:         NodeTypeFrame,
		CornerRadius: 2.0,
	}

	ans = node.BorderRadius()
	want = "2px"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Border Radius", ans, want)
	}
}

func TestNodeIndividualBorderRadius(t *testing.T) {
	var node Node
	var ans string
	var want string

	node = Node{
		Type:                 NodeTypeFrame,
		RectangleCornerRadii: []float64{1.0, 2.0, 1.0, 2.0},
	}

	ans = node.BorderRadius()
	want = "1px 2px"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Individual Border Radius", ans, want)
	}

	node.RectangleCornerRadii = []float64{1.0, 2.0, 3.0, 2.0}
	ans = node.BorderRadius()
	want = "1px 2px 3px"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Individual Border Radius", ans, want)
	}

	node.RectangleCornerRadii = []float64{1.0, 2.0, 3.0, 4.0}
	ans = node.BorderRadius()
	want = "1px 2px 3px 4px"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Individual Border Radius", ans, want)
	}
}

func TestNodeBorder(t *testing.T) {
	var node Node
	var ans map[string]string
	var want map[string]string

	node = Node{
		Type: NodeTypeFrame,
		Strokes: []Paint{
			{
				Type: PaintTypeSolid,
				Color: Color{
					Red:   0.1,
					Green: 0.2,
					Blue:  0.3,
					Alpha: 0.5,
				},
			},
		},
		StrokeDashes: []float64{},
		StrokeWeight: 2.0,
	}

	ans = node.Border()
	want = map[string]string{"border": "2px solid rgba(25,51,76,0.5)"}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Border", ans, want)
	}

	node.StrokeDashes = []float64{1.0, 2.0}
	ans = node.Border()
	want = map[string]string{"border": "2px dashed rgba(25,51,76,0.5)"}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Border", ans, want)
	}
}

func TestNodeindividualBorders(t *testing.T) {
	var node Node
	var ans map[string]string
	var want map[string]string

	node = Node{
		Type: NodeTypeFrame,
		Strokes: []Paint{
			{
				Type: PaintTypeSolid,
				Color: Color{
					Red:   0.1,
					Green: 0.2,
					Blue:  0.3,
					Alpha: 0.5,
				},
			},
		},
		IndividualStrokeWeights: StrokeWeights{
			Top:    1.0,
			Right:  2.0,
			Bottom: 3.0,
			Left:   4.0,
		},
	}

	ans = node.Border()
	want = map[string]string{
		"border-top":    "1px solid rgba(25,51,76,0.5)",
		"border-right":  "2px solid rgba(25,51,76,0.5)",
		"border-bottom": "3px solid rgba(25,51,76,0.5)",
		"border-left":   "4px solid rgba(25,51,76,0.5)",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Individual Borders", ans, want)
	}
}

func TestNodeBorderStyle(t *testing.T) {
	var node Node
	var ans string
	var want string

	dashes := []float64{}
	node = Node{
		Type:         NodeTypeFrame,
		StrokeDashes: dashes,
	}

	ans = node.BorderStyle()
	want = "solid"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "BorderStyle", ans, want)
	}

	node.StrokeDashes = []float64{1.0, 2.0}
	ans = node.BorderStyle()
	want = "dashed"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "BorderStyle", ans, want)
	}
}

func TestNodeBlur(t *testing.T) {
	var node Node
	var ans string
	var want string

	isVisible := true
	node = Node{
		Type: NodeTypeFrame,
		Effects: []Effect{
			{
				Type:    EffectTypeLayerBlur,
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
			},
		},
	}

	ans = node.Blur()
	want = "blur(4px)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "Blur", ans, want)
	}
}

func TestNodeBackgroundBlur(t *testing.T) {
	var node Node
	var ans string
	var want string

	isVisible := true
	node = Node{
		Type: NodeTypeFrame,
		Effects: []Effect{
			{
				Type:    EffectTypeBackgroundBlur,
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
			},
		},
	}

	ans = node.BackgroundBlur()
	want = "blur(4px)"
	if ans != want {
		t.Errorf("%+v = %v; want %v", "BackgroundBlur", ans, want)
	}
}

func TestNodeCss(t *testing.T) {
	var node Node
	var ans map[string]string
	var want map[string]string

	isVisible := true
	parent := Node{
		Type: NodeTypeFrame,
	}
	node = Node{
		Type:     NodeTypeFrame,
		Visible:  &isVisible,
		Rotation: 0.7853982,
		Fills: []Paint{
			{
				Type: PaintTypeSolid,
				Color: Color{
					Red:   1.0,
					Green: 1.0,
					Blue:  1.0,
					Alpha: 1.0,
				},
			},
		},
		Strokes: []Paint{
			{
				Type: PaintTypeSolid,
				Color: Color{
					Red:   0.0,
					Green: 0.0,
					Blue:  0.0,
					Alpha: 1.0,
				},
			},
		},
		StrokeWeight:         1.0,
		RectangleCornerRadii: []float64{3.0, 5.0, 3.0, 5.0},
		AbsoluteBoundingBox: Rectangle{
			Width:  50.0,
			Height: 56.0,
		},
		PaddingLeft:   6.0,
		PaddingRight:  6.0,
		PaddingTop:    8.0,
		PaddingBottom: 8.0,
		ItemSpacing:   4.0,
	}

	ans = node.Css(parent)
	want = map[string]string{
		"background":    "rgba(255,255,255,1)",
		"border":        "1px solid rgba(0,0,0,1)",
		"border-radius": "3px 5px",
		"transform":     "rotate(45deg)",
	}
	if !maps.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Css", ans, want)
	}
}

// TODO TextCss
// TODO Variants

// name -> .name
// property=value -> [property="value"]
// property=pseudo -> :pseudo
// property=value;pseudo -> [property="value"]:pseudo
// property=value;class -> [property="value"].class
func TestNodeClasses(t *testing.T) {
	var node Node
	var ans []string
	var want []string

	node = Node{
		Type: NodeTypeFrame,
		Name: "my-component",
	}

	ans = node.Classes()
	want = []string{".my-component"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "myComponent"
	ans = node.Classes()
	want = []string{".my-component"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=default"
	ans = node.Classes()
	want = []string{}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=test"
	ans = node.Classes()
	want = []string{"[type=\"test\"]"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=hover"
	ans = node.Classes()
	want = []string{":hover"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=test,state=ok"
	ans = node.Classes()
	want = []string{"[type=\"test\"]", "[state=\"ok\"]"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=test,state=hover"
	ans = node.Classes()
	want = []string{"[type=\"test\"]", ":hover"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "state=hover, type=test"
	ans = node.Classes()
	want = []string{"[type=\"test\"]", ":hover"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "state=default;hover, type=test"
	ans = node.Classes()
	want = []string{"[type=\"test\"]", ":hover"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=test,state=default"
	ans = node.Classes()
	want = []string{"[type=\"test\"]"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=test;hover"
	ans = node.Classes()
	want = []string{"[type=\"test\"]", ":hover"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=test;ok"
	ans = node.Classes()
	want = []string{"[type=\"test\"].ok"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=test,hover=true"
	ans = node.Classes()
	want = []string{"[type=\"test\"]", ":hover"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "type=test,hover=false"
	ans = node.Classes()
	want = []string{"[type=\"test\"]"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "hover=true"
	ans = node.Classes()
	want = []string{":hover"}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}

	node.Name = "hover=false"
	ans = node.Classes()
	want = []string{}
	if !slices.Equal(ans, want) {
		t.Errorf("%+v = %v; want %v", "Classes", ans, want)
	}
}
