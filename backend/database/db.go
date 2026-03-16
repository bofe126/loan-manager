package database

import (
	"fmt"
	"loan-manager-wails/backend/models"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 打开数据库连接（使用纯 Go SQLite 驱动）
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 检查表是否存在
	if !DB.Migrator().HasTable(&models.Loan{}) {
		// 表不存在，创建新表
		if err := DB.AutoMigrate(&models.Loan{}, &models.Payment{}); err != nil {
			return fmt.Errorf("数据库迁移失败: %w", err)
		}
	} else {
		// 表已存在，只添加缺失的列，不修改现有列
		if err := DB.Migrator().AutoMigrate(&models.Loan{}, &models.Payment{}); err != nil {
			// 如果 AutoMigrate 失败，忽略错误继续运行
			fmt.Println("警告: 数据库迁移遇到问题，但将继续使用现有结构:", err)
		}
	}

	fmt.Println("数据库初始化成功:", dbPath)
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
