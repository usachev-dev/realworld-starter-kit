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
	FollowingID  uint `gorm:"primaryKey;autoIncrement:false"`
	FollowedByID uint `gorm:"primaryKey;autoIncrement:false"`
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
	err := db.Where(&Follow{FollowedByID: followedByID, FollowingID: followingID}).First(&follow).Error
	if err == nil {
		return true
	}
	return false
}

func AddFollow(followedByID uint, followingID uint) error {
	db := DB.Get()
	follow := Follow{
		FollowingID:  followingID,
		FollowedByID: followedByID,
	}
	err := db.Save(follow).Error
	if err != nil {
		// if error was because follow already exists, we report success
		foundErr := db.Where(&follow).First(&follow).Error
		if foundErr == nil {
			return nil
		}
		return err
	}
	return nil
}
