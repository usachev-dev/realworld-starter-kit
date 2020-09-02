package domain_test

import (
	"../DB"
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
	Email:    "u22@u",
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
	var user models.User
	DB.Get().Where(&models.User{Email: userCreate.Email}).First(&user)
	DB.Get().Exec(fmt.Sprintf("DELETE FROM users WHERE email = '%s'", userCreate.Email))
	DB.Get().Exec(fmt.Sprintf("DELETE FROM follows WHERE following_id = '%d'", user.ID))
	DB.Get().Exec(fmt.Sprintf("DELETE FROM follows WHERE followed_by_id = '%d'", user.ID))

}

var articleCreate = domain.ArticleCreate{
	Title:       "titledsasdfsdf",
	Description: "Description",
	Body:        "Body",
	TagList:     []string{"tag1", "tag2"},
}

func createArticle(t *testing.T) {
	createUser(t)
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token
	_, err := domain.CreateArticle(articleCreate, tokenString)
	if err != nil {
		t.Fatalf("could not create article: %s", err)
	}
}

func destroyArticle() {
	var article models.Article
	DB.Get().Where(&models.Article{Title: articleCreate.Title}).First(&article)
	DB.Get().Exec(fmt.Sprintf("DELETE FROM articles WHERE title = '%s'", articleCreate.Title))
	DB.Get().Exec(fmt.Sprintf("DELETE FROM tags WHERE article_id = '%d'", article.ID))
	destroyUser()
}
