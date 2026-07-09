package controllers

import (
	"errors"
	"strconv"

	"cms-blog-api/models"
	"cms-blog-api/services"
	"cms-blog-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BlogController struct {
	service *services.BlogService
}

func NewBlogController(service *services.BlogService) *BlogController {
	return &BlogController{service: service}
}

// CreateBlog handles POST /blogs
func (bc *BlogController) CreateBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	if err := bc.service.CreateBlog(c.Request.Context(), &blog); err != nil {
		utils.RespondInternalError(c, "Failed to create blog: "+err.Error())
		return
	}

	utils.RespondCreated(c, "Blog created successfully", blog)
}

// GetAllBlogs handles GET /blogs
func (bc *BlogController) GetAllBlogs(c *gin.Context) {
	blogs, err := bc.service.GetAllPublishedBlogs(c.Request.Context())
	if err != nil {
		utils.RespondInternalError(c, "Failed to retrieve blogs: "+err.Error())
		return
	}

	utils.RespondOK(c, "Blogs fetched successfully", blogs)
}

// GetBlogByID handles GET /blogs/:id
func (bc *BlogController) GetBlogByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid blog ID parameter")
		return
	}

	blog, err := bc.service.GetBlogByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondNotFound(c, "Blog not found")
			return
		}
		utils.RespondInternalError(c, "Failed to retrieve blog: "+err.Error())
		return
	}

	utils.RespondOK(c, "Blog fetched successfully", blog)
}

// UpdateBlog handles PUT /blogs/:id
func (bc *BlogController) UpdateBlog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid blog ID parameter")
		return
	}

	// Verify existing blog exists
	existing, err := bc.service.GetBlogByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondNotFound(c, "Blog not found")
			return
		}
		utils.RespondInternalError(c, "Failed to check existing blog: "+err.Error())
		return
	}

	var updatedInput models.Blog
	if err := c.ShouldBindJSON(&updatedInput); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	updatedInput.ID = existing.ID
	updatedInput.CreatedAt = existing.CreatedAt

	if err := bc.service.UpdateBlog(c.Request.Context(), &updatedInput); err != nil {
		utils.RespondInternalError(c, "Failed to update blog: "+err.Error())
		return
	}

	utils.RespondOK(c, "Blog updated successfully", updatedInput)
}

// DeleteBlog handles DELETE /blogs/:id
func (bc *BlogController) DeleteBlog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid blog ID parameter")
		return
	}

	if err := bc.service.DeleteBlog(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondNotFound(c, "Blog not found")
			return
		}
		utils.RespondInternalError(c, "Failed to delete blog: "+err.Error())
		return
	}

	utils.RespondOK(c, "Blog deleted successfully", nil)
}

// GetBlogsByCategory handles GET /blogs/category/:category
func (bc *BlogController) GetBlogsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		utils.RespondBadRequest(c, "Category parameter is required")
		return
	}

	blogs, err := bc.service.GetBlogsByCategory(c.Request.Context(), category)
	if err != nil {
		utils.RespondInternalError(c, "Failed to retrieve blogs by category: "+err.Error())
		return
	}

	utils.RespondOK(c, "Blogs fetched successfully", blogs)
}
