package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/overm-app/api-auth/internal/interface/api"
)

func main() {
	sugar := setupLogger()

	if err := godotenv.Load(); err != nil {
		sugar.Warnw("No .env file found, using environment variables or defaults")
	}

	setupTimezone(sugar)

	r := api.SetupRouter(sugar)

	port := os.Getenv("SERVER_PORT")

	sugar.Infow("Starting server","port", port)

	if err := r.Run(":" + port); err != nil {
		sugar.Errorw("Failed to start server",
			"error", err,
			"port", port,
		)
		os.Exit(1)
	}
}

func setupLogger() *zap.SugaredLogger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		level = zapcore.InfoLevel
	}

	logger, _ := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()

	sugar := logger.Sugar()
	return sugar
}

func setupTimezone(sugar *zap.SugaredLogger) {
	timezone := os.Getenv("SERVER_TIMEZONE")
	if timezone == "" {
		timezone = "UTC"
	}
	loc, err := time.LoadLocation(timezone)
	sugar.Infow("Setting server timezone", "timezone", timezone)
	if err != nil {
		sugar.Errorw("Failed to set timezone", "timezone", timezone, "error", err)
		loc = time.UTC
	}
	time.Local = loc
}