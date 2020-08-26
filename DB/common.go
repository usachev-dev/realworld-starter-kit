package DB

import "github.com/jinzhu/gorm"

var DB *gorm.DB

func Close() error {
	return DB.Close()
}

func Get() *gorm.DB {
	return DB
}
