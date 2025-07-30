package repo

import (
	"fmt"
	"mypage-backend/internal/config"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// data model definitions
type User struct {
	ID         uint   `json:"uid" gorm:"primarykey"`
	Username   string `json:"username" gorm:"unique;not null"`
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"-" gorm:"not null"`
	Role       string `json:"role" gorm:"default:user"`
	Registered int64  `json:"registered" gorm:"autoCreateTime"`
	Status     string `json:"status" gorm:"default:processing"` // e.g., active, inactive, processing, banned, etc.
}

type ArticleMeta struct {
	ID        uint     `json:"id" gorm:"primarykey"`
	Title     string   `json:"title" gorm:"not null"`
	AuthorID  uint     `json:"uid" gorm:"not null"`
	CreatedAt int64    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64    `json:"updated_at" gorm:"autoUpdateTime"`
	Tags      []string `json:"tags" gorm:"type:text"` // Tags stored as JSON
}

type ArticleStatistics struct {
	ID        uint  `json:"id" gorm:"primarykey"`
	Liked     uint  `json:"liked" gorm:"default:0"`
	Disliked  uint  `json:"disliked" gorm:"default:0"`
	Views     uint  `json:"views" gorm:"default:0"`
	Comments  uint  `json:"comments" gorm:"default:0"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime"`
}

type Comment struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	UserID    uint   `json:"user_id" gorm:"not null"`
	ArticleID string `json:"article_id" gorm:"not null"`
	ParentID  uint   `json:"parent_id"` // ID of the parent comment, if any
	Content   string `json:"content" gorm:"not null"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
}

// Initializes the database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {

	dbPath := cfg.DB_Path

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	if _, err := os.Stat(dbDir); err != nil {
		return nil, fmt.Errorf("database directory not accessible: %w", err)
	}

	fmt.Printf("Attempting to open database at: %s\n", dbPath)
	fmt.Printf("Database directory: %s\n", dbDir)

	// Add grom config here
	config := &gorm.Config{}

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移数据库结构
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// 自动迁移数据库表结构
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&ArticleMeta{},
		&ArticleStatistics{},
		&Comment{},
	)
}

// Close Database connection
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
