package handlers

import (
	"../api_errors"
	"../models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func readRequest(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request")
	}
	return body, nil
}

func createUserSerialize(data []byte) (models.UserCreate, error) {
	var requestData map[string]models.UserCreate
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		return models.UserCreate{}, fmt.Errorf("could not read request json")
	}
	return requestData["user"], nil
}

func userSignInSerialize(data []byte) (models.UserSignIn, error) {
	var requestData map[string]models.UserSignIn
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		return models.UserSignIn{}, fmt.Errorf("could not read request json")
	}
	return requestData["user"], nil
}

func createUserRead(r *http.Request) (models.UserCreate, *api_errors.E) {
	userBytes, readErr := readRequest(r)
	if readErr != nil {
		return models.UserCreate{}, api_errors.NewError(http.StatusBadRequest).Add("body", readErr.Error())
	}
	userData, serErr := createUserSerialize(userBytes)
	if serErr != nil {
		return models.UserCreate{}, api_errors.NewError(http.StatusBadRequest).Add("body", serErr.Error())
	}
	return userData, nil
}

func signInRead(r *http.Request) (models.UserSignIn, *api_errors.E) {
	userBytes, readErr := readRequest(r)
	if readErr != nil {
		return models.UserSignIn{}, api_errors.NewError(http.StatusBadRequest).Add("body", readErr.Error())
	}
	userData, serErr := userSignInSerialize(userBytes)
	if serErr != nil {
		return models.UserSignIn{}, api_errors.NewError(http.StatusBadRequest).Add("body", serErr.Error())
	}
	return userData, nil
}

func respToByte(value interface{}, field string) []byte {
	m := map[string]interface{}{field: value}
	result, _ := json.Marshal(m)
	return result
}

func createUserHandle(w http.ResponseWriter, r *http.Request) {
	userData, readErr := createUserRead(r)
	if !readErr.IsOk() {
		readErr.Send(w)
		return
	}
	user, err := models.CreateUser(userData)
	if !err.IsOk() {
		err.Send(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respToByte(user, "user"))
}

func signInHandle(w http.ResponseWriter, r *http.Request) {
	userData, readErr := signInRead(r)
	if !readErr.IsOk() {
		readErr.Send(w)
		return
	}
	user, err := models.SignIn(userData)
	if !err.IsOk() {
		err.Send(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respToByte(user, "user"))
}
