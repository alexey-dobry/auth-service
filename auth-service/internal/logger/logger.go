package logger

import (
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

const logDirPath = "./logs"

func NewLogger(env string) *zap.SugaredLogger {
	os.MkdirAll(logDirPath, os.ModePerm)

	logFile, err := os.OpenFile(filepath.Join(logDirPath, "auth.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to create log file: %s", err.Error())
	}

	var logCfg zapcore.EncoderConfig
	switch env {
	case EnvLocal:
		logCfg = zap.NewDevelopmentEncoderConfig()
	case EnvDev:
		logCfg = zap.NewDevelopmentEncoderConfig()
	case EnvProd:
		logCfg = zap.NewProductionEncoderConfig()
	}

	logCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(logCfg)

	logCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(logCfg)

	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	}

	core := zapcore.NewTee(cores...)

	return zap.New(core).Sugar()
}
