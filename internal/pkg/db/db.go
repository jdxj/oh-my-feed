package db

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jdxj/oh-my-feed/internal/pkg/config"
	"github.com/jdxj/oh-my-feed/internal/pkg/log"
)

var (
	db *gorm.DB
)

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB.User, config.DB.Password, config.DB.Address, config.DB.Port, config.DB.Dbname)
	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("open db err: %s", err)
	}
}

func WithContext(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func Debug() {
	db = db.Debug()
}
