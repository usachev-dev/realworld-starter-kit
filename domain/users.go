package domain

import (
	"../api_errors"
	"../auth"
	"../models"
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

type Profile struct {
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

func CreateUser(u UserCreate) (UserResponse, *api_errors.E) {
	user := models.User{
		Username:     u.Username,
		PasswordHash: auth.PasswordToHash(u.Password),
		Email:        u.Email,
	}

	err := user.Save()
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
	}, nil
}

func SignIn(u UserSignIn) (UserResponse, *api_errors.E) {
	user, err := models.GetUser(u.Email)

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
	}, nil
}

func GetUser(token string) (UserResponse, *api_errors.E) {
	email, emailErr := auth.GetEmailFromTokenString(token)
	if emailErr != nil {
		return UserResponse{}, api_errors.NewError(http.StatusUnauthorized).Add("token", "token is invalid")

	}
	user, err := models.GetUser(email)
	if err != nil {
		return UserResponse{}, api_errors.NewError(http.StatusNotFound).Add("email", fmt.Sprintf("user not found"))
	}
	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
		Token:    token,
	}, nil
}

func UpdateUser(userUpdate UserUpdate, token string) (*UserResponse, *api_errors.E) {
	email, emailErr := auth.GetEmailFromTokenString(token)
	if emailErr != nil {
		return nil, api_errors.NewError(http.StatusUnauthorized).Add("token", "token is invalid")
	}
	user, userErr := models.GetUser(email)
	if userErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("email", fmt.Sprintf("could not find user with email %s", *userUpdate.Email))
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
	saveErr := user.Save()
	if saveErr != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("user", saveErr.Error())
	}

	tokenString := auth.GetTokenString(user.Email)
	return &UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
		Token:    tokenString,
	}, nil
}

// If request is not authenticated with token, pass empty string as token and profile's follow will be false
func GetProfile(username string, token string) (*Profile, *api_errors.E) {
	user, userErr := models.GetUserByUsername(username)

	if userErr != nil {
		return nil, api_errors.NewError(404).Add("username", fmt.Sprintf("could not find user with this username: %s", username))
	}

	following := false
	if token != "" {
		email, emailErr := auth.GetEmailFromTokenString(token)
		if emailErr == nil {
			follower, followerErr := models.GetUser(email)
			if followerErr == nil {
				following = models.IsFollowing(follower.ID, user.ID)
			}
		}
	}

	return &Profile{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: following,
	}, nil
}
