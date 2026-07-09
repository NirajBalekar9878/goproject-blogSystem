package config

import (
	"fmt"
	"log"
	"os"

	"cms-blog-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase initializes the MySQL connection using GORM and performs auto-migration.
func ConnectDatabase() *gorm.DB {
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "secret")
	dbName := getEnv("DB_NAME", "cms_blog_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
	}

	log.Println("MySQL Database connected successfully")

	// Perform Auto Migration for the Blog model
	err = DB.AutoMigrate(&models.Blog{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database schema: %v", err)
	}
	log.Println("Database auto-migration completed successfully")

	return DB
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
