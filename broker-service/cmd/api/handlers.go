package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"errors"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request){
	payload := jsonResponse {
		Error: false,
		Message: "Hit the Broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter,r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w,r,&requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return 
	}
	
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w,requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("Unknown Action"))
		
	}

}

func (app *Config) authenticate (w http.ResponseWriter,a AuthPayload) {
	// create some json we'll send to the Auth Microcontroller
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// Call the service
	request, err := http.NewRequest("POST","http://authentication-service/authenticate",bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// Get the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("Invalid Credentials!"))
		return 
	}else if response.StatusCode != http.StatusAccepted {
		// app.errorJSON(w, errors.New(`Error in the Service Call! Code: ${response.StatusCode}`))
		app.errorJSON(w, errors.New(response.StatusCode))
		
		return 
	}

	// Create a variable will response body
	var jsonFromService jsonResponse

	// Decode the Json from the Auth Service

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w,err)
		return 
	}

	if jsonFromService.Error {
		app.errorJSON(w,err, http.StatusUnauthorized)
		return
	}

	// Now the valid login
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}