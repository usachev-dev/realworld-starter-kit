package models

import (
	"../DB"
)

func AutoMigrate() {
	db := DB.Get()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Follow{})
	db.AutoMigrate(&Article{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Favorite{})
}
