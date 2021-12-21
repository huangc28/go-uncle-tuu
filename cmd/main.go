package main

import (
	"context"
	"fmt"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func init() {
	db.InitDB()
}

func main() {
	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			OK bool
		}{
			true,
		})
	})

	app.StartApp(r)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", 3009),

		// Good practice: enforce timeouts for servers created.
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {

		log.Infof("listen on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("graceful shutdown...")

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown with error: %s", err.Error())
	}

	log.Info("shutdown complete..")
}
