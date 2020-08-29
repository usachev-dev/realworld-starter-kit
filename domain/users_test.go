package domain_test

import (
	"../DB"
	"../auth"
	"../domain"
	"../models"
	"../utils"
	"fmt"
	"log"
	"testing"
)

func initDb() {
	err := DB.InitPostgres(utils.DbHost(), utils.DbPort(), utils.DbUser(), utils.DbName(), utils.DbPassword())
	if err != nil {
		panic(fmt.Sprintf("could not connect to test db: %s", err))
	}
	models.AutoMigrate()
}

func closeDb() {
	log.Printf("%s", DB.Close())
}

var userCreate domain.UserCreate = domain.UserCreate{
	Email:    "u@u",
	Password: "fretewrts",
	Username: "54tersdfg",
}

var userSignIn domain.UserSignIn = domain.UserSignIn{
	Email:    userCreate.Email,
	Password: userCreate.Password,
}

func createUser(t *testing.T) {
	_, err := domain.CreateUser(userCreate)
	if err != nil {
		t.Fatalf("could not create user: %s", err)
	}
}

func destroyUser() {
	DB.Get().Exec(fmt.Sprintf("DELETE FROM users WHERE email = '%s'", userCreate.Email))
}

func TestCreateUser(t *testing.T) {
	initDb()
	defer closeDb()
	r, err := domain.CreateUser(userCreate)
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
	defer closeDb()
	createUser(t)
	defer destroyUser()

	result, err := domain.SignIn(userSignIn)
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
	defer closeDb()
	createUser(t)
	defer destroyUser()

	userResponse, err := domain.SignIn(userSignIn)
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
	defer closeDb()
	createUser(t)
	defer destroyUser()
	userResponse, _ := domain.SignIn(userSignIn)
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
	defer closeDb()
	createUser(t)
	defer destroyUser()
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token
	userResponse, err := domain.GetUser(tokenString)
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
	defer closeDb()
	createUser(t)
	defer destroyUser()
	si, _ := domain.SignIn(userSignIn)
	tokenString := si.Token

	userName := "asdasd"
	userUpdate := domain.UserUpdate{
		Username: &userName,
	}

	userResponse, err := domain.UpdateUser(userUpdate, tokenString)
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
	defer closeDb()
	createUser(t)
	defer destroyUser()
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token
	profile, err := domain.GetProfile(userCreate.Username, tokenString)
	if err != nil {
		t.Fatalf("could not get profile: %s", err.Error())
	}
	if profile.Following {
		t.Fatalf("profile should not be following himself after signup")
	}
	if profile.Username != userCreate.Username {
		t.Fatalf("profile username %s is not the same as user's %s", profile.Username, userCreate.Username)
	}
}

func TestFollowUser(t *testing.T) {
	initDb()
	defer closeDb()
	createUser(t)
	defer destroyUser()
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token
	profile, err := domain.FollowUser(userCreate.Username, tokenString)
	if err != nil {
		t.Fatalf("could not follow user: %s", err.Error())
	}
	if !profile.Following {
		t.Fatalf("profile should be following after follow")
	}
}

func TestUnfollowUser(t *testing.T) {
	initDb()
	defer closeDb()
	createUser(t)
	defer destroyUser()
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token
	profile, err := domain.UnfollowUser(userCreate.Username, tokenString)
	if err != nil {
		t.Fatalf("could not unfollow user: %s", err.Error())
	}
	if profile.Following {
		t.Fatalf("profile should not be following after follow")
	}
}
