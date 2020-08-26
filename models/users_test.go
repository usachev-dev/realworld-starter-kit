package models_test

import (
	"../DB"
	"../models"
	"../utils"
	"fmt"
	"testing"
)

func initDb() {
	err := DB.InitPostgres(utils.DbHost(), utils.DbPort(), utils.DbUser(), utils.DbName(), utils.DbPassword())
	if err != nil {
		panic(fmt.Sprintf("could not connect to test db: %s", err))
	}
	models.AutoMigrate()
}

func TestCreateUser(t *testing.T) {
	initDb()
	defer DB.Close()
	u := models.UserCreate{
		Email:    "u@u",
		Password: "fretewrts",
		Username: "54tersdfg",
	}
	defer DB.Get().Exec(fmt.Sprintf("DELETE FROM users WHERE email = '%s'", u.Email))

	r, err := models.CreateUser(u)
	fmt.Printf("%+v", r)
	if !err.IsOk() {
		t.Fatalf("could not create user: %s", err)
	}
	if r.Token == "" {
		t.Fatalf("create user response has no token")
	}
}
