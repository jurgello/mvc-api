package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	fmt.Println("Hit Broker")

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) Index(w http.ResponseWriter, r *http.Request) {
	var data struct {
		BrokerURL string
	}
	//data.BrokerURL = os.Getenv("BROKER_URL")
	data.BrokerURL = "http://10.107.16.159:8080"
	app.render(w, "index.page.gohtml", data)

}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	log.Println("authentication here")
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println("bad request", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	//user, err := app.Models.User.GetByEmail(requestPayload.Email)
	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println("wrong email", err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		log.Println("wrong password", err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// log authentication
	// err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	// if err != nil {
	// 	app.errorJSON(w, err)
	// 	return
	// }

	payLoad := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payLoad)
}
