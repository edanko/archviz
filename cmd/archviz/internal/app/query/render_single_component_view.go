package query

import (
	"context"
	"math"

	"github.com/edanko/archviz/cmd/archviz/internal/domain"
)

// RenderSingleComponentViewHandler renders a single component view.
type RenderSingleComponentViewHandler struct {
	s  sourceRepository
	df diagramBuilderFactory
}

// NewRenderSingleComponentViewHandler creates a new RenderSingleComponentViewHandler.
func NewRenderSingleComponentViewHandler(
	s sourceRepository,
	df diagramBuilderFactory,
) RenderSingleComponentViewHandler {
	return RenderSingleComponentViewHandler{
		s:  s,
		df: df,
	}
}

// Handle renders a single component view.
func (h RenderSingleComponentViewHandler) Handle(
	ctx context.Context,
	componentID string,
) ([]byte, error) {
	g, err := h.s.GetGraph(ctx)
	if err != nil {
		return nil, err
	}

	filter := &domain.Filter{
		IDs:        []string{componentID},
		DepthRange: [2]int{0, math.MaxInt},
		// ConnectedTo: viewID,
		// ConnectedFrom: viewID,
		IsolateSubgraph: true,
		// WithChildren:    true,
		WithDependents: true,
	}

	cs, ls := g.Filter(filter)

	d, err := h.df.Create(ctx, cs, ls)
	if err != nil {
		return nil, err
	}

	return d, nil
}
