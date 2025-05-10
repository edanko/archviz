package application

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/sourcegraph/conc/pool"
)

// Adapter interface
type Adapter interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// App represents application service
type App struct {
	logger          *slog.Logger
	adapters        []Adapter
	shutdownTimeout time.Duration
}

// New provides new service application
func New(logger *slog.Logger) *App {
	return &App{
		shutdownTimeout: 5 * time.Second, // Default shutdown timeout
		logger:          logger,
	}
}

// AddAdapters adds adapters to application service
func (app *App) AddAdapters(adapters ...Adapter) {
	app.adapters = append(app.adapters, adapters...)
}

// WithShutdownTimeout overrides default shutdown timout
func (app *App) WithShutdownTimeout(timeout time.Duration) {
	app.shutdownTimeout = timeout
}

// Run runs the service application
func (app *App) Run(ctx context.Context) {
	runPool := pool.New().WithContext(ctx).WithFirstError()
	for _, adapter := range app.adapters {
		runPool.Go(func(ctx context.Context) error {
			return adapter.Start(ctx)
		})
	}

	go func() {
		<-ctx.Done()
		app.stop()
	}()

	go func() {
		err := runPool.Wait()
		if err != nil {
			app.logger.Error("run error", "error", err)
			os.Exit(1)
			return
		}
	}()

	<-ctx.Done()
	app.stop()
}

func (app *App) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), app.shutdownTimeout)
	defer cancel()

	app.logger.Info("shutting down...")

	closePool := pool.New().WithContext(ctx)

	for _, adapter := range app.adapters {
		closePool.Go(func(ctx context.Context) error {
			return adapter.Stop(ctx)
		})
	}

	err := closePool.Wait()
	if err != nil {
		app.logger.Error("shutdown error", "error", err)
		os.Exit(1)
		return
	}

	app.logger.Info("gracefully stopped")
}
