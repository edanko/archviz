package app

import "github.com/edanko/archviz/cmd/archviz/internal/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct{}

type Queries struct {
	RenderView                query.RenderViewHandler
	ConstructMenu             query.ConstructMenuHandler
	RenderSingleComponentView query.RenderSingleComponentViewHandler
}
