package repo

import (
	"fmt"
	"log"

	// "sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"erp-backend/internal/models"
)

var DB *gorm.DB

// var dbOnce sync.Once
// var initErr error

// InitDatabase 初始化数据库连接
func InitDatabase(dbPath string) error {
	var err error

	// 配置GORM日志
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connected successfully")

	// 自动迁移所有模型
	err = AutoMigrate()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}

// AutoMigrate 自动迁移所有数据表
func AutoMigrate() error {
	return DB.AutoMigrate(
		// 核心实体
		&models.User{},
		&models.Location{},
		&models.Project{},
		&models.Donor{},
		&models.Volunteer{},
		&models.Employee{},

		// 财务相关
		&models.Transaction{},
		&models.Donation{},
		&models.Fund{},
		&models.Expense{},
		&models.Purchase{},
		&models.Payroll{},

		// 库存和礼品
		&models.Inventory{},
		&models.GiftType{},
		&models.Gift{},
		&models.InventoryTransaction{},
		&models.Delivery{},

		// 关联表
		&models.VolunteerProject{},
		&models.EmployeeProject{},
		&models.FundProject{},
		&models.DonationInventory{},
		&models.DeliveryInventory{},

		// Schedule
		&models.Schedule{},
	)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
