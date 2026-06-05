package database

import (
	"log"
	"os"
	"path/filepath"

	"chalaoshi-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接，如果数据库文件不存在则创建，并标记需要导入
func Init(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	needImport := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		needImport = true
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		// 启用 WAL 模式提升并发读取性能
	})
	if err != nil {
		return err
	}

	// 启用 WAL 模式 & 优化 SQLite 性能
	DB.Exec("PRAGMA journal_mode=WAL")
	DB.Exec("PRAGMA synchronous=NORMAL")
	DB.Exec("PRAGMA cache_size=-64000")
	DB.Exec("PRAGMA busy_timeout=5000")
	DB.Exec("PRAGMA temp_store=MEMORY")
	DB.Exec("PRAGMA mmap_size=268435456")

	// 自动迁移
	if err := DB.AutoMigrate(&models.Teacher{}, &models.Comment{}, &models.Course{}); err != nil {
		return err
	}

	if needImport {
		log.Println("[DB] 新数据库，需导入数据")
	}

	// SQLite LIKE 查询配合 B-tree 索引对10k数据已足够快

	return nil
}

// createFTSIndex 创建额外的搜索优化索引（使用LIKE查询已经足够快，不需要FTS5）

// NeedsImport 检查是否需要数据导入
func NeedsImport() bool {
	var count int64
	DB.Model(&models.Teacher{}).Count(&count)
	return count == 0
}
