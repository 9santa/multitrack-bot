package gateway

import (
	"github.com/gin-gonic/gin"
	"multitrack-bot/internal/core"
	"net/http"
)

type GatewayServer struct {
	router          *gin.Engine
	trackingService *core.TrackingService
}

func NewGateway(trackingService *core.TrackingService) *GatewayServer {
	g := &GatewayServer{
		router:          gin.Default(),
		trackingService: trackingService,
	}

	g.setupRoutes()
	return g
}

func (g *GatewayServer) setupRoutes() {
	// REST API
	api := g.router.Group("/api/v1")
	{
		api.POST("/track", g.trackPackage)
		api.GET("/health", g.healthCheck)
	}
}

func (g *GatewayServer) trackPackage(c *gin.Context) {
	var req struct {
		TrackingNumber string `json:"tracking_number" binding:"required"`
		Courier        string `json:"courier,omitempty"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := g.trackingService.Track(c.Request.Context(), req.TrackingNumber, req.Courier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (g *GatewayServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (g *GatewayServer) Run(port string) error {
	return g.router.Run(":" + port)
}
