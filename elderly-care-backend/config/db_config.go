package config

import (
	"elderly-care-backend/global"
	"elderly-care-backend/models"

	//"elderly-care-backend/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.Database.User,
		Config.Database.Password,
		Config.Database.Host,
		Config.Database.Port,
		Config.Database.Db)

	var logLevel logger.LogLevel
	if Config.Server.Profile == "dev" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second * 5,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		fmt.Fprintf(logFile, "Cannot connect to database: %v\n", err)
		// 关键修复：连接失败时必须返回，不能继续执行
		log.Fatalf("数据库连接失败: %v", err)
		return
	}

	// 测试数据库连接是否真正可用
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Fprintf(logFile, "获取数据库实例失败: %v\n", err)
		log.Fatalf("获取数据库实例失败: %v", err)
		return
	}

	// 验证连接
	if err := sqlDB.Ping(); err != nil {
		fmt.Fprintf(logFile, "数据库连接测试失败: %v\n", err)
		log.Fatalf("数据库连接测试失败: %v", err)
		return
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(Config.Database.IdleConnections)
	sqlDB.SetMaxOpenConns(Config.Database.MaxConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(Config.Database.MaxLifetime) * time.Second)

	// 设置全局变量
	global.Db = db

	// 执行数据库迁移
	modelsList := []interface{}{
		&models.Account{},
		&models.Message{},
		&models.AccountEvaluation{},
		&models.ContactList{},
		&models.UserLocation{}, // 添加用户位置表
		//&models.Task{},
		//&models.SOSRecord{},
	}

	if err = db.AutoMigrate(modelsList...); err != nil {
		fmt.Fprintf(logFile, "数据库迁移失败: %v\n", err)
		log.Fatalf("数据库迁移失败: %v", err)
		return
	}

	fmt.Fprintf(logFile, "数据库初始化成功\n")
	log.Println("数据库初始化和迁移完成")
}
