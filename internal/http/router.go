package http

import (
	nethttp "net/http"

	"zapier-project/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(deploymentsHandler *handlers.DeploymentsHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(nethttp.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/deployments", deploymentsHandler.List)
	router.GET("/deployments/:id", deploymentsHandler.GetByID)

	return router
}
