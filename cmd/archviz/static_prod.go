//go:build !dev
// +build !dev

package main

import (
	"embed"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

//go:embed static
var publicFS embed.FS

func setupStaticHandler(app *fiber.App) {
	app.Get("/static/*", static.New("/static", static.Config{
		FS:     publicFS,
		Browse: false,
	}))
}
