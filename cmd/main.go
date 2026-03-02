package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/overm-app/api-auth/internal/interface/api"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if err := godotenv.Load(); err != nil {
		sugar.Warn("No .env file found, using environment variables or defaults")
	}

	timezone := os.Getenv("SERVER_TIMEZONE")
	if timezone == "" {
		timezone = "UTC"
		sugar.Warn("SERVER_TIMEZONE not set, defaulting to UTC")
	}
	loc, err := time.LoadLocation(timezone)
	sugar.Infof("Setting server timezone to %s", timezone)
	if err != nil {
		sugar.Errorf("Failed to set timezone %s: %v", timezone, err)
		loc = time.UTC
	}
	time.Local = loc

	r := api.SetupRouter()

	port := os.Getenv("SERVER_PORT")

	sugar.Infof("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		sugar.Errorw("Failed to start server", 
		"error", err,
		"port", port,
		"timestamp", time.Now().Format(time.RFC3339),
		)
		os.Exit(1)
	}
}
