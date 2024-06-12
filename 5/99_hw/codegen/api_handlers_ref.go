package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		response   string
		statusCode int
	)
	switch r.URL.Path {
	case "/user/create":
		response, statusCode, err = srv.handleCreate(w, r)
	case "/user/profile":
		response, statusCode, err = srv.handleProfile(w, r)
	default:
		response, statusCode, err = "", 404, errors.New("unknown method")
	}
	w.WriteHeader(statusCode)
	if err != nil {
		w.Write([]byte(`{"error": "`))
		w.Write([]byte(err.Error()))
		w.Write([]byte(`"}`))
	} else {
		w.Write([]byte(`{"error":"", "response":`))
		w.Write([]byte(response))
		w.Write([]byte(`}`))
	}
}

func (srv *MyApi) handleProfile(w http.ResponseWriter, r *http.Request) (string, int, error) {
	w.Header().Set("Content-Type", "application/json")
	var login string
	if r.Method == "GET" {
		login = r.URL.Query().Get("login")
	} else {
		r.ParseForm()
		login = r.Form.Get("login")
	}
	if login == "" {
		return "", http.StatusBadRequest, errors.New("login must be not empty")
	}
	params := ProfileParams{
		Login: login,
	}

	user, err := srv.Profile(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			return "", err.(ApiError).HTTPStatus, err
		default:
			return "", http.StatusInternalServerError, err
		}
	}
	response, err := json.Marshal(user)
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
	login := r.Form.Get("login")

	if login == "" {
		return "", http.StatusBadRequest, errors.New("login must be not empty")
	}
	if len(login) < 10 {
		return "", http.StatusBadRequest, errors.New("login len must be >= 10")
	}

	full_name := r.Form.Get("full_name")

	status := r.Form.Get("status")
	if status == "" {
		status = "user"
	} else if status != "user" && status != "moderator" && status != "admin" {
		return "", http.StatusBadRequest, errors.New("status must be one of [user, moderator, admin]")
	}

	ageString := r.Form.Get("age")
	age, err := strconv.Atoi(ageString)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("age must be int")
	}
	if age < 0 {
		return "", http.StatusBadRequest, errors.New("age must be >= 0")
	}
	if age > 128 {
		return "", http.StatusBadRequest, errors.New("age must be <= 128")
	}

	params := CreateParams{
		Login:  login,
		Name:   full_name,
		Status: status,
		Age:    age,
	}

	user, err := srv.Create(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			return "", err.(ApiError).HTTPStatus, err
		default:
			return "", http.StatusInternalServerError, err
		}
	}
	response, err := json.Marshal(user)
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
		w.Write([]byte(`{"error": "`))
		w.Write([]byte(err.Error()))
		w.Write([]byte(`"}`))
	} else {
		w.Write([]byte(`{"error":"", "response":`))
		w.Write([]byte(response))
		w.Write([]byte(`}`))
	}
}

func (srv *OtherApi) handleCreate(w http.ResponseWriter, r *http.Request) (string, int, error) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("X-Auth") != "100500" {
		return "", http.StatusUnauthorized, errors.New("unauthorized")
	}
	if r.Method != "POST" {
		return "", http.StatusNotAcceptable, errors.New("bad method")
	}

	r.ParseForm()
	username := r.Form.Get("username")
	if username == "" {
		return "", http.StatusBadRequest, errors.New("username must be not empty")
	}
	if len(username) < 3 {
		return "", http.StatusBadRequest, errors.New("username must be longer then 3")
	}

	accountName := r.Form.Get("account_name")

	class := r.Form.Get("class")
	if class == "" {
		class = "warrior"
	} else if class != "warrior" && class != "sorcerer" && class != "rouge" {
		return "", http.StatusBadRequest, errors.New("class must be one of [warrior, sorcerer, rouge]")
	}

	levelString := r.Form.Get("level")
	level, err := strconv.Atoi(levelString)
	if err != nil {
		return "", http.StatusBadRequest, errors.New("level must be int")
	}

	if level < 1 {
		return "", http.StatusBadRequest, errors.New("level must be >= 1")
	}
	if level > 50 {
		return "", http.StatusBadRequest, errors.New("level must be <= 50")
	}

	params := OtherCreateParams{
		Username: username,
		Name:     accountName,
		Class:    class,
		Level:    level,
	}

	user, err := srv.Create(r.Context(), params)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	response, err := json.Marshal(user)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return string(response), http.StatusOK, nil
}
