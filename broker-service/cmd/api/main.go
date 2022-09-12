package main

import (
	"fmt"
	"log"
	"net/http"
)


const webPort = "8010"

type Config struct {}

func main() {
	app := Config{}

	log.Printf("Starting the Broker Service On Port: %s \n",webPort)

	// define HTTP Service
	src := &http.Server {
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start the Server
	err := src.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

