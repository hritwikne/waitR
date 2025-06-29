package app

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func StartHTTPServer(ctx *AppContext) {
	log := ctx.Logger()
	config := ctx.GetConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info(
			"Received request",
			zap.String("method", r.Method), zap.String("url", r.URL.String()),
		)
		w.Write([]byte("Hello from WaitR"))
	})

	address := config.ListenAddress
	srv := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx.httpServer = srv

	go func() {
		log.Info("Starting WaitR server", zap.String("address", address))
		// Ignore http.ErrServerClosed error since it is returned on graceful shutdown
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error", zap.Error(err))
		}
	}()
}
