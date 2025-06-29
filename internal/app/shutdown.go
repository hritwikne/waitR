package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func HandleShutdown(ctx *AppContext) {
	log := ctx.Logger()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	log.Info("Shutdown signal received", zap.String("signal", sig.String()))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ctx.Shutdown(shutdownCtx); err != nil {
		log.Error("Error during shutdown", zap.Error(err))
	} else {
		log.Info("Shutdown complete")
	}
}
