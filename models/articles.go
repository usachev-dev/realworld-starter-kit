package models

import (
	"../DB"
	"fmt"
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

type ArticlesList struct {
	Article
	User
	Tag
	Favorite
}

func CreateArticle(a *Article, tags []string) (*Article, error) {
	db := DB.Get()
	err := db.Transaction(func(tx *gorm.DB) error {
		saveErr := tx.Omit("Author").Create(&a).Error
		if saveErr != nil {
			return saveErr
		}
		getErr := tx.Where(&a).Preload("Author").First(&a).Error
		if getErr != nil {
			return getErr
		}
		for _, tag := range tags {
			tagErr := tx.Create(&Tag{ArticleID: a.ID, Name: tag}).Error
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

func UnFavoriteArticle(articleID uint, userID uint) error {
	db := DB.Get()
	return db.Delete(&Favorite{ArticleID: articleID, UserID: userID}).Error
}

func DeleteArticle(articleID uint) error {
	db := DB.Get()
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&Article{}, articleID).Error
		if err != nil {
			return err
		}
		tagRmErr := tx.Where(&Tag{ArticleID: articleID}).Delete(&Tag{}).Error
		if tagRmErr != nil {
			return tagRmErr
		}
		favRmErr := tx.Where(&Favorite{ArticleID: articleID}).Delete(&Favorite{}).Error
		if favRmErr != nil {
			return favRmErr
		}
		return nil
	})
}

func UpdateArticle(slug string, a *Article, tags *[]string) (*Article, error) {
	db := DB.Get()
	var article Article
	getErr := db.Where(&Article{Slug: slug}).First(&article).Error
	if getErr != nil {
		return nil, getErr
	}
	// TODO remove mutation
	// If updated without ID, it tries to insert
	a.ID = article.ID
	err := db.Transaction(func(tx *gorm.DB) error {
		saveErr := tx.Save(a).First(&article).Error
		if saveErr != nil {
			return saveErr
		}
		if tags != nil {
			tagRmErr := tx.Where(&Tag{ArticleID: article.ID}).Delete(&Tag{}).Error
			if tagRmErr != nil {
				return tagRmErr
			}
			for _, tag := range *tags {
				tagErr := tx.Create(&Tag{ArticleID: article.ID, Name: tag}).Error
				if tagErr != nil {
					return tagErr
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func ListArticles(tag string, authorID uint, favoritedByID uint, limit uint, offset uint, userID uint) (*[]ArticlesList, error) {
	db := DB.Get()

	var result []ArticlesList

	query := "SELECT * " +
		fmt.Sprintf("FROM (SELECT * FROM articles LIMIT %d OFFSET %d) as articles ", limit, offset) +
		"LEFT JOIN tags on tags.article_id = articles.id " +
		"LEFT JOIN users on users.id = articles.author_id " +
		fmt.Sprintf("LEFT JOIN favorites on favorites.article_id = articles.id and favorites.user_id = %d ", userID)

	filter := []string{}

	if tag != "" {
		filter = append(filter, fmt.Sprintf("'%s' in (SELECT tags.name FROM articles as a JOIN tags on tags.article_id = articles.id WHERE articles.id = a.id)", tag))
	}

	if authorID != 0 {
		filter = append(filter, fmt.Sprintf("articles.author_id = %d", authorID))
	}

	if favoritedByID != 0 {
		filter = append(filter, fmt.Sprintf("favorites.user_id = %d", favoritedByID))
	}

	if len(filter) > 0 {
		query = query + "WHERE "
		for i, f := range filter {
			query = query + f + " "
			if i < len(filter)-1 {
				query = query + "and "
			}
		}
	}

	err := db.Raw(query).Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}
