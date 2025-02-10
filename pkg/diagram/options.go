package diagram

// WithGridRows sets the number of rows in the grid to entire diagram.
// func WithGridRows(rows int) DiagramOption {
// 	return func(d *Diagram) error {
// 		var err error
// 		d.graph, err = d2oracle.Set(
// 			d.graph,
// 			nil,
// 			"grid-rows",
// 			nil,
// 			go2.Pointer(fmt.Sprintf("%d", rows)),
// 		)
// 		return err
// 	}
// }

// // WithGridColumns sets the number of columns in the grid to entire diagram.
// func WithGridColumns(columns int) DiagramOption {
// 	return func(d *Diagram) error {
// 		var err error
// 		d.graph, err = d2oracle.Set(
// 			d.graph,
// 			nil,
// 			"grid-columns",
// 			nil,
// 			go2.Pointer(fmt.Sprintf("%d", columns)),
// 		)
// 		return err
// 	}
// }

// // WithGridGap sets the gap between the nodes in the grid.
// func WithGridGap(gap int) DiagramOption {
// 	return func(d *Diagram) error {
// 		var err error
// 		d.graph, err = d2oracle.Set(
// 			d.graph,
// 			nil,
// 			"grid-gap",
// 			nil,
// 			go2.Pointer(fmt.Sprintf("%d", gap)),
// 		)
// 		return err
// 	}
// }

// // WithDirection sets the direction of the diagram.
// func WithDirection(direction DiagramDirection) DiagramOption {
// 	return func(d *Diagram) error {
// 		var err error
// 		d.graph, err = d2oracle.Set(
// 			d.graph,
// 			nil,
// 			"direction",
// 			nil,
// 			go2.Pointer(string(direction)),
// 		)
// 		return err
// 	}
// }

// WithRootStyle sets the style of the root node.
// Only some styles are applicable at the root level:
// // fill, fill-pattern, stroke, stroke-width, stroke-dash and double-border.
// func WithRootStyle(opts ...styleOption) DiagramOption {
// 	return func(d *Diagram) error {
// 		validOptions := map[string]bool{
// 			"WithFill":         true,
// 			"WithFillPattern":  true,
// 			"WithStroke":       true,
// 			"WithStrokeWidth":  true,
// 			"WithStrokeDash":   true,
// 			"WithDoubleBorder": true,
// 		}

// 		for _, opt := range opts {
// 			fnName := getFunctionName(opt)
// 			if !validOptions[fnName] {
// 				return errors.Errorf("style option %s is not applicable at the root level", fnName)
// 			}

// 			k, v, err := opt()
// 			if err != nil {
// 				return errors.Wrap(err, "failed to apply style option")
// 			}

// 			graph, err := d2oracle.Set(
// 				d.graph,
// 				nil,
// 				"style."+k,
// 				nil,
// 				&v,
// 			)
// 			if err != nil {
// 				return errors.Wrap(err, "failed to set style")
// 			}

// 			d.graph = graph
// 		}
// 		return nil
// 	}
// }

// // getFunctionName returns the name of the function provided as an argument.
// func getFunctionName(i any) string {
// 	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
// 	parts := strings.Split(name, ".")
// 	return parts[len(parts)-2]
// }

// func WithClasses(classes []map[string]string) DiagramOption {
// 	return func(d *Diagram) error {

// 		var err error
// 		d.graph, _, err = d2oracle.Create(
// 			d.graph,
// 			nil,
// 			"classes",
// 		)
// 		if err != nil {
// 			return errors.Wrap(err, "failed to create classes")
// 		}

// 		for _, attributes := range classes {
// 			className := "classes." + attributes["name"]

// 			d.graph, _, err = d2oracle.Create(
// 				d.graph,
// 				nil,
// 				className,
// 			)
// 			if err != nil {
// 				return errors.Wrap(err, "failed to create class")
// 			}

// 			for key, value := range attributes {
// 				if key == "name" {
// 					continue
// 				}

// 				var err error
// 				d.graph, err = d2oracle.Set(
// 					d.graph,
// 					nil,
// 					className+"."+key,
// 					nil,
// 					&value,
// 				)
// 				if err != nil {
// 					return errors.Wrap(err, "failed to set class attribute")
// 				}
// 			}
// 		}
// 		return nil
// 	}
// }
