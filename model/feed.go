package model

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Feed struct {
	gorm.Model

	Address    string
	AddressMD5 string
}

func AddFeed(tx *gorm.DB, address string) (uint64, error) {
	address = strings.TrimSuffix(address, "/")
	err := myValidator.Var(address, "url")
	if err != nil {
		return 0, err
	}

	sum := md5.Sum([]byte(address))
	sumStr := hex.EncodeToString(sum[:])

	feed := &Feed{
		Address:    address,
		AddressMD5: sumStr,
	}
	return uint64(feed.ID),
		// 更新updated_at以使feed.ID被赋值
		tx.Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]any{"updated_at": time.Now()})}).
			Create(feed).Error
}

func GetFeed(tx *gorm.DB, id uint64) (Feed, error) {
	feed := Feed{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	return feed, tx.Take(&feed).Error
}

func GetFeeds(ctx context.Context) ([]Feed, error) {
	var feeds []Feed
	return feeds, db.WithContext(ctx).
		Find(&feeds).Error
}

func UpdateFeedTitle(ctx context.Context, id uint, title string) error {
	return db.WithContext(ctx).Model(Feed{}).
		Where("id = ?", id).
		Where("title != ?", title).
		Update("title", title).
		Error
}
