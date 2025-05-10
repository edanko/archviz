package application

import (
	"context"
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/go-faster/errors"
)

// DebugAdapter ./...
type DebugAdapter struct {
	*http.Server
}

// NewDebugAdapter provides new debug adapter.
// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
// /debug/vars - Added to the default mux by importing the expvar package.
func NewDebugAdapter(address string) *DebugAdapter {
	return &DebugAdapter{
		&http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			Addr:         address,
			Handler:      http.DefaultServeMux,
		},
	}
}

// Start starts http application adapter.
func (adapter *DebugAdapter) Start(_ context.Context) error {
	err := adapter.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stops http application adapter.
func (adapter *DebugAdapter) Stop(ctx context.Context) error {
	return adapter.Shutdown(ctx)
}
