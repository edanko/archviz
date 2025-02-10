package service

import (
	"context"

	"github.com/edanko/archviz/cmd/archviz/internal/adapters"
	"github.com/edanko/archviz/cmd/archviz/internal/app"
	"github.com/edanko/archviz/cmd/archviz/internal/app/query"
)

var (
	// TODO: move to config
	paths = []string{
		"/home/egor/work/DocHub/public/workspace/ac/components",
		"/home/egor/work/DocHub/public/workspace/ac/contexts",
	}
)

func NewApplication(
	ctx context.Context,
	cfg *adapters.Config,
) (application *app.Application, cleanup func()) {
	dataStore, err := adapters.NewLocalDirectorySource(
		ctx,
		paths,
	)
	if err != nil {
		panic(err)
	}

	df := app.NewDiagramBuilderFactory(&cfg.Render)

	return &app.Application{
			Commands: app.Commands{},
			Queries: app.Queries{
				RenderView:                query.NewRenderViewHandler(dataStore, df),
				ConstructMenu:             query.NewConstructMenuHandler(dataStore),
				RenderSingleComponentView: query.NewRenderSingleComponentViewHandler(dataStore, df),
			},
		},
		func() {
			_ = dataStore.Close()
		}
}
