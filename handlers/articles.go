package handlers

import (
	"../api_errors"
	"../domain"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func createArticleSerialize(data []byte) (*domain.ArticleCreate, error) {
	var requestData map[string]domain.ArticleCreate
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		return nil, fmt.Errorf("could not read request json")
	}
	result := requestData["article"]
	return &result, nil
}

func createArticleRead(r *http.Request) (*domain.ArticleCreate, *api_errors.E) {
	bytes, readErr := readRequest(r)
	if readErr != nil {
		return nil, api_errors.NewError(http.StatusBadRequest).Add("body", readErr.Error())
	}
	createArticle, serErr := createArticleSerialize(bytes)
	if serErr != nil {
		return nil, api_errors.NewError(http.StatusBadRequest).Add("body", serErr.Error())
	}
	return createArticle, nil
}

func createArticleHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	data, readErr := createArticleRead(r)
	if readErr != nil {
		readErr.Send(w)
		return
	}
	user, err := domain.CreateArticle(*data, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(user, "article")))
}

func getArticleHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		api_errors.NewError(http.StatusBadRequest).Add("slug", "article request should contain slug").Send(w)
		return
	}

	article, err := domain.GetArticle(slug, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(article, "article")))
}

func favoriteArticleHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		api_errors.NewError(http.StatusBadRequest).Add("slug", "article request should contain slug").Send(w)
		return
	}
	article, err := domain.FavoriteArticle(slug, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(article, "article")))
}

func unfavoriteArticleHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		api_errors.NewError(http.StatusBadRequest).Add("slug", "article request should contain slug").Send(w)
		return
	}
	article, err := domain.UnfavoriteArticle(slug, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(article, "article")))
}

func deleteArticleHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		api_errors.NewError(http.StatusBadRequest).Add("slug", "article request should contain slug").Send(w)
		return
	}
	err := domain.DeleteArticle(slug, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write([]byte{}))
}

func updateArticleSerialize(data []byte) (*map[string]interface{}, error) {
	var requestData map[string]map[string]interface{}
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		return nil, fmt.Errorf("could not read request json")
	}
	result := requestData["article"]
	return &result, nil
}

func updateArticleRead(r *http.Request) (*map[string]interface{}, *api_errors.E) {
	bytes, readErr := readRequest(r)
	if readErr != nil {
		return nil, api_errors.NewError(http.StatusBadRequest).Add("body", readErr.Error())
	}
	updateData, serErr := updateArticleSerialize(bytes)
	if serErr != nil {
		return nil, api_errors.NewError(http.StatusBadRequest).Add("body", serErr.Error())
	}
	return updateData, nil
}

func updateArticleHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	vars := mux.Vars(r)
	slug := vars["slug"]
	if slug == "" {
		api_errors.NewError(http.StatusBadRequest).Add("slug", "article request should contain slug").Send(w)
		return
	}

	updateData, readErr := updateArticleRead(r)
	if readErr != nil {
		readErr.Send(w)
		return
	}

	article, err := domain.UpdateArticle(slug, *updateData, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(article, "article")))
}

func listArticlesHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	vars := mux.Vars(r)

	var tag *string = nil
	if _, found := vars["tag"]; found {
		t, _ := vars["tag"]
		tag = &t
	}
	var author *string = nil
	if _, found := vars["tag"]; found {
		a, _ := vars["author"]
		author = &a
	}
	var favorited *string = nil
	if _, found := vars["tag"]; found {
		f, _ := vars["favorited"]
		favorited = &f
	}
	var limit uint = 0
	if _, found := vars["limit"]; found {
		f, _ := vars["limit"]
		scan, err := strconv.ParseInt(f, 0, 64)
		if err == nil {
			limit = uint(scan)
		}
	}
	var offset uint = 0
	if _, found := vars["offset"]; found {
		o, _ := vars["offset"]
		scan, err := strconv.ParseInt(o, 0, 64)
		if err == nil {
			offset = uint(scan)
		}
	}

	result, count, err := domain.ListArticles(tag, author, favorited, limit, offset, token)

	if err != nil {
		err.Send(w)
		return
	}
	newResponse().addField("articles", *result).addField("articlesCount", count).print().send(w)
}
