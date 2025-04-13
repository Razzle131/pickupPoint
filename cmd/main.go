package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/handler"
	"github.com/Razzle131/pickupPoint/pkg/logger"
	"github.com/joho/godotenv"
)

//go:generate go tool oapi-codegen -config ../api/schema/apiConfig.yaml ../api/schema/swagger.yaml
//go:generate go tool oapi-codegen -config ../api/schema/modelsConfig.yaml ../api/schema/swagger.yaml

func main() {
	godotenv.Load()

	cfg := handler.Config{
		Port:      os.Getenv("SERVER_PORT"),
		DbPort:    os.Getenv("DATABASE_PORT"),
		DbUser:    os.Getenv("DATABASE_USER"),
		DbPasword: os.Getenv("DATABASE_PASSWORD"),
		DbName:    os.Getenv("DATABASE_NAME"),
		DbHost:    os.Getenv("DATABASE_HOST"),
	}

	logger.SetupLogging(slog.LevelDebug)

	slog.Debug("Debugging info enabled")

	dbConString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPasword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	slog.Debug(dbConString)

	slog.Info("Starting db", slog.String("DSN", os.Getenv("POSTGRES_CONN")))

	// dbConn, err := db.New(context.Background(), cfg.DSN)
	// if err != nil {
	// 	slog.Error("Failed to connect to database", slog.String("error", err.Error()))
	// 	os.Exit(1)
	// }
	// defer dbConn.Close()

	server := handler.NewServer()

	r := http.NewServeMux()

	h := api.HandlerFromMux(server, r)

	srv := &http.Server{
		Addr:              "0.0.0.0:" + cfg.Port,
		Handler:           h,
		ReadHeaderTimeout: consts.HttpTimeout,
		WriteTimeout:      consts.HttpTimeout,
	}
	slog.Info("Starting server on address " + srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error(fmt.Sprintf("Failed to start server: %s", err.Error()))
		os.Exit(1)
	}
}
