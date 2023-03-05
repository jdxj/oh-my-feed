package model

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/jdxj/oh-my-feed/validator"
)

type UserFeed struct {
	gorm.Model

	TelegramID int64
	FeedID     uint64
}

func AddUserFeed(ctx context.Context, telegramID int64, address string) error {
	address, err := validator.ValidateFeed(ctx, address)
	if err != nil {
		return err
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := AddUser(tx, telegramID)
		if err != nil {
			return err
		}

		feedID, err := AddFeed(tx, address)
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{DoNothing: true}).
			Create(&UserFeed{
				TelegramID: telegramID,
				FeedID:     feedID,
			}).Error
	})
}

type ListUserFeedReq struct {
	TelegramID int64
	FeedID     uint64

	Offset int
	Limit  int
}

type ListUserFeedRsp struct {
	Count     int64
	UserFeeds []UserFeed
}

func ListUserFeed(ctx context.Context, req ListUserFeedReq) (ListUserFeedRsp, error) {
	tx := db.WithContext(ctx).Model(UserFeed{})
	if v := req.FeedID; v != 0 {
		tx.Where("feed_id = ?", req.FeedID)
	}
	if v := req.TelegramID; v != 0 {
		tx.Where("telegram_id = ?", req.TelegramID)
	}

	rsp := ListUserFeedRsp{}
	err := tx.Count(&rsp.Count).Error
	if err != nil {
		return rsp, err
	}

	if v := req.Offset; v > 0 {
		tx.Offset(v)
	}
	if v := req.Limit; v > 0 {
		tx.Limit(v)
	}

	return rsp, tx.Find(&rsp.UserFeeds).Error
}
