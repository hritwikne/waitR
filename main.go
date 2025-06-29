package main

import (
	"fmt"
	"waitr/internal/app"
)

func main() {
	pathToConfig := "config/server.json"
	config, err := app.LoadConfig(pathToConfig)
	if err != nil {
		panic("Failed to load server configuration: " + err.Error())
	}

	if err := app.ValidateConfig(config); err != nil {
		panic(fmt.Errorf("invalid server configuration: %w", err))
	}

	logger, err := app.SetupLogging()
	if err != nil {
		panic("Cannot create logger: " + err.Error())
	}
	defer logger.Sync()
	logger.Info("Config file succesfully loaded and parsed")

	// central store init
	ctx := app.NewAppContext(config, logger)

	// goroutines
	app.StartHTTPServer(ctx)
	go app.WatchConfigFile(ctx, pathToConfig)

	app.HandleShutdown(ctx) // waiting for signal
}
