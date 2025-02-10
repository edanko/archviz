package diagram

// Reference: https://d2lang.com/tour/shapes
type NodeShape string

// String returns the string representation of the NodeShape.
func (s NodeShape) String() string {
	return string(s)
}

const (
	NodeShapeRectangle     NodeShape = "rectangle"
	NodeShapeSquare        NodeShape = "square"
	NodeShapePage          NodeShape = "page"
	NodeShapeParallelogram NodeShape = "parallelogram"
	NodeShapeDocument      NodeShape = "document"
	NodeShapeCylinder      NodeShape = "cylinder"
	NodeShapeQueue         NodeShape = "queue"
	NodeShapePackage       NodeShape = "package"
	NodeShapeStep          NodeShape = "step"
	NodeShapeCallout       NodeShape = "callout"
	NodeShapeStoredData    NodeShape = "stored_data"
	NodeShapePerson        NodeShape = "person"
	NodeShapeDiamond       NodeShape = "diamond"
	NodeShapeOval          NodeShape = "oval"
	NodeShapeCircle        NodeShape = "circle"
	NodeShapeHexagon       NodeShape = "hexagon"
	NodeShapeCloud         NodeShape = "cloud"
	NodeShapeImage         NodeShape = "image"
	NodeShapeText          NodeShape = "text"
)

// Ref: https://d2lang.com/tour/connections#arrowheads
type ArrowHead string

const (
	// ArrowHeadTriangle is the triangle arrow head.
	// Can be further styled as style.filled: false.
	ArrowHeadTriangle ArrowHead = "triangle"
	// ArrowHeadArrow like triangle but pointier.
	ArrowHeadArrow ArrowHead = "arrow"
	// ArrowHeadDiamond
	// Can be further styled as style.filled: false.
	ArrowHeadDiamond ArrowHead = "diamond"
	// ArrowHeadCircle
	// Can be further styled as style.filled: false.
	ArrowHeadCircle            ArrowHead = "circle"
	ArrowCrowsFootOne          ArrowHead = "cf-one"
	ArrowCrowsFootOneRequired  ArrowHead = "cf-one-required"
	ArrowCrowsFootMany         ArrowHead = "cf-many"
	ArrowCrowsFootManyRequired ArrowHead = "cf-many-required"
)

type DiagramDirection string

const (
	Left  DiagramDirection = "left"
	Right DiagramDirection = "right"
	Up    DiagramDirection = "up"
	Down  DiagramDirection = "down"
)

// Reference: https://d2lang.com/tour/positions
type Position string

const (
	TopLeft      Position = "top-left"
	TopCenter    Position = "top-center"
	TopRight     Position = "top-right"
	CenterLeft   Position = "center-left"
	CenterRight  Position = "center-right"
	BottomLeft   Position = "bottom-left"
	BottomCenter Position = "bottom-center"
	BottomRight  Position = "bottom-right"
)
