package app

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func WatchConfigFile(ctx *AppContext, path string) {
	log := ctx.Logger()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("Failed to create fsnotify watcher", zap.Error(err))
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		log.Error("Failed to watch config file", zap.String("path", path), zap.Error(err))
		return
	}

	log.Info("Watching for config changes", zap.String("path", path))

	debounce := time.NewTimer(0)
	if !debounce.Stop() {
		<-debounce.C
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// handle file writes
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Info("Detected config file change")

				debounce.Reset(3000 * time.Millisecond)

				go func() {
					<-debounce.C
					newCfg, err := LoadConfig(path)
					if err != nil {
						log.Error("Failed to reload config", zap.Error(err))
						log.Info("Continuing with previous config")
						return
					}

					if err := ValidateConfig(newCfg); err != nil {
						log.Error("Invalid configuration", zap.Error(err))
						log.Info("Continuing with previous config")
						return
					}

					ctx.ReloadConfig(newCfg)
				}()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("Watcher error", zap.Error(err))
		}
	}
}
