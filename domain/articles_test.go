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
	if len(result.TagList) == 0 {
		t.Fatalf("could not load article tags")
	}
}

func TestGetArticle(t *testing.T) {
	initDb()
	defer closeDb()
	createArticle(t)
	defer destroyArticle()

	result, err := domain.GetArticle(domain.SlugFromTitle(articleCreate.Title), "")
	if err != nil {
		t.Fatalf("could not get article: %s", err)
	}

	if len(result.TagList) == 0 {
		t.Fatalf("could not load article tags")
	}
}

func TestFavoriteArticle(t *testing.T) {
	initDb()
	defer closeDb()
	createArticle(t)
	defer destroyArticle()
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token

	result, err := domain.FavoriteArticle(domain.SlugFromTitle(articleCreate.Title), tokenString)
	if err != nil {
		t.Fatalf("could not favorite article: %s", err)
	}

	if !result.Favorited || result.FavoritesCount == 0 {
		t.Fatalf("could not favorite an article")
	}

	unfavor, err2 := domain.UnfavoriteArticle(domain.SlugFromTitle(articleCreate.Title), tokenString)
	if err2 != nil {
		t.Fatalf("could not unfavorite article: %s", err)
	}

	if unfavor.Favorited || unfavor.FavoritesCount != 0 {
		t.Fatalf("could not unfavorite an article")
	}
}
