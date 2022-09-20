package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"context"
	"fmt"
	"net/http"
	"time"
	"log-service/data"
)

const (
	webPort= "80"
	rpcPort= "5001"
	grpcPort= "50001"
	mongoURL= "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	// Connect to MongoDb
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context to disconnect

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Start Web Server
	// go app.serve() // it will run in back and terminate the main
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	log.Println(err)
	if err != nil {
		log.Panic()
	}

}

// func (app *Config) serve() {
	
// }

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	
	// Change the setting from the environment to pick
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err !=nil {
		log.Println("Error Connecting:",err)
		return nil,err
	}

	return conn,nil

}
