package main

import (
	"log"
	"os"
	"path/filepath"

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

	// Get the directory where the executable is located
	execDir := filepath.Dir(os.Args[0])
	// Also try current working directory as fallback
	cwd, _ := os.Getwd()
	frontendDist := filepath.Join(cwd, "frontend", "dist")
	
	// If not found in cwd, try relative to executable
	if _, err := os.Stat(frontendDist); err != nil {
		frontendDist = filepath.Join(execDir, "..", "frontend", "dist")
	}

	// Serve the frontend build output for static assets
	r.Static("/assets", filepath.Join(frontendDist, "assets"))

	// Serve index.html for all other routes (SPA fallback)
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(frontendDist, "index.html"))
	})
	// Start server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
