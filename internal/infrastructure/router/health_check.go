package router

import (
	"fmt"
	"net/http"
	"order-placement-system/env"
	"order-placement-system/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
)

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   env.ServiceName,
		"version":   env.AppVersion,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func LogRoutes(engine *gin.Engine) {
	routes := engine.Routes()
	log.Infof("Registered routes", log.S("count", fmt.Sprintf("%d", len(routes))))

	for _, route := range routes {
		log.Infof("Route registered",
			log.S("method", route.Method),
			log.S("path", route.Path))
	}
}
