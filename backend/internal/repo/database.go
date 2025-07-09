package repo

import (
	"fmt"
	"mypage-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据模型
type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"default:user"`
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
}

type Comment struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	UserID    uint   `json:"user_id" gorm:"not null"`
	User      User   `json:"user" gorm:"foreignKey:UserID"`
	ArticleID string `json:"article_id" gorm:"not null"`
	Content   string `json:"content" gorm:"not null"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}



// 初始化数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 自动迁移
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Comment{})
}