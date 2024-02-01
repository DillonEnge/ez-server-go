package service

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DillonEnge/ez-server-go/middleware"
	"golang.org/x/sync/errgroup"
)

func MakeServe(addr string, mux *http.ServeMux, bundles []*middleware.ContextBundle) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	s := &http.Server{
		BaseContext:    func(_ net.Listener) context.Context { return ctx },
		Addr:           addr,
		Handler:        middleware.Context(middleware.Logger(mux), bundles...),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	g.Go(func() error {
		slog.Info(fmt.Sprintf("Starting server on %s", addr))
		return s.ListenAndServe()
	})

	<-ctx.Done()
	slog.Info("Shutting down...")
	s.Shutdown(ctx)
}
