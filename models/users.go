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

type Profile struct {
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

type Follow struct {
	ID           uint `gorm:"primary_key"`
	FollowingID  uint
	FollowedByID uint
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
	}, nil
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
	}, nil
}

func getUser(token string) (User, *api_errors.E) {
	email, emailErr := auth.GetEmailFromTokenString(token)
	if emailErr != nil {
		return User{}, api_errors.NewError(http.StatusUnauthorized).Add("email", emailErr.Error())
	}
	db := DB.Get()
	var user User
	dbErr := db.Where(&User{Email: email}).First(&user).Error
	if dbErr != nil {
		return User{}, api_errors.NewError(http.StatusNotFound).Add("email", fmt.Sprintf("could not find user with email %s", email))
	}
	return user, nil
}

func GetUser(token string) (UserResponse, *api_errors.E) {
	user, err := getUser(token)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
		Token:    token,
	}, nil
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
	}, nil
}

func isFollowing(followedByID uint, followingID uint) bool {
	db := DB.Get()
	var follow Follow
	followErr := db.Where(&Follow{FollowedByID: followedByID, FollowingID: followingID}).First(&follow).Error
	if followErr == nil {
		return true
	}
	return false
}

// If request is not authenticated with token, pass empty string as token and profile's follow will be false
func GetProfile(username string, token string) (Profile, *api_errors.E) {
	db := DB.Get()

	// TODO parallel
	var user User
	dbErr := db.Where(&User{Username: username}).First(&user).Error
	if dbErr != nil {
		return Profile{}, api_errors.NewError(404).Add("username", fmt.Sprintf("could not find user with this username: %s", username))
	}

	following := false
	if token != "" {
		follower, followerErr := getUser(token)
		if followerErr == nil {
			following = isFollowing(follower.ID, user.ID)
		}
	}

	return Profile{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: following,
	}, nil
}
