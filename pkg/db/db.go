package db

import (
	"errors"
	"log"

	"homework03/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局DB实例
var DB *gorm.DB

type Config struct {
	Driver string
	DSN    string
}

func Init(cfg Config) error {
	log.Printf("传入配置 Driver=%q, DSN=%q", cfg.Driver, cfg.DSN)

	var dialector gorm.Dialector
	switch cfg.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	default:
		return errors.New("只支持 sqlite")
	}

	dbInstance, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("【错误】gorm.Open 打开数据库失败：%v", err)
		return err
	}

	// 只有上面Open无错误，才会执行这里！
	DB = dbInstance
	log.Printf("✅ 全局DB赋值成功，地址=%p", DB)

	err = DB.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	)
	if err != nil {
		log.Printf("【错误】AutoMigrate迁移表失败：%v", err)
		return err
	}
	log.Println("✅ 数据表迁移完成")
	return nil
}