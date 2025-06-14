package main

import (
	"context"
	"github.com/madhavarora03/byte-cloud/internal/config"
	"github.com/madhavarora03/byte-cloud/internal/storage/postgres"
	"log/slog"
	"os"
)

func main() {
	//	Load config
	cfg := config.MustLoad()

	//	Connect to db
	db, err := postgres.Init(cfg)
	if err != nil {
		slog.Error("cannot init postgres: %s", err.Error())
		os.Exit(1)
	}

	defer db.Conn.Close(context.Background())
	slog.Info("connected to postgres", slog.String("env", cfg.Env), slog.String("version", cfg.AppVersion))
}
