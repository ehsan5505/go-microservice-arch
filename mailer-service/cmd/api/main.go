package main

import (
	"net/http"
	"log"
	"fmt"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting the Mail Server On Port ");

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m:= Mail {
		Domain: 		os.Getenv("MAIL_DOMAIN"),
		Host: 			os.Getenv("MAIL_HOST"),
		Port: 			port,
		Username: 	os.Genenv("MAIL_USERNAME"),
		Password: 	os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromName: 	os.Getenv("MAIL_FROM"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
	}

	return m
}