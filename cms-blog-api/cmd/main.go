package main

import (
	"log"
	"os"

	"cms-blog-api/config"
	"cms-blog-api/controllers"
	"cms-blog-api/middleware"
	"cms-blog-api/repositories"
	"cms-blog-api/routes"
	"cms-blog-api/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if available
	if err := godotenv.Load(); err != nil {
		log.Println("Notice: .env file not found, using system environment variables or defaults")
	}

	// Initialize MySQL connection and Auto Migration
	config.ConnectDatabase()

	// Initialize Redis connection
	config.ConnectRedis()

	// Dependency Injection setup (Layered Architecture)
	blogRepo := repositories.NewBlogRepository()
	blogService := services.NewBlogService(blogRepo)
	blogController := controllers.NewBlogController(blogService)

	// Initialize Gin router
	router := gin.New()

	// Use Gin Recovery and custom Logger middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())

	// Register REST API routes
	routes.RegisterRoutes(router, blogController)

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting CMS Blog API Server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
