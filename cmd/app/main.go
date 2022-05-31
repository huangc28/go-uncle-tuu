package main

import (
	"context"
	"fmt"
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/app"
	"huangc28/go-ios-iap-vendor/internal/app/deps"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func init() {
	config.InitConfig()
	db.InitDB()
	if err := deps.Get().Run(); err != nil {
		log.Fatalf("failed to initialize dependency container %s", err.Error())
	}
}

func main() {
	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		db := db.GetDB()

		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, struct {
			OK bool
		}{
			true,
		})
	})

	app.StartApp(r)
	conf := config.GetAppConf()

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", conf.APIPort),

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
