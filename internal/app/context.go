package app

import (
	"sync"

	"go.uber.org/zap"
)

// central store
type AppContext struct {
	logger *zap.Logger

	configMu sync.RWMutex
	config   *ServerConfig
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
func (a *AppContext) GetConfig() *ServerConfig {
	a.configMu.RLock()
	defer a.configMu.RUnlock()
	return a.config
}

// thread-safe config swapping
func (a *AppContext) ReloadConfig(newCfg *ServerConfig) {
	a.configMu.Lock()
	defer a.configMu.Unlock()
	a.config = newCfg
}
