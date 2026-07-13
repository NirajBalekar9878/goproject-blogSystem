package routes

import (
	"cms-blog-api/controllers"
	"cms-blog-api/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all auth and blog endpoints to the Gin engine.
func RegisterRoutes(router *gin.Engine, blogController *controllers.BlogController, authController *controllers.AuthController) {
	// Authentication routes
	authGroup := router.Group("/auth")
	{
		// Public auth endpoints
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)

		// Protected auth endpoints
		protectedAuth := authGroup.Group("")
		protectedAuth.Use(middleware.AuthMiddleware())
		{
			protectedAuth.GET("/profile", authController.GetProfile)
		}
	}

	// Blog routes
	blogGroup := router.Group("/blogs")
	{
		// Public read endpoints
		blogGroup.GET("", blogController.GetAllBlogs)
		blogGroup.GET("/:id", blogController.GetBlogByID)
		blogGroup.GET("/category/:category", blogController.GetBlogsByCategory)

		// Protected write endpoints requiring JWT authentication middleware
		protectedBlogs := blogGroup.Group("")
		protectedBlogs.Use(middleware.AuthMiddleware())
		{
			protectedBlogs.POST("", blogController.CreateBlog)
			protectedBlogs.PUT("/:id", blogController.UpdateBlog)
			protectedBlogs.DELETE("/:id", blogController.DeleteBlog)
		}
	}
}
