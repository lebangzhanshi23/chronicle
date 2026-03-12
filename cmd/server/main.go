package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yuyudeqiu/chronicle/internal/handler"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

func main() {
	// Initialize database
	service.InitDB("data/app.db")

	// Setup router
	r := gin.Default()

	// Register APIs
	handler.RegisterRoutes(r)

	// Serve the frontend build output for static assets
	r.Static("/assets", "./frontend/dist/assets")

	// Serve index.html for all other routes (SPA fallback)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
	// Start server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
