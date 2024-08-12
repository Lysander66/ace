package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lysander66/ace/internal/controller"
	"github.com/Lysander66/zephyr/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetDefault(logger.New(slog.LevelInfo))

	gin.SetMode(gin.DebugMode)
	g := gin.New()
	pprof.Register(g)
	g.Use(cors.Default(), gin.Logger(), gin.Recovery())
	g.ForwardedByClientIP = true

	controller.SetRouter(g)

	srv := &http.Server{Addr: ":8180", Handler: g}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("ListenAndServe")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Debug("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "err", err)
	}
	slog.Info("Server exiting")
}
