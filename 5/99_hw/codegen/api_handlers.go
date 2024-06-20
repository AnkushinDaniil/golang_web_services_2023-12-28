package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (srv *OtherApi) handleCreate(w http.ResponseWriter, r *http.Request) (string, int, error) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("X-Auth") != "100500" {
		return "", http.StatusForbidden, errors.New("unauthorized")
	}
	if r.Method != "POST" {
		return "", http.StatusNotAcceptable, errors.New("bad method")
	}
	r.ParseForm()
	var (
		username string
		account_name string
		class string
		level string
	)

	username = r.Form.Get("username")
	account_name = r.Form.Get("account_name")
	class = r.Form.Get("class")
	level = r.Form.Get("level")


	if class == "" {
		class = "warrior"
	}
	if class != "warrior" && class != "sorcerer" && class != "rouge" {
		return "", http.StatusBadRequest, errors.New("class must be one of [warrior, sorcerer, rouge]")
	}
	levelInt, err := strconv.Atoi(level)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("level must be int")
	}
	if levelInt < 1 {
		return "", http.StatusBadRequest, errors.New("level must be >= 1")
	}
	if levelInt > 50 {
		return "", http.StatusBadRequest, errors.New("level must be <= 50")
	}

	if username == "" {
		return "", http.StatusBadRequest, errors.New("username must be not empty")
	}
	if len(username) < 3 {
		return "", http.StatusBadRequest, errors.New("username len must be >= 3")
	}

	params := OtherCreateParams{
		Username:	username,
		Name:	account_name,
		Class:	class,
		Level:	levelInt,
	}

	responseCreate, err := srv.Create(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			return "", err.(ApiError).HTTPStatus, err
		default:
			return "", http.StatusInternalServerError, err
		}
	}

	response, err := json.Marshal(responseCreate)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return string(response), http.StatusOK, nil
}
func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		response   string
		statusCode int
	)
	switch r.URL.Path {
	case "/user/create":
		response, statusCode, err = srv.handleCreate(w, r)
	default:
		response, statusCode, err = "", 404, errors.New("unknown method")
	}
	w.WriteHeader(statusCode)
	if err != nil {
		w.Write([]byte("{\"error\": \""))
		w.Write([]byte(err.Error()))
		w.Write([]byte("\"}"))
	} else {
		w.Write([]byte("{\"error\":\"\", \"response\":"))
		w.Write([]byte(response))
		w.Write([]byte("}"))
	}
}

func (srv *MyApi) handleProfile(w http.ResponseWriter, r *http.Request) (string, int, error) {
	w.Header().Set("Content-Type", "application/json")
	var (
		login string
	)

	if r.Method != "GET" {
	r.ParseForm()
		login = r.Form.Get("login")
	} else {
		login = r.URL.Query().Get("login")
	}

	if login == "" {
		return "", http.StatusBadRequest, errors.New("login must be not empty")
	}

	params := ProfileParams{
		Login:	login,
	}

	responseProfile, err := srv.Profile(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			return "", err.(ApiError).HTTPStatus, err
		default:
			return "", http.StatusInternalServerError, err
		}
	}

	response, err := json.Marshal(responseProfile)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return string(response), http.StatusOK, nil
}

func (srv *MyApi) handleCreate(w http.ResponseWriter, r *http.Request) (string, int, error) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("X-Auth") != "100500" {
		return "", http.StatusForbidden, errors.New("unauthorized")
	}
	if r.Method != "POST" {
		return "", http.StatusNotAcceptable, errors.New("bad method")
	}
	r.ParseForm()
	var (
		status string
		age string
		login string
		full_name string
	)

	login = r.Form.Get("login")
	full_name = r.Form.Get("full_name")
	status = r.Form.Get("status")
	age = r.Form.Get("age")

	if login == "" {
		return "", http.StatusBadRequest, errors.New("login must be not empty")
	}
	if len(login) < 10 {
		return "", http.StatusBadRequest, errors.New("login len must be >= 10")
	}


	if status == "" {
		status = "user"
	}
	if status != "user" && status != "moderator" && status != "admin" {
		return "", http.StatusBadRequest, errors.New("status must be one of [user, moderator, admin]")
	}
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("age must be int")
	}
	if ageInt < 0 {
		return "", http.StatusBadRequest, errors.New("age must be >= 0")
	}
	if ageInt > 128 {
		return "", http.StatusBadRequest, errors.New("age must be <= 128")
	}

	params := CreateParams{
		Status:	status,
		Age:	ageInt,
		Login:	login,
		Name:	full_name,
	}

	responseCreate, err := srv.Create(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			return "", err.(ApiError).HTTPStatus, err
		default:
			return "", http.StatusInternalServerError, err
		}
	}

	response, err := json.Marshal(responseCreate)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return string(response), http.StatusOK, nil
}
func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		response   string
		statusCode int
	)
	switch r.URL.Path {
	case "/user/profile":
		response, statusCode, err = srv.handleProfile(w, r)
	case "/user/create":
		response, statusCode, err = srv.handleCreate(w, r)
	default:
		response, statusCode, err = "", 404, errors.New("unknown method")
	}
	w.WriteHeader(statusCode)
	if err != nil {
		w.Write([]byte("{\"error\": \""))
		w.Write([]byte(err.Error()))
		w.Write([]byte("\"}"))
	} else {
		w.Write([]byte("{\"error\":\"\", \"response\":"))
		w.Write([]byte(response))
		w.Write([]byte("}"))
	}
}
