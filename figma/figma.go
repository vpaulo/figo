package figma

// Figma files types
type File struct {
	Name          string                  `json:"name"`
	LastModified  string                  `json:"lastModified"`
	Version       string                  `json:"version"`
	Document      Node                    `json:"document"`
	ComponentSets map[string]ComponentSet `json:"componentSets"`
	Components    map[string]Component    `json:"components"`
	Styles        map[string]Style        `json:"styles"`
}

type NodeType string

const (
	NodeTypeDocument       NodeType = "DOCUMENT"
	NodeTypeCanvas                  = "CANVAS"
	NodeTypeFrame                   = "FRAME"
	NodeTypeGroup                   = "GROUP"
	NodeTypeSection                 = "SECTION"
	NodeTypeVector                  = "VECTOR"
	NodeTypeBoolean                 = "BOOLEAN_OPERATION"
	NodeTypeStar                    = "STAR"
	NodeTypeLine                    = "LINE"
	NodeTypeEllipse                 = "ELLIPSE"
	NodeTypeRegularPolygon          = "REGULAR_POLYGON"
	NodeTypeRectangle               = "RECTANGLE"
	NodeTypeText                    = "TEXT"
	NodeTypeSlice                   = "SLICE"
	NodeTypeComponent               = "COMPONENT"
	NodeTypeComponentSet            = "COMPONENT_SET"
	NodeTypeInstance                = "INSTANCE"
)

type Node struct {
	// Node Common
	Type     NodeType `json:"type"`
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Visible  bool     `json:"visible,omitempty"`
	Children []Node   `json:"children"`
	// CANVAS
	BackgroundColor Color           `json:"backgroundColor"`
	ExportSettings  []ExportSetting `json:"exportSettings,omitempty"`
	// FRAME
	Fills                   []Paint           `json:"fills,omitempty"`
	Strokes                 []Paint           `json:"strokes,omitempty"`
	StrokeWeight            float64           `json:"strokeWeight,omitempty"`
	StrokeAlign             StrokeAlign       `json:"strokeAlign,omitempty"`
	StrokeDashes            []float64         `json:"strokeDashes,omitempty"`
	CornerRadius            float64           `json:"cornerRadius,omitempty"`
	RectangleCornerRadii    []float64         `json:"rectangleCornerRadii,omitempty"`
	CornerSmoothing         float64           `json:"cornerSmoothing,omitempty"`
	BlendMode               BlendMode         `json:"blendMode,omitempty"`
	PreserveRatio           bool              `json:"preserveRatio,omitempty"`
	TargetAspectRatio       Vector            `json:"targetAspectRatio"`
	Constraints             LayoutConstraint  `json:"constraints"`
	LayoutAlign             LayoutAlign       `json:"layoutAlign,omitempty"`
	Opacity                 float64           `json:"opacity,omitempty"`
	AbsoluteBoundingBox     Rectangle         `json:"absoluteBoundingBox"`
	AbsoluteRenderBounds    Rectangle         `json:"absoluteRenderBounds"`
	Size                    Vector            `json:"size"`
	MinWidth                float64           `json:"minWidth,omitempty"`
	MaxWidth                float64           `json:"maxWidth,omitempty"`
	MinHeight               float64           `json:"minHeight,omitempty"`
	MaxHeight               float64           `json:"maxHeight,omitempty"`
	RelativeTransform       Transform         `json:"relativeTransform,omitempty"`
	ClipsContent            bool              `json:"clipsContent,omitempty"`
	LayoutMode              LayoutMode        `json:"layoutMode,omitempty"`
	LayoutSizingHorizontal  LayoutSizing      `json:"layoutSizingHorizontal,omitempty"`
	LayoutSizingVertical    LayoutSizing      `json:"layoutSizingVertical,omitempty"`
	LayoutWrap              LayoutWrap        `json:"layoutWrap,omitempty"`
	PrimaryAxisSizingMode   SizingMode        `json:"primaryAxisSizingMode,omitempty"`
	CounterAxisSizingMode   SizingMode        `json:"counterAxisSizingMode,omitempty"`
	PrimaryAxisAlignItems   AlignItems        `json:"primaryAxisAlignItems,omitempty"`
	CounterAxisAlignItems   AlignItems        `json:"counterAxisAlignItems,omitempty"`
	CounterAxisAlignContent AlignContent      `json:"counterAxisAlignContent,omitempty"`
	PaddingLeft             float64           `json:"paddingLeft,omitempty"`
	PaddingRight            float64           `json:"paddingRight,omitempty"`
	PaddingTop              float64           `json:"paddingTop,omitempty"`
	PaddingBottom           float64           `json:"paddingBottom,omitempty"`
	ItemSpacing             float64           `json:"itemSpacing,omitempty"`
	CounterAxisSpacing      float64           `json:"counterAxisSpacing,omitempty"`
	LayoutPositioning       LayoutPositioning `json:"layoutPositioning,omitempty"`
	ItemReverseZIndex       bool              `json:"itemReverseZIndex,omitempty"`
	StrokesIncludedInLayout bool              `json:"strokesIncludedInLayout,omitempty"`
	LayoutGrids             []LayoutGrid      `json:"layoutGrids,omitempty"`
	OverflowDirection       OverflowDirection `json:"overflowDirection,omitempty"`
	Effects                 []Effect          `json:"effects,omitempty"`
	IsMask                  bool              `json:"isMask,omitempty"`
	MaskType                MaskType          `json:"maskType,omitempty"`
	// TODO: this (StyleType) does not seem to match the values returned
	Styles map[StyleType]string `json:"styles,omitempty"`
	// SECTION
	SectionContentsHidden bool `json:"sectionContentsHidden,omitempty"`
	// VECTOR
	LayoutGrow              float64                   `json:"layoutGrow,omitempty"`
	FillGeometry            []Path                    `json:"fillGeometry,omitempty"`
	FillOverrideTable       map[float64]PaintOverride `json:"fillOverrideTable,omitempty"`
	IndividualStrokeWeights StrokeWeights             `json:"individualStrokeWeights"`
	StrokeCap               StrokeCap                 `json:"strokeCap,omitempty"`
	StrokeJoin              StrokeJoin                `json:"strokeJoin,omitempty"`
	StrokeMiterAngle        float64                   `json:"strokeMiterAngle,omitempty"`
	StrokeGeometry          []Path                    `json:"strokeGeometry,omitempty"`
	// BOOLEAN_OPERATION
	BooleanOperation BooleanOperation `json:"booleanOperation,omitempty"`
	// ELLIPSE
	ArcData ArcData `json:"arcData"`
	// TEXT
	Characters              string               `json:"characters,omitempty"`
	Style                   TypeStyle            `json:"style"`
	CharacterStyleOverrides []float64            `json:"characterStyleOverrides,omitempty"`
	StyleOverrideTable      map[string]TypeStyle `json:"styleOverrideTable,omitempty"`
	LineTypes               []LineTypes          `json:"lineTypes,omitempty"`
	LineIndentations        []float64            `json:"lineIndentations,omitempty"`
	// COMPONENT
	ComponentPropertyDefinitions map[string]ComponentPropertyDefinition `json:"componentPropertyDefinitions,omitempty"`
	// INSTANCE
	ComponentId         string                       `json:"componentId,omitempty"`
	IsExposedInstance   bool                         `json:"isExposedInstance,omitempty"`
	ExposedInstances    []string                     `json:"exposedInstances,omitempty"`
	ComponentProperties map[string]ComponentProperty `json:"componentProperties,omitempty"`
	Overrides           []Overrides                  `json:"overrides,omitempty"`
}

type ComponentSet struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Component struct {
	Key            string `json:"key"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ComponentSetId string `json:"componentSetId"`
}

type StyleType string

const (
	StyleTypeFill   StyleType = "FILL"
	StyleTypeText             = "TEXT"
	StyleTypeEffect           = "EFFECT"
	StyleTypeGrid             = "GRID"
)

type Style struct {
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StyleType   StyleType `json:"styleType"`
}

type Color struct {
	Red   float64 `json:"r"`
	Green float64 `json:"g"`
	Blue  float64 `json:"b"`
	Alpha float64 `json:"a"`
}

type ExportSetting struct {
	Suffix     string      `json:"suffix"`
	Format     ImageFormat `json:"format"`
	Constraint Constraint  `json:"constraint"`
}

type ImageFormat string

const (
	ImageFormatJPG  ImageFormat = "JPG"
	ImageFormatPNG              = "PNG"
	ImageFormateSVG             = "SVG"
)

type Constraint struct {
	Type  ConstraintType `json:"type"`
	Value float64        `json:"value"`
}

type ConstraintType string

const (
	ConstraintTypeScale  ConstraintType = "SCALE"
	ConstraintTypeWidth                 = "WIDTH"
	ConstraintTypeHeight                = "HEIGHT"
)

type Paint struct {
	Type                    PaintType                `json:"type"`
	Visible                 bool                     `json:"visible"`
	Opacity                 float64                  `json:"opacity"`
	Color                   Color                    `json:"color"`
	BlendMode               BlendMode                `json:"blendMode"`
	GradientHandlePositions []Vector                 `json:"gradientHandlePositions"`
	GradientStops           []ColorStop              `json:"gradientStops"`
	ScaleMode               ScaleMode                `json:"scaleMode"`
	ImageTransform          Transform                `json:"imageTransform"`
	ScalingFactor           float64                  `json:"scalingFactor"`
	Rotation                float64                  `json:"rotation"`
	ImageRef                string                   `json:"imageRef"`
	Filters                 ImageFilters             `json:"filters"`
	GifRef                  string                   `json:"gifRef"`
	BoundVariables          map[string]VariableAlias `json:"boundVariables"`
}

type PaintType string

const (
	PaintTypeSolid           PaintType = "SOLID"
	PaintTypeGradientLinear            = "GRADIENT_LINEAR"
	PaintTypeGradientRadial            = "GRADIENT_RADIAL"
	PaintTypeGradientAngular           = "GRADIENT_ANGULAR"
	PaintTypeGradientDiamond           = "GRADIENT_DIAMOND"
	PaintTypeImage                     = "IMAGE"
	PaintTypeEmoji                     = "EMOJI"
	PaintTypeVideo                     = "VIDEO"
)

type BlendMode string

const (
	BlendModeScale  BlendMode = "SCALE"
	BlendModeWidth            = "WIDTH"
	BlendModeHeight           = "HEIGHT"
	// Normal blends:
	BlendModePasshrough = "PASS_THROUGH" // (only applicable to objects with children)
	BlendModeNormal     = "NORMAL"

	// Darken:
	BlendModeDarken     = "DARKEN"
	BlendModeMultiply   = "MULTIPLY"
	BlendModeLinearBurn = "LINEAR_BURN " // ("Plus darker" in Figma)
	BlendModeColorBurn  = "COLOR_BURN"

	// Lighten:
	BlendModeLighten     = "LIGHTEN"
	BlendModeScreen      = "SCREEN"
	BlendModeLinearDodge = "LINEAR_DODGE" // ("Plus lighter" in Figma)
	BlendModeColorDodge  = "COLOR_DODGE"

	// Contrast:
	BlendModeOverlay   = "OVERLAY"
	BlendModeSoftLight = "SOFT_LIGHT"
	BlendModeHardLight = "HARD_LIGHT"

	// Inversion:
	BlendModeDifference = "DIFFERENCE"
	BlendModeExclusion  = "EXCLUSION"

	// Component:
	BlendModeHue        = "HUE"
	BlendModeSaturation = "SATURATION"
	BlendModeColor      = "COLOR"
	BlendModeLuminosity = "LUMINOSITY"
)

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ColorStop struct {
	Position       float64                  `json:"position"`
	Color          Color                    `json:"color"`
	BoundVariables map[string]VariableAlias `json:"boundVariables"`
}

type VariableAlias struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type ScaleMode string

const (
	ScaleModeFill    ScaleMode = "FILL"
	ScaleModeFit               = "FIT"
	ScaleModeTile              = "TILE"
	ScaleModeStretch           = "STRETCH"
)

type Transform [][]float64

type ImageFilters struct {
	Exposure    float64 `json:"exposure"`
	Contrast    float64 `json:"contrast"`
	Saturation  float64 `json:"saturation"`
	Temperature float64 `json:"temperature"`
	Tint        float64 `json:"tint"`
	Highlights  float64 `json:"highlights"`
	Shadows     float64 `json:"shadows"`
}

type StrokeAlign string

const (
	StrokeAlignInside  StrokeAlign = "INSIDE"
	StrokeAlignOutside             = "OUTSIDE"
	StrokeAlignCenter              = "CENTER"
)

type LayoutConstraint struct {
	Vertical   VerticalConstraint   `json:"vertical"`
	Horizontal HorizontalConstraint `json:"horizontal"`
}

type VerticalConstraint string

const (
	VerticalConstraintTop       VerticalConstraint = "TOP"
	VerticalConstraintBottom                       = "BOTTOM"
	VerticalConstraintCenter                       = "CENTER"
	VerticalConstraintTopBottom                    = "TOP_BOTTOM"
	VerticalConstraintScale                        = "SCALE"
)

type HorizontalConstraint string

const (
	HorizontalConstraintLeft      HorizontalConstraint = "LEFT"
	HorizontalConstraintRight                          = "RIGHT"
	HorizontalConstraintCenter                         = "CENTER"
	HorizontalConstraintLeftRight                      = "LEFT_RIGHT"
	HorizontalConstraintScale                          = "SCALE"
)

type LayoutAlign string

const (
	LayoutAlignInherit LayoutAlign = "INHERIT"
	LayoutAlignStretch             = "STRETCH"

	// In horizontal auto-layout frames, "MIN" and "MAX" correspond to "TOP" and "BOTTOM". In vertical auto-layout frames, "MIN" and "MAX" correspond to "LEFT" and "RIGHT".
	LayoutAlignMin    = "MIN"
	LayoutAlignCenter = "CENTER"
	LayoutAlignMax    = "MAX"
)

type Rectangle struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type LayoutMode string

const (
	LayoutModeNone       LayoutMode = "NONE"
	LayoutModeHorizontal            = "HORIZONTAL"
	LayoutModeVertical              = "VERTICAL"
)

type LayoutSizing string

const (
	LayoutSizingFixed LayoutSizing = "FIXED"
	LayoutSizingHug                = "HUG"
	LayoutSizingFill               = "FILL"
)

type LayoutWrap string

const (
	LayoutWrapNoWrap LayoutWrap = "NO_WRAP"
	LayoutWrapWrap              = "WRAP"
)

type SizingMode string

const (
	SizingModeAuto  SizingMode = "AUTO"
	SizingModeFixed            = "FIXED"
)

type AlignItems string

const (
	AlignItemsMin          AlignItems = "MIN"
	AlignItemsCenter                  = "CENTER"
	AlignItemsMax                     = "MAX"
	AlignItemsSpaceBetween            = "SPACE_BETWEEN"
	AlignItemsBaseline                = "BASELINE"
)

type AlignContent string

const (
	AlignContentAuto         AlignContent = "AUTO"
	AlignContentSpaceBetween              = "SPACE_BETWEEN"
)

type LayoutPositioning string

const (
	LayoutPositioningAuto     LayoutPositioning = "AUTO"
	LayoutPositioningAbsolute                   = "ABSOLUTE"
)

type OverflowDirection string

const (
	OverflowDirectionNone                           OverflowDirection = "NONE"
	OverflowDirectionHorizontalScrolling                              = "HORIZONTAL_SCROLLING"
	OverflowDirectionVerticalScrolling                                = "VERTICAL_SCROLLING"
	OverflowDirectionHorizontalAndVerticalScrolling                   = "HORIZONTAL_AND_VERTICAL_SCROLLING"
)

type MaskType string

const (
	MaskTypeAlpha     MaskType = "ALPHA"
	MaskTypeVector             = "VECTOR"
	MaskTypeLuminance          = "LUMINANCE"
)

type LayoutGrid struct {
	Pattern        Pattern                  `json:"pattern"`
	SectionSize    float64                  `json:"sectionSize"`
	Visible        bool                     `json:"visible"`
	Color          Color                    `json:"color"`
	Alignment      Alignment                `json:"alignment"`
	GutterSize     float64                  `json:"gutterSize"`
	Offset         float64                  `json:"offset"`
	Count          float64                  `json:"count"`
	BoundVariables map[string]VariableAlias `json:"boundVariables"`
}

type Pattern string

const (
	PatternColumns Pattern = "COLUMNS" // Vertical grid
	PatternRows            = "ROWS"    // Horizontal grid
	PatternGrid            = "GRID"    // Square grid
)

type Alignment string

const (
	AlignmentMin     Alignment = "MIN"
	AlignmentStretch           = "STRETCH"
	Alignmentcenter            = "CENTER"
)

type Effect struct {
	Type                 EffectType               `json:"type"`
	Pattern              Pattern                  `json:"pattern"`
	Visible              bool                     `json:"visible"`
	Radius               float64                  `json:"radius"`
	Color                Color                    `json:"color"`
	BlendMode            BlendMode                `json:"blendMode"`
	Offset               Vector                   `json:"offset"`
	Spread               float64                  `json:"spread"`
	ShowShadowBehindNode bool                     `json:"showShadowBehindNode"`
	BoundVariables       map[string]VariableAlias `json:"boundVariables"`
}

type EffectType string

const (
	EffectTypeInnerShadow    EffectType = "INNER_SHADOW"
	EffectTypeDropShadow                = "DROP_SHADOW"
	EffectTypeLayerBlur                 = "LAYER_BLUR"
	EffectTypeBackgroundBlur            = "BACKGROUND_BLUR"
)

type Path struct {
	Path        string  `json:"path"`
	WindingRule string  `json:"windingRule"`
	OverrideID  float64 `json:"overrideID"`
}

type PaintOverride struct {
	Fills              []Paint `json:"fills"`
	InheritFillStyleId string  `json:"inheritFillStyleId"`
}

type StrokeWeights struct {
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
}

type StrokeCap string

const (
	StrokeCapNone           StrokeCap = "NONE"
	StrokeCapRound                    = "ROUND"
	StrokeCapSquare                   = "SQUARE"
	StrokeCapLineArrow                = "LINE_ARROW"
	StrokeCapTriangleArrow            = "TRIANGLE_ARROW"
	StrokeCapDiamondFilled            = "DIAMOND_FILLED"
	StrokeCapCircleFilled             = "CIRCLE_FILLED"
	StrokeCapTriangleFilled           = "TRIANGLE_FILLED"
	StrokeCapWashiTape1               = "WASHI_TAPE_1"
	StrokeCapWashiTape2               = "WASHI_TAPE_2"
	StrokeCapWashiTape3               = "WASHI_TAPE_3"
	StrokeCapWashiTape4               = "WASHI_TAPE_4"
	StrokeCapWashiTape5               = "WASHI_TAPE_5"
	StrokeCapWashiTape6               = "WASHI_TAPE_6"
)

type StrokeJoin string

const (
	StrokeJoinMitter StrokeJoin = "MITER"
	StrokeJoinBevel             = "BEVEL"
	StrokeJoinRound             = "ROUND"
)

type BooleanOperation string

const (
	BooleanOperationUnion     BooleanOperation = "UNION"
	BooleanOperationIntersect                  = "INTERSECT"
	BooleanOperationSubtract                   = "SUBTRACT"
	BooleanOperationExclude                    = "EXCLUDE"
)

type ArcData struct {
	StartingAngle float64 `json:"startingAngle"`
	EndingAngle   float64 `json:"endingAngle"`
	InnerRadius   float64 `json:"innerRadius"`
}

type TypeStyle struct {
	FontFamily                string             `json:"fontFamily"`
	FontPostScriptName        string             `json:"fontPostScriptName"`
	FontStyle                 string             `json:"fontStyle"`
	ParagraphSpacing          float64            `json:"paragraphSpacing"`
	ParagraphIndent           float64            `json:"paragraphIndent"`
	ListSpacing               float64            `json:"listSpacing"`
	Italic                    bool               `json:"italic"`
	FontWeight                float64            `json:"fontWeight"`
	FontSize                  float64            `json:"fontSize"`
	TextCase                  TextCase           `json:"textCase"`
	TextDecoration            TextDecoration     `json:"textDecoration"`
	TextAutoResize            TextAutoResize     `json:"textAutoResize"`
	TextTruncation            TextTruncation     `json:"textTruncation"`
	MaxLines                  float64            `json:"maxLines"`
	TextAlignHorizontal       TextAlign          `json:"textAlignHorizontal"`
	TextAlignVertical         TextAlign          `json:"textAlignVertical"`
	LetterSpacing             float64            `json:"letterSpacing"`
	Hyperlink                 Hyperlink          `json:"hyperlink"`
	OpentypeFlags             map[string]float64 `json:"opentypeFlags"`
	LineHeightPx              float64            `json:"lineHeightPx"`
	LineHeightPercentFontSize float64            `json:"lineHeightPercentFontSize"`
	LineHeightUnit            LineHeightUnit     `json:"lineHeightUnit"`
	IsOverrideOverTextStyle   bool               `json:"isOverrideOverTextStyle"`
	SemanticWeight            SemanticWeight     `json:"semanticWeight"`
	SemanticItalic            SemanticItalic     `json:"semanticItalic"`
}

type TextCase string

const (
	TextCaseOriginal        TextCase = "ORIGINAL"
	TextCaseUpper                    = "UPPER"
	TextCaseLower                    = "LOWER"
	TextCaseTitle                    = "TITLE"
	TextCaseSmallCaps                = "SMALL_CAPS"
	TextCaseSmallCapsForced          = "SMALL_CAPS_FORCED"
)

type TextDecoration string

const (
	TextDecorationNone          TextDecoration = "NONE"
	TextDecorationStrikethrough                = "STRIKETHROUGH"
	TextDecorationUnderline                    = "UNDERLINE"
)

type TextAutoResize string

const (
	TextAutoResizeNone           TextAutoResize = "NONE"
	TextAutoResizeHeight                        = "HEIGHT"
	TextAutoResizeWidthAndHeight                = "WIDTH_AND_HEIGHT"
)

type TextTruncation string

const (
	TextTruncationDisabled TextTruncation = "DISABLED"
	TextTruncationEnding                  = "ENDING"
)

type TextAlign string

const (
	TextAlignLeft      TextAlign = "LEFT"
	TextAlignRight               = "RIGHT"
	TextAlignCenter              = "CENTER"
	TextAlignJustified           = "JUSTIFIED"
	TextAlignTop                 = "TOP"
	TextAlignBottom              = "BOTTOM"
)

type Hyperlink struct {
	Type   HyperlinkType `json:"type"`
	Url    string        `json:"url"`
	NodeID string        `json:"nodeID"`
}

type HyperlinkType string

const (
	HyperlinkTypeURL  HyperlinkType = "URL"
	HyperlinkTypeNode               = "NODE"
)

type LineHeightUnit string

const (
	LineHeightUnitPixels    LineHeightUnit = "PIXELS"
	LineHeightUnitFontSize                 = "FONT_SIZE_%"
	LineHeightUnitIntrinsic                = "INTRINSIC_%"
)

type SemanticWeight string

const (
	SemanticWeightBold   SemanticWeight = "BOLD"
	SemanticWeightNormal                = "NORMAL"
)

type SemanticItalic string

const (
	SemanticItalicItalic SemanticItalic = "ITALIC"
	SemanticItalicNormal                = "NORMAL"
)

type LineTypes string

const (
	LineTypesOrdered   LineTypes = "ORDERED"
	LineTypesUnordered           = "UNORDERED"
	LineTypesNone                = "NONE"
)

type ComponentPropertyDefinition struct {
	Type            ComponentPropertyType `json:"type"`
	DefaultValue    string                `json:"defaultValue"`
	VariantOptions  []string              `json:"variantOptions,omitempty"`
	PreferredValues StyleType             `json:"preferredValues,omitempty"`
}

type ComponentPropertyType string

const (
	ComponentPropertyTypeBoolean      ComponentPropertyType = "BOOLEAN"
	ComponentPropertyTypeInstanceSwap                       = "INSTANCE_SWAP"
	ComponentPropertyTypeText                               = "TEXT"
	ComponentPropertyTypeVariant                            = "VARIANT"
)

type ComponentProperty struct {
	Type            ComponentPropertyType    `json:"type"`
	Value           string                   `json:"value"`
	PreferredValues StyleType                `json:"preferredValues,omitempty"`
	BoundVariables  map[string]VariableAlias `json:"boundVariables"`
}

type Overrides struct {
	ID               string   `json:"id"`
	OverriddenFields []string `json:"overriddenFields"`
}

type Token struct {
	Name      string
	Variable  string
	Value     string
	Theme     string
	ClassName string
}

// Figma Variables types
type Variables struct {
	Status float64 `json:"status"`
	Error  bool    `json:"error"`
	Meta   Meta    `json:"meta"`
}

type Meta struct {
	Variables           map[string]Variable           `json:"variables"`
	VariableCollections map[string]VariableCollection `json:"variableCollections"`
}

type Variable struct {
	ID                   string                 `json:"id"`
	Name                 string                 `json:"name"`
	VariableCollectionId string                 `json:"variableCollectionId"`
	ResolvedType         ResolvedType           `json:"resolvedType"`
	ValuesByMode         map[string]interface{} `json:"valuesByMode"`
	Remote               bool                   `json:"remote"`
	Description          string                 `json:"description"`
	HiddenFromPublishing bool                   `json:"hiddenFromPublishing"`
	Scopes               []VariableScope        `json:"scopes"`
	CodeSyntax           VariableCodeSyntax     `json:"codeSyntax"`
	DeletedButReferenced bool                   `json:"deletedButReferenced"`
}

type VariableCollection struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Key                  string   `json:"key"`
	Modes                []Modes  `json:"modes"`
	DefaultModeId        string   `json:"defaultModeId"`
	Remote               bool     `json:"remote"`
	HiddenFromPublishing bool     `json:"hiddenFromPublishing"`
	VariableIds          []string `json:"variableIds"`
}

type Modes struct {
	ModeId string `json:"modeId"`
	Name   string `json:"name"`
}

type ResolvedType string

const (
	ResolvedTypeBoolean ResolvedType = "BOOLEAN"
	ResolvedTypeFloat                = "FLOAT"
	ResolvedTypeString               = "STRING"
	ResolvedTypeColor                = "COLOR"
)

type VariableScope string

const (
	// Valid scopes for FLOAT variables:
	VariableScopeAllScopes        VariableScope = "ALL_SCOPES"
	VariableScopeCornerRadius                   = "CORNER_RADIUS"
	VariableScopeTextContent                    = "TEXT_CONTENT"
	VariableScopeWidthHeight                    = "WIDTH_HEIGHT"
	VariableScopeGap                            = "GAP"
	VariableScopeStrokeFloat                    = "STROKE_FLOAT"
	VariableScopeOpacity                        = "OPACITY"
	VariableScopeEffectFloat                    = "EFFECT_FLOAT"
	VariableScopeFontWeight                     = "FONT_WEIGHT"
	VariableScopeFontSize                       = "FONT_SIZE"
	VariableScopeLineHeight                     = "LINE_HEIGHT"
	VariableScopeLetterSpacing                  = "LETTER_SPACING"
	VariableScopeParagraphSpacing               = "PARAGRAPH_SPACING"
	VariableScopeParagraphIndent                = "PARAGRAPH_INDENT"
	// Valid scopes for STRING variables:
	// VariableScopeAllScopes 	= "ALL_SCOPES"
	// VariableScopeTextContent = "TEXT_CONTENT"
	VariableScopeFontFamily     = "FONT_FAMILY"
	VariableScopeFontStyle      = "FONT_STYLE"
	VariableScopeFontVariations = "FONT_VARIATIONS"
	// Valid scopes for COLOR variables:
	// VariableScopeAllScopes = "ALL_SCOPES"
	VariableScopeAllFills    = "ALL_FILLS"
	VariableScopeFrameFill   = "FRAME_FILL"
	VariableScopeShapeFill   = "SHAPE_FILL"
	VariableScopeTextFill    = "TEXT_FILL"
	VariableScopeStrokeColor = "STROKE_COLOR"
	VariableScopeEffectColor = "EFFECT_COLOR"
)

type VariableCodeSyntax struct {
	Web     string `json:"WEB"`
	Android string `json:"ANDROID"`
	Ios     string `json:"iOS"`
}
