package domain_test

import (
	"../DB"
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
		t.Fatalf("could not get article: %s", err.Error())
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
		t.Fatalf("could not favorite article: %s", err.Error())
	}

	if !result.Favorited || result.FavoritesCount == 0 {
		t.Fatalf("could not favorite an article")
	}

	unfavor, err2 := domain.UnfavoriteArticle(domain.SlugFromTitle(articleCreate.Title), tokenString)
	if err2 != nil {
		t.Fatalf("could not unfavorite article: %s", err.Error())
	}

	if unfavor.Favorited || unfavor.FavoritesCount != 0 {
		t.Fatalf("could not unfavorite an article")
	}
}

func TestDeleteArticle(t *testing.T) {
	initDb()
	defer closeDb()
	createArticle(t)
	defer destroyArticle()
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token

	_, aErr := domain.GetArticle(domain.SlugFromTitle(articleCreate.Title), "")
	if aErr != nil {
		t.Fatalf("could not get article: %s", aErr.Error())
	}

	err := domain.DeleteArticle(domain.SlugFromTitle(articleCreate.Title), tokenString)
	if err != nil {
		t.Fatalf("could not delete an article")
	}

	result, rErr := domain.GetArticle(domain.SlugFromTitle(articleCreate.Title), "")
	if result != nil || rErr == nil {
		t.Fatalf("deleted article returned")
	}
}

func TestUpdateArticle(t *testing.T) {
	initDb()
	defer closeDb()
	createArticle(t)
	defer destroyArticle()
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token
	update := map[string]interface{}{
		"title":       "newtitle",
		"tagList":     []string{"newtag"},
		"body":        "",
		"description": "",
	}

	result, err := domain.UpdateArticle(domain.SlugFromTitle(articleCreate.Title), update, tokenString)
	defer func() {
		DB.Get().Exec("DELETE FROM articles WHERE title = 'newtitle'")
	}()
	if err != nil {
		t.Fatalf("could not update article: %s", err.Error())
	}
	if result.Title != update["title"].(string) {
		t.Fatalf("article title did not update: %s, %s", result.Title, update["title"].(string))
	}
	if len(result.TagList) != 1 || result.TagList[0] != "newtag" {
		t.Fatalf("taglist did not properly update, %+v", result.TagList)
	}
}
