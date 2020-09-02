package domain

import (
	"../api_errors"
	"../auth"
	"../models"
	"net/http"
	"net/url"
	"strings"
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

	tags, tagErr := models.GetTagsForArticle(article.ID)
	if tagErr != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("tags", tagErr.Error())
	}

	authorProfile, authorErr := GetProfile(article.Author.Username, tokenString)
	if authorErr != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("author", authorErr.Error())
	}

	return &ArticleResponse{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        *tagsToTagList(*tags),
		CreatedAt:      article.CreatedAt.String(),
		UpdatedAt:      article.UpdatedAt.String(),
		Favorited:      false,
		FavoritesCount: 0,
		Author:         *authorProfile,
	}, nil
}

func GetArticle(slug string, tokenString string) (*ArticleResponse, *api_errors.E) {
	article, err := models.GetArticle(slug)
	if err != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("slug", err.Error())
	}

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
		CreatedAt:      article.CreatedAt.String(),
		UpdatedAt:      article.UpdatedAt.String(),
		Favorited:      favorited,
		FavoritesCount: models.GetFavoriteCount(article.ID),
		Author:         *authorProfile,
	}, nil
}

func FavoriteArticle(slug string, tokenString string) (*ArticleResponse, *api_errors.E) {
	email, _ := auth.GetEmailFromTokenString(tokenString)
	user, userErr := models.GetUser(email)
	if userErr != nil {
		return nil, api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")

	}
	article, err := models.GetArticle(slug)
	if err != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("slug", err.Error())
	}
	models.FavoriteArticle(article.ID, user.ID)
	return GetArticle(slug, tokenString)
}
