package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"context"
	"log-service/data"
)

const (
	webPort: "80"
	rpcPort: "5001"
	grpcPort: "50001"
	mongoURL: "mongodb://mongo/27017"
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
	go app.serve()

}

funct (app *Config) serve() {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	
	// Change the setting from the environment to pick
	clientOptions.SetAuth(options.Credentials{
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
