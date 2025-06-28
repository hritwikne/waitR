package main

import (
	"waitr/internal/app"
)

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		panic("Failed to load server configuration: " + err.Error())
	}

	logger, err := app.SetupLogging()
	if err != nil {
		panic("Cannot create logger: " + err.Error())
	}
	defer logger.Sync()

	ctx := app.NewAppContext(config, logger)

	log := ctx.Logger()
	log.Info("Config file succesfully loaded and parsed")

	app.Start(ctx)
}
