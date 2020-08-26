package models

import (
	"../DB"
)

func AutoMigrate() {
	db := DB.Get()
	db.AutoMigrate(&User{})
}
