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
	timeouts := config.Timeouts

	srv := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(timeouts.ReadHeaderTimeout) * time.Millisecond,
		ReadTimeout:       time.Duration(timeouts.ReadTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(timeouts.WriteTimeout) * time.Millisecond,
		IdleTimeout:       time.Duration(timeouts.IdleTimeout) * time.Millisecond,
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
