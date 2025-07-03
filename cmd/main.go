package main

import (
	"context"
	"fmt"
	"order-placement-system/env"
	"order-placement-system/pkg/log"
	"syscall"

	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}
	env.LoadEnv()
}

func main() {
	log.Init(env.LogLevel)
	log.Infof("Starting",
		log.S("serviceName", env.ServiceName),
		log.S("version", env.AppVersion))

	gin.SetMode(env.GinMode)
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(gin.Logger())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", env.Port),
		Handler: ginEngine,
	}
	log.Infof("Server Listening on port", env.Port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), env.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown", err)
	}

	log.Info("Server exiting gracefully")
}
