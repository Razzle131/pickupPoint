package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/handler"
	"github.com/Razzle131/pickupPoint/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const shutdownTimeout = time.Second * 5

//go:generate go tool oapi-codegen -config ../api/schema/apiConfig.yaml ../api/schema/swagger.yaml
//go:generate go tool oapi-codegen -config ../api/schema/modelsConfig.yaml ../api/schema/swagger.yaml

func main() {
	godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runServer(ctx); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx context.Context) error {
	cfg := handler.Config{
		Port:      os.Getenv("SERVER_PORT"),
		DbPort:    os.Getenv("DATABASE_PORT"),
		DbUser:    os.Getenv("DATABASE_USER"),
		DbPasword: os.Getenv("DATABASE_PASSWORD"),
		DbName:    os.Getenv("DATABASE_NAME"),
		DbHost:    os.Getenv("DATABASE_HOST"),
	}

	logger.SetupLogging(slog.LevelInfo)

	slog.Debug("Debugging info enabled")

	dbConString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPasword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	slog.Info("Connecting to db on url: " + dbConString)
	db, err := sqlx.Open("postgres", dbConString)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to connect to database: %s", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if db.Ping() != nil {
		slog.Error("ping result is error")
		os.Exit(1)
	}

	server := handler.NewServer(db)
	mux := http.NewServeMux()
	h := api.HandlerFromMux(server, mux)

	srv := &http.Server{
		Addr:              "0.0.0.0:" + cfg.Port,
		Handler:           h,
		ReadHeaderTimeout: consts.HttpTimeout,
		WriteTimeout:      consts.HttpTimeout,
	}

	slog.Info("Starting server on address " + srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	slog.Info(fmt.Sprintf("Listening on %s", srv.Addr))
	<-ctx.Done()

	slog.Info("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("server shutdown timeout: %s", ctx.Err())
	default:
		slog.Info("Server down")
	}

	return nil
}
