package domain

import (
	"../api_errors"
	"../auth"
	"../models"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ArticleCreate struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type ArticleResponse struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount uint     `json:"favoritesCount"`
	Author         Profile  `json:"author"`
}

func formatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.999Z")
}

func SlugFromTitle(title string) string {
	return url.QueryEscape(strings.ToLower(strings.Join(strings.Split(title, " "), "-")))
}

func tagsFromTagList(list []string) *[]models.Tag {
	var result []models.Tag
	for _, tag := range list {
		result = append(result, models.Tag{Name: tag})
	}
	return &result
}

func tagsToTagList(tags []models.Tag) *[]string {
	var result []string
	for _, tag := range tags {
		result = append(result, tag.Name)
	}
	return &result
}

func CreateArticle(articleCreate ArticleCreate, tokenString string) (*ArticleResponse, *api_errors.E) {
	email, emErr := auth.GetEmailFromTokenString(tokenString)
	if emErr != nil {
		return nil, api_errors.NewError(http.StatusUnauthorized).Add("token", "invalid token")
	}

	user, uErr := models.GetUser(email)
	if uErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("user", "user not found")
	}

	article, err := models.CreateArticle(&models.Article{
		Title:       articleCreate.Title,
		Body:        articleCreate.Body,
		Description: articleCreate.Description,
		AuthorID:    user.ID,
		Slug:        SlugFromTitle(articleCreate.Title),
	}, articleCreate.TagList)
	if err != nil {
		return nil, api_errors.NewError(http.StatusUnprocessableEntity).Add("article", err.Error())
	}

	return articleToResponse(article, tokenString)
}

func articleToResponse(article *models.Article, tokenString string) (*ArticleResponse, *api_errors.E) {
	tags, tagErr := models.GetTagsForArticle(article.ID)
	if tagErr != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("tagList", tagErr.Error())
	}

	authorProfile, authorErr := GetProfile(article.Author.Username, tokenString)
	if authorErr != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("author", authorErr.Error())
	}

	var favorited bool = false
	if tokenString != "" {
		email, _ := auth.GetEmailFromTokenString(tokenString)
		user, _ := models.GetUser(email)
		if user != nil {
			favorited = models.IsArticleFavorited(article.ID, user.ID)
		}
	}

	return &ArticleResponse{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        *tagsToTagList(*tags),
		CreatedAt:      formatTime(article.CreatedAt),
		UpdatedAt:      formatTime(article.UpdatedAt),
		Favorited:      favorited,
		FavoritesCount: models.GetFavoriteCount(article.ID),
		Author:         *authorProfile,
	}, nil
}

func GetArticle(slug string, tokenString string) (*ArticleResponse, *api_errors.E) {
	article, err := models.GetArticle(slug)
	if err != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("slug", err.Error())
	}
	return articleToResponse(article, tokenString)
}

func FavoriteArticle(slug string, tokenString string) (*ArticleResponse, *api_errors.E) {
	email, _ := auth.GetEmailFromTokenString(tokenString)
	user, userErr := models.GetUser(email)
	if userErr != nil {
		return nil, api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")
	}

	article, articleErr := models.GetArticle(slug)
	if articleErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("slug", articleErr.Error())
	}

	err := models.FavoriteArticle(article.ID, user.ID)
	if err != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("article", articleErr.Error())
	}

	return GetArticle(slug, tokenString)
}

func UnfavoriteArticle(slug string, tokenString string) (*ArticleResponse, *api_errors.E) {
	email, _ := auth.GetEmailFromTokenString(tokenString)
	user, userErr := models.GetUser(email)
	if userErr != nil {
		return nil, api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")
	}

	article, articleErr := models.GetArticle(slug)
	if articleErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("slug", articleErr.Error())
	}

	err := models.UnFavoriteArticle(article.ID, user.ID)
	if err != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("article", articleErr.Error())
	}

	return GetArticle(slug, tokenString)
}

func DeleteArticle(slug string, tokenString string) *api_errors.E {
	email, _ := auth.GetEmailFromTokenString(tokenString)
	user, userErr := models.GetUser(email)
	if userErr != nil {
		return api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")
	}

	article, articleErr := models.GetArticle(slug)
	if articleErr != nil {
		return api_errors.NewError(http.StatusNotFound).Add("slug", articleErr.Error())
	}

	if article.AuthorID != user.ID {
		return api_errors.NewError(http.StatusForbidden).Add("token", "Cannot delete articles of other users")
	}

	err := models.DeleteArticle(article.ID)
	if err != nil {
		return api_errors.NewError(http.StatusInternalServerError).Add("article", err.Error())
	}
	return nil
}

func isString(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case string:
		return true
	default:
		return false
	}
}

func isStrSlice(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case []string:
		return true
	default:
		return false
	}
}

func UpdateArticle(slug string, updateData map[string]interface{}, tokenString string) (*ArticleResponse, *api_errors.E) {
	email, emErr := auth.GetEmailFromTokenString(tokenString)
	if emErr != nil {
		return nil, api_errors.NewError(http.StatusUnauthorized).Add("token", "invalid token")
	}

	user, uErr := models.GetUser(email)
	if uErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("user", "user not found")
	}

	article, articleErr := models.GetArticle(slug)
	if articleErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("slug", articleErr.Error())
	}

	if article.AuthorID != user.ID {
		return nil, api_errors.NewError(http.StatusForbidden).Add("token", "Cannot update articles of other users")
	}

	update := ArticleCreate{}
	if isString(updateData["title"]) {
		update.Title = updateData["title"].(string)
	}
	if isString(updateData["body"]) {
		update.Body = updateData["body"].(string)
	}
	if isString(updateData["description"]) {
		update.Description = updateData["description"].(string)
	}
	if isStrSlice(updateData["tagList"]) {
		update.TagList = updateData["tagList"].([]string)
	}

	result, err := models.UpdateArticle(slug, &models.Article{
		Title:       update.Title,
		Body:        update.Body,
		Description: update.Description,
		AuthorID:    user.ID,
		Slug:        SlugFromTitle(update.Title),
	}, update.TagList)
	if err != nil {
		return nil, api_errors.NewError(http.StatusUnprocessableEntity).Add("article", err.Error())
	}

	return articleToResponse(result, tokenString)
}
