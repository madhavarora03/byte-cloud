package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/madhavarora03/byte-cloud/internal/config"
	"github.com/madhavarora03/byte-cloud/internal/db"
	"github.com/madhavarora03/byte-cloud/internal/http/handlers/health"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	//go:embed migrations/*.sql
	migrationFiles embed.FS
)

func main() {
	//	Load config
	cfg := config.MustLoad()

	//	Connect to db
	ctx := context.Background()
	dbConfig := &db.Config{
		MigrationsTable:       "schema_migrations",
		DatabaseName:          "your_db_name", // will be auto-detected if blank
		SchemaName:            "public",       // will be auto-detected if blank
		StatementTimeout:      10 * time.Second,
		MultiStatementEnabled: true,
	}

	pg, err := db.NewPostgres(ctx, dbConfig, cfg.DbUri)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pg.Db.Close()

	slog.Info("connected to db", slog.String("env", cfg.Env), slog.String("version", cfg.AppVersion))

	//	run migrations
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		log.Fatalf("Failed to read embedded migrations: %v", err)
	}

	driver, err := db.WithInstance(pg.Db, dbConfig)
	if err != nil {
		log.Fatalf("Failed to wrap db: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "pgx5", driver)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	slog.Info("migrations ran successfully")
	//	setup router
	if cfg.Env == "production" {
		//	in production use release mode
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		slog.Info("endpoint", slog.String("http-method", httpMethod), slog.String("path", absolutePath), slog.String("handler", handlerName), slog.String("num-handlers", fmt.Sprint(nuHandlers)))
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//	create API version group
	v1 := router.Group("/api/v1")

	//	register routes for different services
	health.SetupRouter(v1)

	//	start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Warn("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Warn("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	slog.Info("timeout of 5 seconds.")
	slog.Info("Server exiting")
}
