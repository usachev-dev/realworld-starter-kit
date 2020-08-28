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
	Email:    userCreate.Email,
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
	if err != nil {
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
	if err != nil {
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
	if err != nil {
		t.Fatalf("could not sign in: %s", err.Error())
	}
	result := auth.ValidateTokenStringWithEmail(userResponse.Token, userCreate.Email)
	if result != nil {
		t.Fatalf("%s; recieved invalid token: %s", result.Error(), userResponse.Token)
	}
}

func TestTokenHasEmail(t *testing.T) {
	initDb()
	defer DB.Close()
	createUser()
	defer destroyUser()
	userResponse, _ := models.SignIn(userSignIn)
	tokenString := userResponse.Token
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

func TestGetUser(t *testing.T) {
	initDb()
	defer DB.Close()
	createUser()
	defer destroyUser()
	userResponse, _ := models.SignIn(userSignIn)
	tokenString := userResponse.Token
	userResponse, err := models.GetUser(tokenString)
	if err != nil {
		t.Fatalf("could not get user: %s", err.Error())
	}
	if userResponse.Email != userCreate.Email {
		t.Fatalf("got user with wrong email")
	}
	if userResponse.Username != userCreate.Username {
		t.Fatalf("got user with wrong username")
	}
}

func TestUpdateUser(t *testing.T) {
	initDb()
	defer DB.Close()
	createUser()
	defer destroyUser()
	userResponse, _ := models.SignIn(userSignIn)
	tokenString := userResponse.Token

	userName := "asdasd"
	userUpdate := models.UserUpdate{
		Username: &userName,
	}

	userResponse, err := models.UpdateUser(userUpdate, tokenString)
	if err != nil {
		t.Fatalf("could not update user: %s", err.Error())
	}
	if userResponse.Username != userName {
		t.Fatalf("userName did not update")
	}
	if userResponse.Token == "" {
		t.Fatalf("user update response has no token")
	}
	validateError := auth.ValidateTokenString(userResponse.Token)
	if validateError != nil {
		t.Fatalf("user update response token is invalid: %s", validateError.Error())
	}

}

func TestGetProfile(t *testing.T) {
	initDb()
	defer DB.Close()
	createUser()
	defer destroyUser()
	userResponse, _ := models.SignIn(userSignIn)
	tokenString := userResponse.Token
	profile, err := models.GetProfile(userCreate.Username, tokenString)
	if err != nil {
		t.Fatalf("could not get profile")
	}
	if profile.Following {
		t.Fatalf("profile should not be following himself after signup")
	}
	if profile.Username != userCreate.Username {
		t.Fatalf("profile username %s is not the same as user's %s", profile.Username, userCreate.Username)
	}
}
