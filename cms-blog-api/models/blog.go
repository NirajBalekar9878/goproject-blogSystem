package models

import (
	"time"
)

// Blog represents a blog post in our CMS.
// GORM tags define database column properties (primaryKey, type, not null, index).
// JSON tags define the JSON key names returned in API responses.
// binding:"required" instructs Gin to validate that these fields are present in requests.
type Blog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title" binding:"required"`
	Content   string    `gorm:"type:text;not null" json:"content" binding:"required"`
	Author    string    `gorm:"type:varchar(100);not null" json:"author" binding:"required"`
	Category  string    `gorm:"type:varchar(100);not null;index" json:"category" binding:"required"`
	Status    string    `gorm:"type:varchar(50);not null;index" json:"status" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
