package models

import (
	"../DB"
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Slug        string `gorm:"unique_index"`
	Title       string
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
	AuthorID    uint
	Author      User `gorm:"foreignKey:AuthorID"`
}

type Tag struct {
	ArticleID uint
	Name      string
}

type Favorite struct {
	ArticleID uint
	UserID    uint
}

func CreateArticle(a *Article, tags []string) (*Article, error) {
	db := DB.Get()
	err := db.Transaction(func(tx *gorm.DB) error {
		saveErr := db.Omit("Author").Create(&a).Error
		if saveErr != nil {
			return saveErr
		}
		getErr := db.Where(&a).Preload("Author").First(&a).Error
		if getErr != nil {
			return getErr
		}
		for _, tag := range tags {
			tagErr := db.Create(&Tag{ArticleID: a.ID, Name: tag}).Error
			if tagErr != nil {
				return tagErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

func GetTagsForArticle(articleID uint) (*[]Tag, error) {
	db := DB.Get()
	var result []Tag
	err := db.Where(&Tag{ArticleID: articleID}).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetArticle(slug string) (*Article, error) {
	db := DB.Get()
	var a Article
	err := db.Where(&Article{
		Slug: slug,
	}).Preload("Author").First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func IsArticleFavorited(articleID uint, userID uint) bool {
	db := DB.Get()
	count := db.Where(&Favorite{ArticleID: articleID, UserID: userID}).Find(&[]Favorite{}).RowsAffected
	return count > 0
}

func GetFavoriteCount(articleID uint) uint {
	db := DB.Get()
	count := uint(db.Where(&Favorite{ArticleID: articleID}).Find(&[]Favorite{}).RowsAffected)
	return count
}

func FavoriteArticle(articleID uint, userID uint) error {
	db := DB.Get()
	return db.Save(&Favorite{ArticleID: articleID, UserID: userID}).Error
}

func DeleteArticle(articleID *uint) error {
	db := DB.Get()
	err := db.Delete(&Article{}, articleID).Error
	return err
}
