package models_test

import (
	"../DB"
	"../api_errors"
	"../auth"
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

var userCreate models.UserCreate = models.UserCreate{
	Email:    "u@u",
	Password: "fretewrts",
	Username: "54tersdfg",
}

var userSignIn models.UserSignIn = models.UserSignIn{
	Email: userCreate.Email,
	Password: userCreate.Password,
}

func createUser() (models.UserResponse, *api_errors.E) {
	return models.CreateUser(userCreate)
}

func destroyUser() {
	DB.Get().Exec(fmt.Sprintf("DELETE FROM users WHERE email = '%s'", userCreate.Email))
}

func TestCreateUser(t *testing.T) {
	initDb()
	defer DB.Close()
	r, err := createUser()
	defer destroyUser()
	if !err.IsOk() {
		t.Fatalf("could not create user: %s", err)
	}
	if r.Token == "" {
		t.Fatalf("create user response has no token")
	}
}

func TestSignIn(t *testing.T) {
	initDb()
	defer DB.Close()
	createUser()
	defer destroyUser()

	result, err := models.SignIn(userSignIn)
	if !err.IsOk() {
		t.Fatalf("could not sign in, %s", err.Error())
	}
	if result.Token == "" {
		t.Fatalf("sign in response has no token")
	}
	if result.Email != userCreate.Email {
		t.Fatalf("emails do not match: %s %s", result.Email, userCreate.Email)
	}
	if result.Username != userCreate.Username {
		t.Fatalf("usernames do not match: %s %s", result.Username, userCreate.Username)
	}
}

func TestSignedInUserHasValidToken(t *testing.T) {
	initDb()
	defer DB.Close()
	createUser()
	defer destroyUser()

	userResponse, err := models.SignIn(userSignIn)
	if err != nil && !err.IsOk() {
		t.Fatalf("could not sign in: %s", err.Error())
	}
	result := auth.ValidateTokenString(userResponse.Token, userCreate.Email)
	if !result.IsOk() {
		t.Fatalf("%s; recieved invalid token: %s", result.Error(), userResponse.Token)
	}
}

func TestTokenHasEmail(t *testing.T) {
	initDb()
	defer DB.Close()
	createUser()
	defer destroyUser()
	userResponse, _ := models.SignIn(userSignIn)
	tokenString :=  userResponse.Token
	email, err := auth.GetEmailFromTokenString(tokenString)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if email == "" {
		t.Fatalf("recieved empty email from token")
	}
	if email != userCreate.Email {
		t.Fatalf("token contained wrong email")
	}
}
