package handlers

import (
	"../api_errors"
	"../auth"
	"../domain"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func readRequest(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request")
	}
	return body, nil
}

func createUserSerialize(data []byte) (domain.UserCreate, error) {
	var requestData map[string]domain.UserCreate
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		return domain.UserCreate{}, fmt.Errorf("could not read request json")
	}
	return requestData["user"], nil
}

func userSignInSerialize(data []byte) (domain.UserSignIn, error) {
	var requestData map[string]domain.UserSignIn
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		return domain.UserSignIn{}, fmt.Errorf("could not read request json")
	}
	return requestData["user"], nil
}

func createUserRead(r *http.Request) (domain.UserCreate, *api_errors.E) {
	userBytes, readErr := readRequest(r)
	if readErr != nil {
		return domain.UserCreate{}, api_errors.NewError(http.StatusBadRequest).Add("body", readErr.Error())
	}
	userData, serErr := createUserSerialize(userBytes)
	if serErr != nil {
		return domain.UserCreate{}, api_errors.NewError(http.StatusBadRequest).Add("body", serErr.Error())
	}
	return userData, nil
}

func signInRead(r *http.Request) (domain.UserSignIn, *api_errors.E) {
	userBytes, readErr := readRequest(r)
	if readErr != nil {
		return domain.UserSignIn{}, api_errors.NewError(http.StatusBadRequest).Add("body", readErr.Error())
	}
	userData, serErr := userSignInSerialize(userBytes)
	if serErr != nil {
		return domain.UserSignIn{}, api_errors.NewError(http.StatusBadRequest).Add("body", serErr.Error())
	}
	return userData, nil
}

type response map[string]interface{}

func newResponse() *response {
	return &response{}
}

func (r *response) addField(field string, value interface{}) *response {
	(*r)[field] = value
	return r
}

func (r *response) toByte() []byte {
	result, _ := json.Marshal(r)
	return result
}

func (r *response) send(w http.ResponseWriter) {
	result, _ := json.Marshal(r)
	log.Println(w.Write(result))
}

func (r *response) print() *response {
	fmt.Printf("%+v\n\n\n\n", r)
	return r
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
	user, err := domain.CreateUser(userData)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(user, "user")))
}

func signInHandle(w http.ResponseWriter, r *http.Request) {
	userData, readErr := signInRead(r)
	if readErr != nil {
		readErr.Send(w)
		return
	}
	user, err := domain.SignIn(userData)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(user, "user")))
}

func getUserHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	user, userErr := domain.GetUser(token)
	if userErr != nil {
		userErr.Send(w)
		return
	}
	log.Println(w.Write(respToByte(user, "user")))
}

func userUpdateSerialize(data []byte) (domain.UserUpdate, error) {
	var requestData map[string]domain.UserUpdate
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		return domain.UserUpdate{}, fmt.Errorf("could not read request json")
	}
	return requestData["user"], nil
}

func userUpdateRead(r *http.Request) (domain.UserUpdate, *api_errors.E) {
	data, dataErr := readRequest(r)
	if dataErr != nil {
		return domain.UserUpdate{}, api_errors.NewError(http.StatusBadRequest).Add("body", dataErr.Error())
	}
	userUpdate, serErr := userUpdateSerialize(data)
	if serErr != nil {
		return domain.UserUpdate{}, api_errors.NewError(http.StatusBadRequest).Add("body", serErr.Error())
	}
	return userUpdate, nil
}

func updateUserHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	user, userErr := domain.GetUser(token)
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
	userResponse, err := domain.UpdateUser(userUpdate, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(userResponse, "user")))
}

func getProfileHandle(w http.ResponseWriter, r *http.Request) {
	token, tErr := GetTokenFromRequest(r)
	if tErr != nil {
		token = ""
	}
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		api_errors.NewError(http.StatusBadRequest).Add("username", "profile request should contain username").Send(w)
		return
	}

	profile, err := domain.GetProfile(username, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(profile, "profile")))
}

func followHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		api_errors.NewError(http.StatusBadRequest).Add("username", "follow request should contain username").Send(w)
		return
	}
	profile, err := domain.FollowUser(username, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(profile, "profile")))
}

func unfollowHandle(w http.ResponseWriter, r *http.Request) {
	token, _ := GetTokenFromRequest(r)
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		api_errors.NewError(http.StatusBadRequest).Add("username", "follow request should contain username").Send(w)
		return
	}
	profile, err := domain.UnfollowUser(username, token)
	if err != nil {
		err.Send(w)
		return
	}
	log.Println(w.Write(respToByte(profile, "profile")))
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" || len(h) == 0 {
		return "", fmt.Errorf("could not get authorization header")
	}
	split := strings.Split(h, " ")
	if len(split) == 0 || !(split[0] == "Bearer" || split[0] == "Token") {
		return "", fmt.Errorf("authorization header should contain Bearer or Token token")
	}
	token := split[1]
	return token, nil
}

func AuthRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetTokenFromRequest(r)
		if err != nil {
			api_errors.NewError(http.StatusUnauthorized).Add("Auth", err.Error()).Send(w)
			return
		}
		valErr := auth.ValidateTokenString(token)
		if valErr != nil {
			api_errors.NewError(http.StatusUnauthorized).Add("Auth", valErr.Error()).Send(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
