package http

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"oss.terrastruct.com/d2/lib/log"
)

type FiberAdapter struct {
	address  string
	fiberApp *fiber.App
}

// NewFiberAdapter initializes and returns a new instance of the FiberAdapter.
func NewFiberAdapter(
	fiberApp *fiber.App,
	addr string,
) *FiberAdapter {
	return &FiberAdapter{
		address:  addr,
		fiberApp: fiberApp,
	}
}

// Start starts the HTTP listener.
func (a *FiberAdapter) Start(ctx context.Context) error {
	log.Info(ctx, "starting HTTP listener", slog.String("address", a.address))

	return a.fiberApp.Listen(
		a.address,
		fiber.ListenConfig{
			DisableStartupMessage: true,
		},
	)
}

// Stop gracefully stops the FiberAdapter.
func (a *FiberAdapter) Stop(ctx context.Context) error {
	return a.fiberApp.ShutdownWithContext(ctx)
}
