package routes

import (
	"cms-blog-api/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all blog endpoints to the Gin engine.
func RegisterRoutes(router *gin.Engine, blogController *controllers.BlogController) {
	blogGroup := router.Group("/blogs")
	{
		blogGroup.POST("", blogController.CreateBlog)
		blogGroup.GET("", blogController.GetAllBlogs)
		blogGroup.GET("/:id", blogController.GetBlogByID)
		blogGroup.PUT("/:id", blogController.UpdateBlog)
		blogGroup.DELETE("/:id", blogController.DeleteBlog)
		blogGroup.GET("/category/:category", blogController.GetBlogsByCategory)
	}
}
