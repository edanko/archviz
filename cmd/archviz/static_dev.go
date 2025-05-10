//go:build dev
// +build dev

package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func setupStaticHandler(app *fiber.App) {
	app.Get("/static/*", static.New("./static"))
}
