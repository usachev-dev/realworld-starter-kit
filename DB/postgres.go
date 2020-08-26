package DB

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


func InitPostgres(
	host string,
	port string,
	user string,
	dbname string,
	password string,
) error {
	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		user,
		dbname,
		password,
	)

	result, err := connectPostgres(connectString)
	if err != nil {
		return err
	}
	DB = result
	return nil
}


func connectPostgres(connectString string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connectString)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to DB: %s", err)
	}
	return db, nil
}
