package app

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"strings"

	"oss.terrastruct.com/d2/lib/log"

	"github.com/edanko/archviz/cmd/archviz/internal/adapters"
	"github.com/edanko/archviz/cmd/archviz/internal/domain"
	"github.com/edanko/archviz/pkg/diagram"
)

// DiagramBuilderFactory is a factory for creating diagrams.
type DiagramBuilderFactory struct {
	baseConfig *adapters.RenderConfig
}

// NewDiagramBuilderFactory creates a new DiagramBuilderFactory.
func NewDiagramBuilderFactory(
	baseConfig *adapters.RenderConfig,
) *DiagramBuilderFactory {
	return &DiagramBuilderFactory{
		baseConfig: baseConfig,
	}
}

// buildOptionsFromConfig builds options for the diagram from the config.
// TODO: there should be overrides, for example - theme choice based on the frontend theme toggle.
// Or overrides can be set in the yamls with components.
func (b *DiagramBuilderFactory) buildOptionsFromConfig() []diagram.RenderOption {
	var opts []diagram.RenderOption

	if b.baseConfig.Minify {
		opts = append(opts, diagram.WithMinifier())
	}

	if b.baseConfig.ImageBundler {
		opts = append(opts, diagram.WithImageBundler(true))
	}

	switch b.baseConfig.LayoutEngine {
	case "dagre":
		// it is the default, so nothing to do here.
	case "tala":
		opts = append(opts, diagram.WithTALALayoutEngine())
	case "elk":
		opts = append(opts, diagram.WithELKLayoutEngine())
	}

	// if b.baseConfig.Direction != "" {
	// 	opts = append(
	// 		opts,
	// 		diagram.WithDirection(
	// 			diagram.DiagramDirection(b.baseConfig.Direction),
	// 		),
	// 	)
	// }

	// if b.baseConfig.Classes != nil {
	// 	opts = append(
	// 		opts,
	// 		diagram.WithClasses(b.baseConfig.Classes),
	// 	)
	// }

	return opts
}

// some interesting code here: https://github.com/laupse/cloud-view/blob/main/internal/app/aws.go
func (b *DiagramBuilderFactory) Create(
	ctx context.Context,
	components []*domain.Component,
	links []domain.Link,
) ([]byte, error) {
	renderer, err := diagram.NewRender(
		b.buildOptionsFromConfig()...,
	)
	if err != nil {
		return nil, err
	}

	r, err := diagram.NewDiagram()
	if err != nil {
		return nil, err
	}

	for _, attributes := range b.baseConfig.Classes {
		copywoname := make(map[string]string)
		maps.Copy(copywoname, attributes)
		delete(copywoname, "name")

		r.AddClass(attributes["name"], copywoname)
	}

	for _, c := range components {
		nodeOpts := []diagram.NodeOption{
			diagram.WithLabel(""),
			diagram.WithNodeAttribute("icon.near", "top-center"),
		}
		if c.Properties["class"].(string) != "" {
			nodeOpts = append(nodeOpts, diagram.WithClass(c.Properties["class"].(string)))
		}

		r.AddNode(
			c.ID,
			nodeOpts...,
		)

		markdownContent := constructComponentMarkdownExplaination(c)
		r.AddNode(c.ID+".explanation",
			diagram.WithCode("md", markdownContent),
		)
	}

	for _, link := range links {
		r.AddDirectedConnection(
			link.SourceID,
			link.TargetID,
			diagram.WithConnectionLabel(link.Description),
		)
		log.Info(ctx, "added single link", slog.String("source", link.SourceID), slog.String("target", link.TargetID))
	}

	return renderer.Render(ctx, r.D2())
}

// wordWrapLineLength is the default line length for word wrapping.
const wordWrapLineLength = 65

// constructComponentMarkdownExplaination constructs a
// markdown explanation for the given component.
func constructComponentMarkdownExplaination(
	component *domain.Component,
) string {
	description := strings.TrimSpace(component.Description)
	description = WordWrap(description, wordWrapLineLength)
	return fmt.Sprintf(
		"# [%s](/components/%s)\n\n%s",
		component.Title, component.ID, description,
	)
}

// WordWrap wraps lines in the input text that exceed the specified line length.
// Markdown style, with backslashes at the end of wrapped lines.
func WordWrap(text string, lineLength int) string {
	var result strings.Builder

	for line := range strings.Lines(text) {
		if line == "" {
			result.WriteString("\n")
			continue
		}

		if len(line) <= lineLength {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			result.WriteString("\n")
			continue
		}

		currentLineLength := 0
		for i, word := range words {
			wordLength := len(word)

			if currentLineLength > 0 {
				if currentLineLength+1+wordLength > lineLength {
					result.WriteString(" \\\n")
					currentLineLength = 0
				}
			}

			if currentLineLength > 0 {
				result.WriteString(" ")
				currentLineLength++
			}

			result.WriteString(word)
			currentLineLength += wordLength

			if i == len(words)-1 {
				result.WriteString("\n")
			}
		}
	}

	return strings.TrimSuffix(result.String(), "\n") + "\n"
}
