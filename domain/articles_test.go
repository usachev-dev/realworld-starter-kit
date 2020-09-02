package domain_test

import (
	"../domain"
	"testing"
)

func TestCreateArticle(t *testing.T) {
	initDb()
	defer closeDb()
	createUser(t)
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token

	result, err := domain.CreateArticle(articleCreate, tokenString)
	defer destroyArticle()

	if err != nil {
		t.Fatalf("could not save article: %s", err.Error())
	}
	if result.Author.Username == "" {
		t.Fatalf("could not load author of saved article")
	}
}

//
//func TestGetArticle(t *testing.T) {
//	initDb()
//	defer closeDb()
//	createArticle(t)
//	defer destroyArticle()
//
//	result, err := domain.GetArticle(domain.SlugFromTitle(articleCreate.Title))
//	if err != nil {
//		t.Fatalf("could not get article: %s", err)
//	}
//	if  result.AuthorID == 0 || result.Author.ID == 0 {
//		t.Fatalf("could not load author of saved article")
//	}
//	if len(result.) == 0 {
//		t.Fatalf("could not load article tags")
//	}
//}
