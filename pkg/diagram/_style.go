package diagram

import (
	"strconv"

	"github.com/go-faster/errors"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2oracle"
)

// styleOption is a function that applies a style option to a node.
type styleOption func() (key string, value string, err error)

// WithStyle sets the style of the whole diagram, node or connection.
func WithStyle(opts ...styleOption) GenericOption {
	return func(graph *d2graph.Graph, _ callerKind, key string) (*d2graph.Graph, string, error) {
		for _, opt := range opts {
			k, v, err := opt()
			if err != nil {
				return graph, key, errors.Wrap(err, "failed to apply style option")
			}
			graph, err = d2oracle.Set(
				graph,
				nil,
				key+".style."+k,
				nil,
				&v,
			)
			if err != nil {
				return graph, key, errors.Wrap(err, "failed to set style")
			}
		}

		return graph, key, nil
	}
}

// WithOpacity sets the opacity of the node.
// The value must be between 0 and 1.
func WithOpacity(opacity float64) styleOption {
	return func() (key string, value string, err error) {
		if opacity < 0 || opacity > 1 {
			return "", "", errors.New("opacity must be between 0 and 1")
		}

		return "opacity", strconv.FormatFloat(opacity, 'f', -1, 64), nil
	}
}

// WithStroke sets the stroke of the node.
// The value must be a CSS color name, hex code, or a subset of CSS gradient strings.
// For sql_tables and classes, stroke is applied as fill to the body
// (since fill is already used to control header's fill).
func WithStroke(stroke string) styleOption {
	return func() (key string, value string, err error) {
		return "stroke", stroke, nil
	}
}

// WithFill sets the fill of the node.
// Must be applied to shapes only.
// CSS color name, hex code, or a subset of CSS gradient strings.
// For sql_tables and classes, fill is applied to the header.
func WithFill(fill string) styleOption {
	return func() (key string, value string, err error) {
		return "fill", fill, nil
	}
}

// WithFillPattern sets the fill pattern of the node.
// Must be applied to shapes only.
func WithFillPattern(fillPattern string) styleOption {
	return func() (key string, value string, err error) {
		switch fillPattern {
		case "dots", "lines", "grain", "none":
			// ok
		default:
			return "", "", errors.New("fill pattern must be one of dots, lines, grain, or none")
		}
		return "fill-pattern", fillPattern, nil
	}
}

// WithStrokeWidth sets the stroke width of the node.
// The value must be between 1 and 15.
func WithStrokeWidth(strokeWidth int) styleOption {
	return func() (key string, value string, err error) {
		if strokeWidth < 1 || strokeWidth > 15 {
			return "", "", errors.New("strocke width must be between 1 and 15")
		}

		return "stroke-width", strconv.Itoa(strokeWidth), nil
	}
}

// WithStrokeDash sets the stroke dash of the node.
// Must be applied to shapes only.
// The value must be between 1 and 10.
func WithStrokeDash(strokeDash int) styleOption {
	return func() (key string, value string, err error) {
		if strokeDash < 1 || strokeDash > 10 {
			return "", "", errors.New("stroke dash must be between 1 and 10")
		}

		return "stroke-dash", strconv.Itoa(strokeDash), nil
	}
}

// WithBorderRadius sets the border radius of the node.
// The value must be between 1 and 20.
// border-radius works on connections too, which controls how rounded the corners are.
// This only applies to layout engines that use corners (e.g. ELK), and of course,
// only has effect on connections whose routes have corners.
// Specifying a very high value creates a "pill" effect.
func WithBorderRadius(borderRadius int) styleOption {
	return func() (key string, value string, err error) {
		if borderRadius < 1 || borderRadius > 20 {
			return "", "", errors.New("border radius must be between 1 and 20")
		}

		return "border-radius", strconv.Itoa(borderRadius), nil
	}
}

// WithShadow sets the shadow of the node.
// Must be applied to shapes only.
func WithShadow(shadow bool) styleOption {
	return func() (key string, value string, err error) {
		return "shadow", strconv.FormatBool(shadow), nil
	}
}

// WithIs3d sets the is3d of the node.
// Must be applied to rectangle, square or hexagon shapes only.
func With3D(is3D bool) styleOption {
	return func() (key string, value string, err error) {
		return "3d", strconv.FormatBool(is3D), nil
	}
}

// WithMultiple sets the multiple of the node.
// Must be applied to shapes only.
func WithMultiple(multiple bool) styleOption {
	return func() (key string, value string, err error) {
		return "multiple", strconv.FormatBool(multiple), nil
	}
}

// WithDoubleBorder sets the double border of the node.
// Must be applied to rectangles and ovals shapes only.
func WithDoubleBorder(doubleBorder bool) styleOption {
	return func() (key string, value string, err error) {
		return "double-border", strconv.FormatBool(doubleBorder), nil
	}
}

// WithFont sets the font of the node.
// The only option for now is to specify mono.
func WithFont(font string) styleOption {
	return func() (key string, value string, err error) {
		if font != "mono" {
			return "", "", errors.New("currently the only option is to specify mono")
		}
		return "font", font, nil
	}
}

// WithFontSize sets the font size of the node.
// The value must be between 8 and 100.
func WithFontSize(size int) styleOption {
	return func() (key string, value string, err error) {
		if size < 8 || size > 100 {
			return "", "", errors.New("font size must be between 8 and 100")
		}
		return "font-size", strconv.Itoa(size), nil
	}
}

// WithFontColor sets the font color of the node.
// The value must be a CSS color name, hex code, or a subset of CSS gradient strings.
// For sql_tables and classes, font-color is applied to the header text only (theme controls other colors in the body).
func WithFontColor(color string) styleOption {
	return func() (key string, value string, err error) {
		return "font-color", color, nil
	}
}

// WithAnimated sets the animated of the node.
// Must be applied to connections only.
func WithAnimated(animated bool) styleOption {
	return func() (key string, value string, err error) {
		return "animated", strconv.FormatBool(animated), nil
	}
}

// WithBold sets the bold of the node.
func WithBold(bold bool) styleOption {
	return func() (key string, value string, err error) {
		return "bold", strconv.FormatBool(bold), nil
	}
}

// WithItalic sets the italic of the node.
func WithItalic(italic bool) styleOption {
	return func() (key string, value string, err error) {
		return "italic", strconv.FormatBool(italic), nil
	}
}

// WithUnderline sets the underline of the node.
func WithUnderline(underline bool) styleOption {
	return func() (key string, value string, err error) {
		return "underline", strconv.FormatBool(underline), nil
	}
}

// WithTextTransform changes the casing of labels.
func WithTextTransform(textTransform string) styleOption {
	return func() (key string, value string, err error) {
		switch textTransform {
		case "uppercase", "lowercase", "title", "none":
			// ok
		default:
			return "", "", errors.New("text transform must be one of uppercase, lowercase, title, or none")
		}
		return "text-transform", textTransform, nil
	}
}
