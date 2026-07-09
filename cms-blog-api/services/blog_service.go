package services

import (
	"context"
	"fmt"
	"log"

	"cms-blog-api/cache"
	"cms-blog-api/models"
	"cms-blog-api/repositories"
)

type BlogService struct {
	repo *repositories.BlogRepository
}

func NewBlogService(repo *repositories.BlogRepository) *BlogService {
	return &BlogService{repo: repo}
}

// CreateBlog creates a new blog post and invalidates relevant caches.
func (s *BlogService) CreateBlog(ctx context.Context, blog *models.Blog) error {
	err := s.repo.Create(blog)
	if err != nil {
		return err
	}

	// Cache Invalidation
	cache.InvalidateBlogCaches(ctx, blog.ID, blog.Category)
	return nil
}

// GetAllPublishedBlogs returns published blogs, checking Redis cache first.
func (s *BlogService) GetAllPublishedBlogs(ctx context.Context) ([]models.Blog, error) {
	var blogs []models.Blog

	// Try Redis cache first
	err := cache.GetCache(ctx, cache.KeyAllBlogs, &blogs)
	if err == nil {
		log.Println("Cache HIT: GET /blogs")
		return blogs, nil
	}

	log.Println("Cache MISS: GET /blogs - Fetching from MySQL")
	blogs, err = s.repo.GetAllPublished()
	if err != nil {
		return nil, err
	}

	// Store in Redis with 5 min expiry
	_ = cache.SetCache(ctx, cache.KeyAllBlogs, blogs, cache.CacheTTL)
	return blogs, nil
}

// GetBlogByID returns a blog by ID, checking Redis cache first.
func (s *BlogService) GetBlogByID(ctx context.Context, id uint) (*models.Blog, error) {
	var blog models.Blog
	cacheKey := fmt.Sprintf("%s%d", cache.KeyBlogByIDPrefix, id)

	err := cache.GetCache(ctx, cacheKey, &blog)
	if err == nil {
		log.Printf("Cache HIT: GET /blogs/%d\n", id)
		return &blog, nil
	}

	log.Printf("Cache MISS: GET /blogs/%d - Fetching from MySQL\n", id)
	foundBlog, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Store in Redis with 5 min expiry
	_ = cache.SetCache(ctx, cacheKey, foundBlog, cache.CacheTTL)
	return foundBlog, nil
}

// UpdateBlog updates an existing blog post and invalidates relevant caches.
func (s *BlogService) UpdateBlog(ctx context.Context, blog *models.Blog) error {
	err := s.repo.Update(blog)
	if err != nil {
		return err
	}

	// Cache Invalidation
	cache.InvalidateBlogCaches(ctx, blog.ID, blog.Category)
	return nil
}

// DeleteBlog removes a blog post and invalidates relevant caches.
func (s *BlogService) DeleteBlog(ctx context.Context, id uint) error {
	// First fetch to know category if we want precise invalidation, or delete directly
	existing, _ := s.repo.GetByID(id)
	category := ""
	if existing != nil {
		category = existing.Category
	}

	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Cache Invalidation
	cache.InvalidateBlogCaches(ctx, id, category)
	return nil
}

// GetBlogsByCategory returns blogs belonging to a specific category, checking Redis cache first.
func (s *BlogService) GetBlogsByCategory(ctx context.Context, category string) ([]models.Blog, error) {
	var blogs []models.Blog
	cacheKey := fmt.Sprintf("%s%s", cache.KeyCategoryPrefix, category)

	err := cache.GetCache(ctx, cacheKey, &blogs)
	if err == nil {
		log.Printf("Cache HIT: GET /blogs/category/%s\n", category)
		return blogs, nil
	}

	log.Printf("Cache MISS: GET /blogs/category/%s - Fetching from MySQL\n", category)
	blogs, err = s.repo.GetByCategory(category)
	if err != nil {
		return nil, err
	}

	// Store in Redis with 5 min expiry
	_ = cache.SetCache(ctx, cacheKey, blogs, cache.CacheTTL)
	return blogs, nil
}
