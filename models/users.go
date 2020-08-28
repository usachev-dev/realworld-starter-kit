package models

import (
	"../DB"
	"../api_errors"
	"../auth"
	"fmt"
	"net/http"
)

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           uint    `gorm:"primary_key" json:"id"`
	Username     string  `gorm:"column:username" json:"username"`
	Email        string  `gorm:"column:email;unique_index" json:"email"`
	Bio          string  `gorm:"column:bio;size:1024" json:"bio"`
	Image        *string `gorm:"column:image" json:"image"`
	PasswordHash string  `gorm:"column:password;not null"`
}

type UserResponse struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Bio      string  `json:"bio"`
	Image    *string `json:"image"`
	Token    string  `json:"token"`
}

type UserUpdate struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Image    *string `json:"image"`
	Bio      *string `json:"bio"`
}

func CreateUser(u UserCreate) (UserResponse, *api_errors.E) {
	user := User{
		Username:     u.Username,
		PasswordHash: auth.PasswordToHash(u.Password),
		Email:        u.Email,
	}

	db := DB.Get()
	err := db.Save(&user).Error
	if err != nil {
		return UserResponse{}, api_errors.NewError(http.StatusInternalServerError).Add("body", err.Error())
	}

	tokenString := auth.GetTokenString(u.Email)
	return UserResponse{
		Username: u.Username,
		Email:    u.Email,
		Bio:      "",
		Image:    nil,
		Token:    tokenString,
	}, &api_errors.Ok
}

func SignIn(u UserSignIn) (UserResponse, *api_errors.E) {
	db := DB.Get()
	var user User
	err := db.Where(&User{Email: u.Email}).First(&user).Error

	if err != nil || auth.CheckPassword(u.Password, user.PasswordHash) != nil {
		return UserResponse{},
			api_errors.NewError(http.StatusUnauthorized).Add("body", "email and password do not match")
	}

	tokenString := auth.GetTokenString(user.Email)

	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
		Token:    tokenString,
	}, &api_errors.Ok
}

func GetUser(token string) (UserResponse, *api_errors.E) {
	email, emailErr := auth.GetEmailFromTokenString(token)
	if emailErr != nil {
		return UserResponse{}, api_errors.NewError(http.StatusUnauthorized).Add("email", emailErr.Error())
	}
	db := DB.Get()
	var user User
	dbErr := db.Where(&User{Email: email}).First(&user).Error
	if dbErr != nil {
		return UserResponse{}, api_errors.NewError(http.StatusNotFound).Add("email", fmt.Sprintf("could not find user with email %s", email))
	}
	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
		Token:    token,
	}, &api_errors.Ok
}

func UpdateUser(userUpdate UserUpdate, token string) (UserResponse, *api_errors.E) {
	email, emailErr := auth.GetEmailFromTokenString(token)
	if emailErr != nil {
		return UserResponse{}, api_errors.NewError(http.StatusUnauthorized).Add("token", "token shpuld contain email info")
	}
	db := DB.Get()
	var user User
	dbErr := db.Where(&User{Email: email}).First(&user).Error
	if dbErr != nil {
		return UserResponse{}, api_errors.NewError(http.StatusNotFound).Add("email", fmt.Sprintf("could not find user with email %s", *userUpdate.Email))
	}
	if userUpdate.Email != nil {
		user.Email = *userUpdate.Email
	}
	if userUpdate.Username != nil {
		user.Username = *userUpdate.Username
	}
	if userUpdate.Bio != nil {
		user.Bio = *userUpdate.Bio
	}
	if userUpdate.Image != nil {
		user.Image = userUpdate.Image
	}
	if userUpdate.Password != nil {
		user.PasswordHash = auth.PasswordToHash(*userUpdate.Password)
	}
	saveErr := db.Save(&user).Error
	if saveErr != nil {
		return UserResponse{}, api_errors.NewError(http.StatusInternalServerError).Add("user", saveErr.Error())
	}

	tokenString := auth.GetTokenString(user.Email)
	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
		Token:    tokenString,
	}, &api_errors.Ok
}
