package main

import (
	"context"
	"fmt"
	"net/http"
	"order-placement-system/env"
	"order-placement-system/internal/adapter/handler"
	"order-placement-system/internal/adapter/presenter"
	"order-placement-system/internal/infrastructure/middleware"
	"order-placement-system/internal/infrastructure/router"
	"order-placement-system/internal/usecases/implementation"
	"order-placement-system/pkg/log"
	"order-placement-system/pkg/utils/parser"
	"os"
	"os/signal"
	"syscall"

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
	engine := gin.New()

	middleware.Setup(engine)
	router.SetupHealthCheck(engine)

	productParser := parser.NewProductParser()

	complementaryCalculator := implementation.NewComplementaryCalculator()

	orderProcessor := implementation.NewOrderProcessor(
		productParser,
		complementaryCalculator,
	)

	orderPresenter := presenter.NewOrderPresenter()

	orderHandler := handler.NewOrderHandler(orderProcessor, orderPresenter)

	router.OrderPlacementV1Routes(engine, orderHandler)

	router.LogRoutes(engine)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", env.Port),
		Handler: engine,
	}

	go func() {
		log.Infof("Starting HTTP server", log.S("port", env.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server", log.E(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), env.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown", log.E(err))
	}

	log.Info("Server exited gracefully")

}
