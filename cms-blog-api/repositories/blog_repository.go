package repositories

import (
	"errors"

	"cms-blog-api/config"
	"cms-blog-api/models"

	"gorm.io/gorm"
)

// BlogRepository handles database operations for blogs.
type BlogRepository struct{}

func NewBlogRepository() *BlogRepository {
	return &BlogRepository{}
}

// Create inserts a new blog into MySQL.
func (r *BlogRepository) Create(blog *models.Blog) error {
	return config.DB.Create(blog).Error
}

// GetAllPublished returns all published blogs from MySQL.
func (r *BlogRepository) GetAllPublished() ([]models.Blog, error) {
	var blogs []models.Blog
	err := config.DB.Where("status = ?", "published").Find(&blogs).Error
	return blogs, err
}

// GetByID finds a blog by its ID. Returns gorm.ErrRecordNotFound if not found.
func (r *BlogRepository) GetByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	err := config.DB.First(&blog, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &blog, nil
}

// Update updates existing blog fields in MySQL.
func (r *BlogRepository) Update(blog *models.Blog) error {
	return config.DB.Save(blog).Error
}

// Delete removes a blog by its ID.
func (r *BlogRepository) Delete(id uint) error {
	result := config.DB.Delete(&models.Blog{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetByCategory returns all blogs belonging to the specified category.
func (r *BlogRepository) GetByCategory(category string) ([]models.Blog, error) {
	var blogs []models.Blog
	err := config.DB.Where("category = ?", category).Find(&blogs).Error
	return blogs, err
}
