package router

import (
	"order-placement-system/internal/adapter/handler"

	"github.com/gin-gonic/gin"
)

func SetupHealthCheck(engine *gin.Engine) {
	engine.GET("/health", healthCheck)
}

func OrderPlacementV1Routes(engine *gin.Engine, order handler.OrderHandlerInterface) {
	v1 := engine.Group("/api/v1")

	orders := v1.Group("/orders")
	{
		orders.POST("/process", order.ProcessOrders)
	}
}
