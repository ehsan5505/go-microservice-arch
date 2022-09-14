package main

import (
	"database/sql"
	"authentication/data"
	"log"
)

const webPort = "80"

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main(){
	log.Println("Starting the Authentication Service")

	// TODO: Connect to the SQL database

	// Setup the Config
	app := Config{}

	src := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort)
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
