package main

import (
	"net/http"
	"log"
	"fmt"
)

type Config struct {

}

const webPort = "80"

func main() {
	app := Config{}

	log.Println("Starting the Mail Server On Port ");

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort)
		Handler: app.routes()
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}