package query

import (
	"context"

	"github.com/edanko/archviz/cmd/archviz/internal/domain"
)

// sourceRepository is a source of data for rendering.
type sourceRepository interface {
	GetGraph(ctx context.Context) (*domain.ComponentGraph, error)
	GetView(ctx context.Context, viewID string) (domain.View, error)
	ListViews(ctx context.Context) (domain.Views, error)
}

// diagramBuilderFactory is a factory for creating diagrams.
type diagramBuilderFactory interface {
	Create(
		ctx context.Context,
		cs []*domain.Component,
		ls []domain.Link,
	) ([]byte, error)
}
