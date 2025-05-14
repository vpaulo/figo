package figma

import "testing"

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

	node = Node{
		Type: NodeTypeFrame,
		Effects: []Effect{
			{
				Type:    EffectTypeDropShadow,
				Visible: true,
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
