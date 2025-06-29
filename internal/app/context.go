package app

import (
	"context"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// central store
type AppContext struct {
	logger *zap.Logger

	configMu sync.RWMutex
	config   *ServerConfig

	serverMu   sync.Mutex
	httpServer *http.Server
}

type LimitedLogger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
}

// store init
func NewAppContext(cfg *ServerConfig, logger *zap.Logger) *AppContext {
	return &AppContext{
		config: cfg,
		logger: logger,
	}
}

// exposing logger with limited set of methods
func (a *AppContext) Logger() LimitedLogger {
	return a.logger
}

// thread-safe getter
func (ctx *AppContext) GetConfig() *ServerConfig {
	ctx.configMu.RLock()
	defer ctx.configMu.RUnlock()
	return ctx.config
}

// thread-safe config swapping
func (ctx *AppContext) ReloadConfig(newCfg *ServerConfig) {
	ctx.configMu.Lock()
	oldConfig := ctx.config
	ctx.config = newCfg
	ctx.configMu.Unlock()

	log := ctx.Logger()
	log.Info("Configuration hot-reloaded successfully")

	if oldConfig.ListenAddress != newCfg.ListenAddress {
		log.Info("Listen address changed, Restarting HTTP Server...")

		oldServer := ctx.GetServer()
		if oldServer == nil {
			log.Info("No running server — starting fresh instance")
			StartHTTPServer(ctx)
			return
		}

		go func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := oldServer.Shutdown((shutdownCtx)); err != nil {
				log.Error("Failed to shutdown old server", zap.Error(err))
			} else {
				log.Info("Previous server instance shut down")
			}

			StartHTTPServer(ctx)
		}()
	}
}

func (ctx *AppContext) Shutdown(shutdownCtx context.Context) error {
	ctx.logger.Info("Performing cleanup...")
	if ctx.httpServer != nil {
		return ctx.httpServer.Shutdown(shutdownCtx) // cleanly closes listener & waits for active reqs
	}
	return nil
}

func (ctx *AppContext) SetServer(srv *http.Server) {
	ctx.serverMu.Lock()
	defer ctx.serverMu.Unlock()
	ctx.httpServer = srv
}

func (ctx *AppContext) GetServer() *http.Server {
	ctx.serverMu.Lock()
	defer ctx.serverMu.Unlock()
	return ctx.httpServer
}
