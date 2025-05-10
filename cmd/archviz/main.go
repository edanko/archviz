package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
	"oss.terrastruct.com/d2/lib/log"

	"github.com/edanko/archviz/cmd/archviz/internal/adapters"
	"github.com/edanko/archviz/cmd/archviz/internal/app/service"
	"github.com/edanko/archviz/cmd/archviz/internal/ports"
	"github.com/edanko/archviz/pkg/application"
	"github.com/edanko/archviz/pkg/http"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := adapters.ReadConfig("./config.yaml")
	if err != nil {
		slog.Error("Failed to read config!", "details", err.Error())
		os.Exit(1)
	}

	logger := slog.New(
		log.NewPrettyHandler(
			slog.NewTextHandler(os.Stdout, nil),
		),
	)
	logger = logger.With(slog.String("app", "archviz"))
	ctx = log.With(ctx, logger)

	if _, err := maxprocs.Set(); err != nil {
		logger.Error(
			"automaxprocs failed",
			slog.String("error", err.Error()),
		)
	}
	logger.Info(
		"startup",
		slog.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)),
	)

	app, cleanup := service.NewApplication(ctx, cfg)
	defer cleanup()

	httpHandlers := ports.NewHTTPServer(app)

	fiberApp := runServer(httpHandlers, logger)

	a := application.New(logger)
	a.AddAdapters(
		http.NewFiberAdapter(fiberApp, cfg.HTTP.Address()),
		application.NewDebugAdapter(cfg.Debug.Address()),
	)

	a.WithShutdownTimeout(cfg.App.ShutdownTimeout)
	a.Run(ctx)
}
