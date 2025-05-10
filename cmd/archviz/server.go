package main

import (
	"log/slog"
	"time"

	"github.com/a-h/templ"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"oss.terrastruct.com/d2/lib/log"

	"github.com/edanko/archviz/cmd/archviz/internal/ports"
	"github.com/edanko/archviz/cmd/archviz/templates"
	"github.com/edanko/archviz/pkg/slogfiber"
)

// runServer runs a new HTTP server with the loaded environment variables.
func runServer(handlers *ports.HTTPServer, logger *slog.Logger) *fiber.App {
	config := fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		JSONEncoder:  sonic.ConfigFastest.Marshal,
		JSONDecoder:  sonic.ConfigFastest.Unmarshal,
	}

	server := fiber.New(config)
	server.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	server.Use(slogfiber.New(logger))
	server.Use(func(c fiber.Ctx) error {
		ctx := c.Context()
		ctx = log.With(ctx, logger)

		c.SetContext(ctx)

		return c.Next()
	})
	server.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	setupStaticHandler(server)

	server.Get("/", handlers.IndexPage)

	server.Get("/api/hello-world", showContentAPIHandler)

	server.Get("/views/:id", handlers.RenderView)
	server.Get("/components/:id", handlers.RenderSingleComponentView)

	server.Use(notFoundMiddleware)

	return server
}

func notFoundMiddleware(c fiber.Ctx) error {
	c.Status(fiber.StatusNotFound)
	return Render(c, templates.NotFound())
}

func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

// showContentAPIHandler handles an API endpoint to show content.
func showContentAPIHandler(c fiber.Ctx) error {
	// Check, if the current request has a 'HX-Request' header.
	// For more information, see https://htmx.org/docs/#request-headers
	if c.Get("HX-Request") == "" || c.Get("HX-Request") != "true" {
		// If not, return HTTP 400 error.
		return fiber.NewError(fiber.StatusBadRequest, "non-htmx request")
	}

	return c.SendString("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>")
}
