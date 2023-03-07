package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
}

func AddUser(tx *gorm.DB, id int64) error {
	return tx.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&User{Model: gorm.Model{ID: uint(id)}}).Error
}
