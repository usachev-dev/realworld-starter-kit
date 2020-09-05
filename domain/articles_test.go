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

func setupListArticles(t *testing.T) string /* token */ {
	initDb()
	createUser(t)
	userResponse, _ := domain.SignIn(userSignIn)
	tokenString := userResponse.Token
	domain.CreateArticle(domain.ArticleCreate{
		Title:       "t1",
		Description: "d1",
		Body:        "b1",
		TagList:     []string{"t0", "t1"},
	}, tokenString)

	domain.CreateArticle(domain.ArticleCreate{
		Title:       "t2",
		Description: "d2",
		Body:        "b2",
		TagList:     []string{"t1"},
	}, tokenString)

	domain.CreateArticle(domain.ArticleCreate{
		Title:       "t3",
		Description: "d3",
		Body:        "b3",
		TagList:     []string{"t2"},
	}, tokenString)

	domain.FavoriteArticle("t2", tokenString)
	domain.FollowUser(userResponse.Username, tokenString)

	return tokenString
}

func tearDownListArticles() {
	DB.Get().Exec("DELETE from articles")
	destroyUser()
}

func TestListArticlesByTag(t *testing.T) {
	token := setupListArticles(t)
	defer tearDownListArticles()

	tag := "t1"
	result, count, err := domain.ListArticles(&tag, nil, nil, 0, 0, token)

	if err != nil {
		t.Fatalf("could not list articles: %s", err)
	}

	if len(*result) != int(count) {
		t.Fatalf("wrong article count, expected %d, got %d", len(*result), count)
	}

	if len(*result) != 2 {
		t.Fatalf("expected 2 articles with tag t1, got %d", len(*result))
	}

	if len((*result)[0].TagList) != 2 {
		t.Fatalf("expected 2 tags for first element, got %d", len((*result)[0].TagList))
	}
}

func TestListArticlesByFav(t *testing.T) {
	token := setupListArticles(t)
	defer tearDownListArticles()

	userName := userCreate.Username
	result, count, err := domain.ListArticles(nil, nil, &userName, 0, 0, token)

	if err != nil {
		t.Fatalf("could not list articles: %s", err)
	}

	if len(*result) != int(count) {
		t.Fatalf("wrong article count, expected %d, got %d", len(*result), count)
	}

	if len(*result) != 1 {
		t.Fatalf("expected 1 article favored by %s, got %d", userName, len(*result))
	}
}

func TestListArticlesByTagAndFav(t *testing.T) {
	token := setupListArticles(t)
	defer tearDownListArticles()

	userName := userCreate.Username
	tag := "t3"
	result, count, err := domain.ListArticles(&tag, nil, &userName, 0, 0, token)

	if err != nil {
		t.Fatalf("could not list articles: %s", err)
	}

	if len(*result) != int(count) {
		t.Fatalf("wrong article count, expected %d, got %d", len(*result), count)
	}

	if len(*result) != 0 {
		t.Fatalf("expected 0 article favored by %s with tag %s, got %d", userName, tag, len(*result))
	}
}

func TestListArticlesByAuthor(t *testing.T) {
	token := setupListArticles(t)
	defer tearDownListArticles()

	userName := userCreate.Username
	result, count, err := domain.ListArticles(nil, &userName, nil, 0, 0, token)

	if err != nil {
		t.Fatalf("could not list articles: %s", err)
	}

	if len(*result) != int(count) {
		t.Fatalf("wrong article count, expected %d, got %d", len(*result), count)
	}

	if len(*result) != 3 {
		t.Fatalf("expected 3 by author %s, got %d", userName, len(*result))
	}
}

func TestListAllArticles(t *testing.T) {
	token := setupListArticles(t)
	defer tearDownListArticles()

	result, count, err := domain.ListArticles(nil, nil, nil, 0, 0, token)

	if err != nil {
		t.Fatalf("could not list articles: %s", err)
	}

	if len(*result) != int(count) {
		t.Fatalf("wrong article count, expected %d, got %d", len(*result), count)
	}

	if len(*result) != 3 {
		t.Fatalf("expected 3 articles total, got %d", len(*result))
	}
}

func TestFeedArticles(t *testing.T) {
	token := setupListArticles(t)
	defer tearDownListArticles()

	result, count, err := domain.FeedArticles(0, 0, token)

	if err != nil {
		t.Fatalf("could not feed articles: %s", err)
	}

	if len(*result) != int(count) {
		t.Fatalf("wrong article count, expected %d, got %d", len(*result), count)
	}

	if len(*result) != 3 {
		t.Fatalf("expected 3 articles total, got %d", len(*result))
	}
}

func TestGetAllTags(t *testing.T) {
	setupListArticles(t)
	defer tearDownListArticles()

	result, err := domain.GetAllTags()

	if err != nil {
		t.Fatalf("could not get all tags: %s", err)
	}

	if len(*result) != 4 {
		t.Fatalf("expected 3 tags total, got %d: %+v", len(*result), *result)
	}
}

func TestCreateComment(t *testing.T) {
	token := setupListArticles(t)
	defer tearDownListArticles()

	result, err := domain.CreateComment("Hello comment", "t2", token)

	if err != nil {
		t.Fatalf("could not get create comment: %s", err)
	}

	if result.Body != "Hello comment" {
		t.Fatalf("comment body, expected \"Hello comment\", got %s", result.Body)
	}
}

func TestGetAllComments(t *testing.T) {
	token := setupListArticles(t)
	defer tearDownListArticles()

	domain.CreateComment("Hello comment", "t2", token)
	domain.CreateComment("Hello comment 2", "t2", token)
	result, err := domain.GetCommentsForArticle("t2", token)

	if err != nil {
		t.Fatalf("could not get get articles: %s", err)
	}

	if len(*result) == 0 {
		t.Fatalf("got 0 comments, expected at least 1")
	}
}
