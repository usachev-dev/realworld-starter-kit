package models

import (
	"../DB"
)

type User struct {
	ID           uint    `gorm:"primary_key" json:"id"`
	Username     string  `gorm:"column:username" json:"username"`
	Email        string  `gorm:"column:email;unique_index" json:"email"`
	Bio          string  `gorm:"column:bio;size:1024" json:"bio"`
	Image        *string `gorm:"column:image" json:"image"`
	PasswordHash string  `gorm:"column:password;not null"`
}

type Follow struct {
	ID           uint `gorm:"primary_key"`
	FollowingID  uint
	FollowedByID uint
}

func (u *User) Save() error {
	db := DB.Get()
	return db.Save(u).Error
}

func GetUser(email string) (*User, error) {
	db := DB.Get()
	var user User
	err := db.Where(&User{Email: email}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	db := DB.Get()
	var user User
	err := db.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsFollowing(followedByID uint, followingID uint) bool {
	db := DB.Get()
	var follow Follow
	followErr := db.Where(&Follow{FollowedByID: followedByID, FollowingID: followingID}).First(&follow).Error
	if followErr == nil {
		return true
	}
	return false
}
