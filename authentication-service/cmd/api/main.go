package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"context"
	// _ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main(){
	log.Println("Starting the Authentication Service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Setup the Config
	app := Config{
		DB: conn,
		Models: data.New(conn),

	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string) (*sql.DB,error){
	// db, err := pgx.Connect(context.Background(),os.Getenv("DATABASE_URL"))
	Println(sql.Open(pgx,os.Getenv("DATABASE_URL")));
	// db, err := sql.Open("pgx",dsn)
	if err !=nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	
	return db, err
}


func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	log.Println("DSN: %s",dsn)
	for {
		connnection, err := openDB(dsn)
		log.Println(err)
		if err != nil {
			log.Println("Postgress is not yet ready to serve....")
			counts++
		}else {
			log.Println("Connected to Postgress")
			return connnection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for 2 seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
