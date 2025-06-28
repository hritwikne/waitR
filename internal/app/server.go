package app

import (
	"net/http"

	"go.uber.org/zap"
)

func Start(ctx *AppContext) {
	log := ctx.Logger()
	address := ctx.GetConfig().ListenAddress

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Received GET / request")
	})

	log.Info("Starting WaitR server", zap.String("address", address))
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("Failed to start WaitR server", zap.Error(err))
	}

}
