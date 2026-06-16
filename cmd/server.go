package cmd

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/yuyudeqiu/chronicle/internal/config"
	"github.com/yuyudeqiu/chronicle/internal/handler"
	"github.com/yuyudeqiu/chronicle/internal/service"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the chronicle web server",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化配置
		if dataDir != "" {
			config.DataDir = dataDir
		}

		// Initialize database
		service.InitDB(config.GetDBPath())

		// Get current working directory
		dir, _ := os.Getwd()
		log.Printf("Working directory: %s", dir)
		log.Printf("Data directory: %s", config.Load())

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
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
