package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/overm-app/api-auth/internal/infrastructure/db"
	"github.com/overm-app/api-auth/internal/infrastructure/repository"
	"github.com/overm-app/api-auth/internal/infrastructure/service"
	"github.com/overm-app/api-auth/internal/interface/api"
	"github.com/overm-app/api-auth/internal/interface/api/handlers"
	"github.com/overm-app/api-auth/internal/usecase"
)

func main() {
	sugar := setupLogger()
	defer sugar.Sync()

	if err := godotenv.Load(); err != nil {
		sugar.Warnw("No .env file found, using environment variables or defaults")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
		sugar.Warnw("SERVER_PORT not set, defaulting to 8081")
	}

	setupTimezone(sugar)

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		sugar.Warnw("Invalid DB_PORT, defaulting to 5432", "error", err)
		dbPort = 5432
	}

	dbCfg := db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	postgresDB, err := db.NewPostgresConnection(dbCfg, sugar)
	if err != nil {
		sugar.Errorw("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer postgresDB.Close()

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbCfg.User, dbCfg.Password, dbCfg.Host, strconv.Itoa(dbCfg.Port), dbCfg.DBName, dbCfg.SSLMode,
	)
	if os.Getenv("ENV") != "production" {
    	runMigrations(databaseURL, sugar)
	}

	userRepo := repository.NewUserRepository(postgresDB)
	tokenRepo := repository.NewTokenRepository(postgresDB)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		sugar.Errorw("JWT_SECRET is not set")
		os.Exit(1)
	}

	jwtExpiration := 24 * time.Hour
	if exp := os.Getenv("JWT_EXPIRATION_HOURS"); exp != "" {
		if parsed, err := time.ParseDuration(exp + "h"); err == nil {
			jwtExpiration = parsed
		} else {
			sugar.Warnw("Invalid JWT_EXPIRATION_HOURS, defaulting to 24h", "value", exp)
		}
	}

	jwtService := service.NewJWTService(jwtSecret, jwtExpiration)
	
	loginUseCase := usecase.NewLoginUseCase(userRepo, tokenRepo, jwtService)

	cookieCfg := handlers.NewCookieConfig()
	authHandler := handlers.NewAuthHandler(
		loginUseCase,
		cookieCfg,
		sugar,
	)

	r := api.NewRouter(authHandler, jwtService, sugar)
	engine := r.SetupRouter(sugar)

	sugar.Infow("Starting server", "port", port)

	if err := engine.Run(":" + port); err != nil {
		sugar.Errorw("Failed to start server",
			"error", err,
			"port", port,
		)
		os.Exit(1)
	}
}

func runMigrations(databaseURL string, sugar *zap.SugaredLogger) {
    m, err := migrate.New("file://migrations", databaseURL)
    if err != nil {
        sugar.Errorw("Failed to create migrator", "error", err)
        os.Exit(1)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        sugar.Errorw("Failed to run migrations", "error", err)
        os.Exit(1)
    }

    sugar.Infow("Migrations applied successfully")
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
