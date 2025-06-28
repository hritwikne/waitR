package app

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Level      string `json:"level"`
	Mode       string `json:"mode"`
	LogFile    string `json:"log_file"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
}

func SetupLogging() (*zap.Logger, error) {
	pathToConfigFile := "config/logger.json"

	// Open the logger configuration file
	file, err := os.Open(pathToConfigFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open logger config: %w", err)
	}
	defer file.Close()

	// Decode the JSON configuration
	var cfg LoggerConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("invalid logger config: %w", err)
	}

	// Parse log level string into zapcore.Level
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.Level)) // fallback to info if invalid

	// Create atomic level so we can change it later at runtime
	atomicLevel := zap.NewAtomicLevelAt(level)

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		NameKey:      "logger",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeName:   zapcore.FullNameEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	consoleWriter := zapcore.Lock(os.Stdout)
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleWriter, atomicLevel),
		zapcore.NewCore(encoder, fileWriter, atomicLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	return logger, nil
}
