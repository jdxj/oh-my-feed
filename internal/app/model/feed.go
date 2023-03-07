package model

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jdxj/oh-my-feed/internal/pkg/db"
)

type Feed struct {
	gorm.Model

	Address    string
	AddressMD5 string
	LatestPost string
}

func AddFeed(tx *gorm.DB, address string) (uint64, error) {
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

func UpdateLatestPost(ctx context.Context, id uint, url string) error {
	return db.WithContext(ctx).Model(Feed{}).
		Where("id = ?", id).
		Where("latest_post != ?", url).
		Update("latest_post", url).
		Error
}
