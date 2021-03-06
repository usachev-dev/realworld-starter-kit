package domain

import (
	"../api_errors"
	"../auth"
	"../models"
	"fmt"
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

type CommentResponse struct {
	ID        uint    `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Body      string  `json:"body"`
	Author    Profile `json:"author"`
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
	var slugUpdate string = slug
	if isString(updateData["title"]) {
		update.Title = updateData["title"].(string)
		slugUpdate = SlugFromTitle(update.Title)
	}
	if isString(updateData["body"]) {
		update.Body = updateData["body"].(string)
	}
	if isString(updateData["description"]) {
		update.Description = updateData["description"].(string)
	}
	var tagListUpdate *[]string
	var tagList []string
	if isStrSlice(updateData["tagList"]) {
		tagList = updateData["tagList"].([]string)
		tagListUpdate = &tagList
	}

	result, err := models.UpdateArticle(slug, &models.Article{
		Title:       update.Title,
		Body:        update.Body,
		Description: update.Description,
		AuthorID:    user.ID,
		Slug:        slugUpdate,
	}, tagListUpdate)
	if err != nil {
		return nil, api_errors.NewError(http.StatusUnprocessableEntity).Add("article", err.Error())
	}

	return articleToResponse(result, tokenString)
}

func articlesListToResponse(list []models.ArticlesList, userID uint) *[]ArticleResponse {
	var lastId uint
	var result []ArticleResponse
	for _, el := range list {
		if lastId != el.Article.ID {
			lastId = el.Article.ID
			result = append(result, ArticleResponse{
				Slug:           el.Slug,
				Title:          el.Title,
				Description:    el.Description,
				Body:           el.Body,
				TagList:        []string{},
				CreatedAt:      formatTime(el.CreatedAt),
				UpdatedAt:      formatTime(el.UpdatedAt),
				Favorited:      false,
				FavoritesCount: 0,
				Author: Profile{
					Username:  el.User.Username,
					Bio:       el.User.Bio,
					Image:     el.User.Image,
					Following: models.IsFollowing(userID, el.User.ID),
				},
			})
		}

		if len(result) > 0 && el.Tag.Name != "" {
			result[len(result)-1].TagList = append(result[len(result)-1].TagList, el.Tag.Name)
		}

		if len(result) > 0 && el.Favorite.UserID != 0 {
			result[len(result)-1].FavoritesCount++
		}

		if len(result) > 0 && el.Favorite.UserID != 0 && el.Favorite.UserID == userID {
			result[len(result)-1].Favorited = true
		}
	}

	return &result

}

func ListArticles(tag *string, authorUsername *string, favoriteByUsername *string, limit uint, offset uint, tokenString string) (*[]ArticleResponse, uint, *api_errors.E) {
	tagFilter := ""
	var authorID uint = 0
	var favoredById uint = 0
	var userID uint = 0

	if limit == 0 {
		limit = 20
	}

	if tag != nil {
		tagFilter = *tag
	}

	if authorUsername != nil {
		author, aErr := models.GetUserByUsername(*authorUsername)
		if aErr != nil {
			return nil, 0, api_errors.NewError(http.StatusNotFound).Add("author", fmt.Sprintf("author with username %s not found", *authorUsername))
		}
		authorID = author.ID
	}

	if favoriteByUsername != nil {
		favored, fErr := models.GetUserByUsername(*favoriteByUsername)
		if fErr != nil {
			return nil, 0, api_errors.NewError(http.StatusNotFound).Add("favorited", fmt.Sprintf("user with username %s not found", *favoriteByUsername))
		}
		favoredById = favored.ID
	}

	if tokenString != "" {
		email, _ := auth.GetEmailFromTokenString(tokenString)
		user, uErr := models.GetUser(email)
		if uErr == nil {
			userID = user.ID
		}
	}

	list, count, listErr := models.ListArticles(tagFilter, authorID, favoredById, limit, offset, userID)
	if listErr != nil {
		return nil, 0, api_errors.NewError(http.StatusInternalServerError).Add("articles", listErr.Error())
	}

	return articlesListToResponse(*list, userID), count, nil

}

func FeedArticles(limit uint, offset uint, tokenString string) (*[]ArticleResponse, uint, *api_errors.E) {
	if limit == 0 {
		limit = 20
	}

	email, _ := auth.GetEmailFromTokenString(tokenString)
	user, uErr := models.GetUser(email)
	if uErr != nil {
		return nil, 0, api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")
	}

	result, count, err := models.FeedArticles(limit, offset, user.ID)
	if err != nil {
		return nil, 0, api_errors.NewError(http.StatusInternalServerError).Add("articles", err.Error())
	}

	return articlesListToResponse(*result, user.ID), count, nil
}

func GetAllTags() (*[]string, *api_errors.E) {
	tags, err := models.GetAllTags()
	if err != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("tags", "could not get tags")
	}
	result := []string{}
	for _, t := range *tags {
		result = append(result, t.Name)
	}
	return &result, nil
}

func CreateComment(body string, articleSlug string, tokenString string) (*CommentResponse, *api_errors.E) {
	email, mailErr := auth.GetEmailFromTokenString(tokenString)
	if mailErr != nil {
		return nil, api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")
	}

	user, uErr := models.GetUser(email)
	if uErr != nil {
		return nil, api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")
	}

	profile, pErr := GetProfile(user.Username, tokenString)
	if pErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("author", "author not found")
	}

	article, aErr := models.GetArticle(articleSlug)
	if aErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("article", "article not found")
	}

	result, err := models.CreateComment(user.ID, article.ID, body)
	if err != nil {
		return nil, api_errors.NewError(http.StatusInternalServerError).Add("comment", err.Error())
	}

	return &CommentResponse{
		ID:        result.ID,
		CreatedAt: formatTime(result.CreatedAt),
		UpdatedAt: formatTime(result.UpdatedAt),
		Body:      result.Body,
		Author:    *profile,
	}, nil
}

func GetCommentsForArticle(slug string, tokenString string) (*[]CommentResponse, *api_errors.E) {
	article, aErr := models.GetArticle(slug)
	if aErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("article", "article not found")
	}

	comments, cErr := models.GetCommentsForArticle(article.ID)
	if cErr != nil {
		return nil, api_errors.NewError(http.StatusNotFound).Add("article", "article not found")
	}

	result := []CommentResponse{}
	for _, c := range *comments {
		profile, _ := GetProfile(c.User.Username, tokenString)
		result = append(result, CommentResponse{
			ID:        c.ID,
			CreatedAt: formatTime(c.CreatedAt),
			UpdatedAt: formatTime(c.UpdatedAt),
			Body:      c.Body,
			Author:    *profile,
		})
	}
	return &result, nil
}

func DeleteComment(commentID uint, tokenString string) *api_errors.E {
	email, mailErr := auth.GetEmailFromTokenString(tokenString)
	if mailErr != nil {
		return api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")
	}

	user, uErr := models.GetUser(email)
	if uErr != nil {
		return api_errors.NewError(http.StatusUnauthorized).Add("token", "token invalid")
	}

	comment, cErr := models.GetComment(commentID)
	if cErr != nil {
		return nil
	}

	if comment.AuthorID != user.ID {
		return api_errors.NewError(http.StatusForbidden).Add("author", "can not delete other people's articles")
	}

	err := models.DeleteComment(commentID)
	if err != nil {
		return api_errors.NewError(http.StatusInternalServerError).Add("comment", "could not delete comment")
	}

	return nil
}
