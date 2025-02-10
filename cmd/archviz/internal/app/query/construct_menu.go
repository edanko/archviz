package query

import (
	"context"

	"github.com/edanko/archviz/cmd/archviz/internal/domain"
)

// ConstructMenuHandler is a handler for constructing a menu.
type ConstructMenuHandler struct {
	s sourceRepository
}

// NewConstructMenuHandler creates a new ConstructMenuHandler.
func NewConstructMenuHandler(
	s sourceRepository,
) ConstructMenuHandler {
	return ConstructMenuHandler{
		s: s,
	}
}

// Handle constructs a menu.
func (h ConstructMenuHandler) Handle(
	ctx context.Context,
) (map[string][]domain.View, error) {
	views, err := h.s.ListViews(ctx)
	if err != nil {
		return nil, err
	}

	data := views.BuildHierarchy()

	return data, nil
}
