package ports

import (
	"context"
	"io"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"oss.terrastruct.com/d2/lib/log"

	"github.com/edanko/archviz/cmd/archviz/internal/app"
	"github.com/edanko/archviz/cmd/archviz/templates/components"
	"github.com/edanko/archviz/cmd/archviz/templates/pages"
)

type HTTPServer struct {
	app *app.Application
}

func NewHTTPServer(app *app.Application) *HTTPServer {
	return &HTTPServer{
		app: app,
	}
}

func Unsafe(rawContent string) templ.Component {
	return templ.ComponentFunc(
		func(_ context.Context, w io.Writer) (err error) {
			_, err = io.WriteString(w, rawContent)
			return
		},
	)
}

func (s *HTTPServer) IndexPage(c fiber.Ctx) error {
	out := []byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>")

	menuItems, err := s.app.Queries.ConstructMenu.Handle(c.Context())
	if err != nil {
		return err
	}

	var templateHandler *templ.ComponentHandler

	if requestFullPage(c) {
		// Render full page layout
		log.Info(c.Context(), "Rendering full page layout")
		templateHandler = templ.Handler(
			pages.Index(
				"Welcome to example!",
				"You're here because it worked out.",
				"gowebly, htmx example page, go with htmx",
				Unsafe(string(out)),
				menuItems,
			),
		)
	} else {
		// Render partial content (without full page layout)
		log.Info(c.Context(), "Rendering partial content")
		templateHandler = templ.Handler(
			// pages.PartialContent( // Assume this component exists in your templates
			Unsafe(string(out)),
			// h,
			// ),
		)
	}

	return adaptor.HTTPHandler(templateHandler)(c)
}

func (s *HTTPServer) RenderView(c fiber.Ctx) error {
	diagramBytes, err := s.app.Queries.RenderView.Handle(
		c.Context(),
		c.Params("id"),
	)
	if err != nil {
		return err
	}

	component := components.ShowSVG(diagramBytes)

	if requestFullPage(c) {
		return s.wrapInBaseLayout(c, component)
	}

	templateHandler := templ.Handler(
		component,
	)
	return adaptor.HTTPHandler(templateHandler)(c)
}

func (s *HTTPServer) RenderSingleComponentView(c fiber.Ctx) error {
	diagramBytes, err := s.app.Queries.RenderSingleComponentView.Handle(
		c.Context(),
		c.Params("id"),
	)
	if err != nil {
		return err
	}

	component := components.ShowSVG(diagramBytes)

	if requestFullPage(c) {
		return s.wrapInBaseLayout(c, component)
	}

	templateHandler := templ.Handler(
		component,
	)
	return adaptor.HTTPHandler(templateHandler)(c)
}

const (
	// FIXME: hxBoosted should be here too.
	hxRequestHeader               = "HX-Request"
	hxHistoryRestoreRequestHeader = "HX-History-Restore-Request"
)

func requestFullPage(c fiber.Ctx) bool {
	isHXRequest, _ := strconv.ParseBool(c.Get(hxRequestHeader, "false"))
	if !isHXRequest {
		return true
	}
	isHXHistoryRestoreRequest, _ := strconv.ParseBool(c.Get(hxHistoryRestoreRequestHeader, "false"))
	return isHXHistoryRestoreRequest
}

func (s *HTTPServer) wrapInBaseLayout(c fiber.Ctx, component templ.Component) error {
	menuItems, err := s.app.Queries.ConstructMenu.Handle(c.Context())
	if err != nil {
		return err
	}

	return adaptor.HTTPHandler(templ.Handler(
		pages.Index(
			"Welcome to example!",
			"You're here because it worked out.",
			"gowebly, htmx example page, go with htmx",
			component,
			menuItems,
		),
	))(c)
}
