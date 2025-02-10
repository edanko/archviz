package query

import (
	"context"

	"github.com/edanko/archviz/cmd/archviz/internal/domain"
)

// RenderViewHandler renders the diagram for the given view.
type RenderViewHandler struct {
	s  sourceRepository
	df diagramBuilderFactory
}

// NewRenderViewHandler creates a new RenderViewHandler.
func NewRenderViewHandler(
	s sourceRepository,
	df diagramBuilderFactory,
) RenderViewHandler {
	return RenderViewHandler{
		s:  s,
		df: df,
	}
}

// Handle renders the diagram for the given view.
func (h RenderViewHandler) Handle(
	ctx context.Context,
	viewID string,
) ([]byte, error) {
	g, err := h.s.GetGraph(ctx)
	if err != nil {
		return nil, err
	}

	view, err := h.s.GetView(ctx, viewID)
	if err != nil {
		return nil, err
	}

	filter := &domain.Filter{
		DepthRange:      [2]int{0, view.MaxDepth},
		IDs:             view.Components,
		IsolateSubgraph: true,
		Aggregate:       true,
		// WithDependents:  true,
		// WithChildren:    true,
	}

	cs, ls := g.Filter(filter)

	d, err := h.df.Create(ctx, cs, ls)
	if err != nil {
		return nil, err
	}

	return d, nil
}
