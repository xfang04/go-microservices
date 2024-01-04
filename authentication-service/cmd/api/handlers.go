package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Authenticate accepts a json payload and attempts to authenticate a user
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate against database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	err = app.logRequest("authenticate", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
		//Data: User{
		//	ID:        1,
		//	FirstName: "Jack",
		//	LastName:  "Smith",
		//	Email:     "jack@smith.com",
		//	Active:    1,
		//},
		Data: user,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err
	}
	logServiceURL := fmt.Sprintf("http://%s/log", "logger-service")
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	clint := &http.Client{}
	_, err = clint.Do(req)
	if err != nil {
		return err
	}
	return nil
}
