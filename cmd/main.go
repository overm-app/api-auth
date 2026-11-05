package main

import (
	"os"
	"strconv"
	"time"

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
	// Setup logger
	sugar := setupLogger()
	defer sugar.Sync()

	currentEnv := os.Getenv("ENV")
	if currentEnv != "develop" && currentEnv != "QA" && currentEnv != "prod" {
		if err := godotenv.Load(); err != nil {
			sugar.Errorw("No .env file found, relying on system environment variables")
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
		sugar.Warnw("SERVER_PORT not set, defaulting to 8081")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		sugar.Warnw("Invalid DB_PORT, defaulting to 5432", "error", err)
		dbPort = 5432
	}

	jwtExpiration := 24 * time.Hour
	if exp := os.Getenv("JWT_EXPIRATION_HOURS"); exp != "" {
		if parsed, err := time.ParseDuration(exp + "h"); err == nil {
			jwtExpiration = parsed
		} else {
			sugar.Warnw("Invalid JWT_EXPIRATION_HOURS, defaulting to 24h", "value", exp)
		}
	}

	// Connect to database
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(postgresDB)
	tokenRepo := repository.NewTokenRepository(postgresDB)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		sugar.Errorw("JWT_SECRET is not set")
		os.Exit(1)
	}

	// Initialize services
	jwtService := service.NewJWTService(jwtSecret, jwtExpiration)
	
	// Initialize use cases
	loginUseCase := usecase.NewLoginUseCase(userRepo, tokenRepo, jwtService)
	registerUseCase := usecase.NewRegisterUseCase(userRepo, tokenRepo, jwtService)

	// Initialize API handlers
	cookieCfg := handlers.NewCookieConfig()
	authHandler := handlers.NewAuthHandler(loginUseCase,cookieCfg,sugar,)
	userHandler := handlers.NewUserHandler(registerUseCase,cookieCfg,sugar,)
	healthHandler := handlers.NewHealthHandler(postgresDB, sugar)


	// Setup and start server
	r := api.NewRouter(authHandler, userHandler, healthHandler, jwtService, sugar)
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

func setupLogger() *zap.SugaredLogger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		zapcore.ISO8601TimeEncoder(t.UTC(), enc)
	}
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

	logger = logger.WithOptions(zap.WithCaller(true), zap.AddStacktrace(zapcore.FatalLevel))
	sugar := logger.Sugar()
	return sugar
}