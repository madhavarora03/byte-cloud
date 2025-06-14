package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/madhavarora03/byte-cloud/internal/config"
	"github.com/madhavarora03/byte-cloud/internal/http/handlers/health"
	"github.com/madhavarora03/byte-cloud/internal/storage/postgres"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
