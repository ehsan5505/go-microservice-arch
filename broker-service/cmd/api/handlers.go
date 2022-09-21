package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
	"log"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth 	AuthPayload `json:"auth,omitempty"`
	Log		LogPayload 	`json:"log,omitempty"`
	Mail	MailPayload	`json:"mail,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From		string	`json:"from"`
	To			string	`json:"to"`
	Subject	string	`json:"subject"`
	Message	string	`json:"message"`
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
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("Unknown Action"))
		
	}

}

func (app *Config) logItem(w http.ResponseWriter,entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry,"","\t",)

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return 
	}

	request.Header.Set("Content-Type","application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w,err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)

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
		app.errorJSON(w,errors.New(fmt.Sprintf("Status Code: %d",response.StatusCode)))
		// app.errorJSON(w, errors.New(`Error in the Service Call!`))		
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

func (app *Config) sendMail(w http.ResponseWriter,msg MailPayload) {
	
	jsonData, _ := json.MarshalIndent(msg,"","\t")

	// Call the Mail Service
	mailServiceUrl := "http://mailer-service/send"
	
	// Post Mail Request
	request, err := http.NewRequest("POST",mailServiceUrl,bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w,err)
		return 
	}

	request.Header.Set("Content-Type","application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.errorJSON(w,err)
		return 
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w,errors.New("Error Calling Mail Service"))
		return 
	}

	// send back the payload
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message Sent to "+ msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}