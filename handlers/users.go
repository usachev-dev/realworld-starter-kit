package handlers

import (
	"../api_errors"
	"../auth"
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
	if readErr != nil {
		readErr.Send(w)
		return
	}
	user, err := models.CreateUser(userData)
	if err != nil {
		err.Send(w)
		return
	}
	w.Write(respToByte(user, "user"))
}

func signInHandle(w http.ResponseWriter, r *http.Request) {
	userData, readErr := signInRead(r)
	if readErr != nil {
		readErr.Send(w)
		return
	}
	user, err := models.SignIn(userData)
	if err != nil {
		err.Send(w)
		return
	}
	w.Write(respToByte(user, "user"))
}

func getUserHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := auth.GetTokenFromRequest(r)
	user, userErr := models.GetUser(token)
	if userErr != nil {
		userErr.Send(w)
		return
	}
	w.Write(respToByte(user, "user"))
}

func userUpdateSerialize(data []byte) (models.UserUpdate, error) {
	var requestData map[string]models.UserUpdate
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		return models.UserUpdate{}, fmt.Errorf("could not read request json")
	}
	return requestData["user"], nil
}

func userUpdateRead(r *http.Request) (models.UserUpdate, *api_errors.E) {
	data, dataErr := readRequest(r)
	if dataErr != nil {
		return models.UserUpdate{}, api_errors.NewError(http.StatusBadRequest).Add("body", dataErr.Error())
	}
	userUpdate, serErr := userUpdateSerialize(data)
	if serErr != nil {
		return models.UserUpdate{}, api_errors.NewError(http.StatusBadRequest).Add("body", serErr.Error())
	}
	return userUpdate, nil
}

func updateUserHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := auth.GetTokenFromRequest(r)
	user, userErr := models.GetUser(token)
	if userErr != nil {
		userErr.Send(w)
		return
	}
	tokenEmail, emailErr := auth.GetEmailFromTokenString(token)
	if emailErr != nil {
		api_errors.NewError(http.StatusForbidden).Add("email", "token should contain email info").Send(w)
		return
	}
	if user.Email != tokenEmail {
		api_errors.NewError(http.StatusForbidden).Add("email", "cannot update other users").Send(w)
		return
	}
	userUpdate, readErr := userUpdateRead(r)
	if readErr != nil {
		readErr.Send(w)
		return
	}
	userResponse, err := models.UpdateUser(userUpdate, token)
	if err != nil {
		err.Send(w)
		return
	}
	w.Write(respToByte(userResponse, "user"))
}

func getProfileHandle(w http.ResponseWriter, r *http.Request) {
	token, tErr := auth.GetTokenFromRequest(r)
	if tErr != nil {
		token = ""
	}
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		api_errors.NewError(http.StatusBadRequest).Add("username", "profile request should contain username").Send(w)
		return
	}

	profile, err := models.GetProfile(username, token)
	if err != nil {
		err.Send(w)
		return
	}
	w.Write(respToByte(profile, "profile"))
}
