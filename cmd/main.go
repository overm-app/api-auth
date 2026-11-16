package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/overm-app/api-auth/internal/config"
	"github.com/overm-app/api-auth/internal/infrastructure/db"
	"github.com/overm-app/api-auth/internal/infrastructure/repository"
	"github.com/overm-app/api-auth/internal/infrastructure/service"
	"github.com/overm-app/api-auth/internal/interface/api"
	"github.com/overm-app/api-auth/internal/interface/api/handlers"
	"github.com/overm-app/api-auth/internal/usecase"
)

func main() {
	// Setup logger
	sugar, level := newLogger()
	defer sugar.Sync()

	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "/etc/config/config.yaml"
	}
	cfg, err := config.Load(path)
	if err != nil {
		sugar.Errorw("invalid configuration", "error", err)
		os.Exit(1)
	}

	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(cfg.Log.Level)); err != nil {
		sugar.Errorw("invalid log.level", "value", cfg.Log.Level, "error", err)
		os.Exit(1)
	}
	level.SetLevel(lvl)

	sugar = sugar.With("service", cfg.Env)
	sugar = sugar.With("env", cfg.Server.Name)

	port := cfg.Server.Port

	// Connect to database
	dbCfg := db.Config{
		Host:     cfg.Database.DBHost,
		Port:     cfg.Database.DBPort,
		User:     cfg.Database.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.DBSslMode,
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

	jwtSecret := []byte(cfg.JWTSecret)
	accessTokenTTL := cfg.Server.Auth.AccessTokenTTL
	refreshTokenTTL := cfg.Server.Auth.RefreshTokenTTL

	// Initialize services
	jwtService := service.NewJWTService(jwtSecret, accessTokenTTL)
	
	// Initialize use cases
	loginUseCase := usecase.NewLoginUseCase(userRepo, tokenRepo, jwtService, refreshTokenTTL)
	registerUseCase := usecase.NewRegisterUseCase(userRepo, tokenRepo, jwtService)

	cookieConfig := handlers.CookieConfig{
		Secure: cfg.Server.Cookie.Secure,
		Domain: cfg.Server.Cookie.Domain,
		AccessTTL: cfg.Server.Auth.AccessTokenTTL,
		RefreshTTL: cfg.Server.Auth.RefreshTokenTTL,
	}

	// Initialize API handlers
	cookieManager := handlers.NewCookieManager(cookieConfig) 
	authHandler := handlers.NewAuthHandler(loginUseCase, cookieManager, sugar,)
	userHandler := handlers.NewUserHandler(registerUseCase, cookieManager, sugar,)
	healthHandler := handlers.NewHealthHandler(postgresDB, sugar)


	// Setup and start server
	r := api.NewRouter(authHandler, userHandler, healthHandler, jwtService, sugar)
	engine := r.SetupRouter(sugar)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		sugar.Infow("Starting server", "port", port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		sugar.Errorw("Server error", "error", err)
		os.Exit(1)
	case <-ctx.Done():
	}

	stop()

	sugar.Infow("Shutting down server...")
	healthHandler.SetShuttingDown()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		sugar.Errorw("Graceful shutdown failed", "error", err)
	}

	sugar.Infow("Server stopped")
}

func newLogger() (*zap.SugaredLogger, zap.AtomicLevel) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		zapcore.ISO8601TimeEncoder(t.UTC(), enc)
	}
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := zap.Config{
		Level:            level,
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(zap.WithCaller(true), zap.AddStacktrace(zapcore.FatalLevel))
	if err != nil {
		panic(err)
	}
	return logger.Sugar(), level
}