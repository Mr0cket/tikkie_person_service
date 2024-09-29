package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/flow"

	"github.com/Mr0cket/tikkie_person_service/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	logger  *log.Logger
	service *service.Service
}

func main() {
	mongoURI := flag.String("uri", "NOT_IMPLEMENTED_YET", "connection string (URI) for Mongo")
	sqsQueueID := flag.String("queue", "NOT_IMPLEMENTED_YET", "SQS queue ID")
	port := flag.Int("port", 6666, "Port") // TODO: use env var

	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)

	// Setup MongoDB connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(*mongoURI).SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := mongoClient.Database("persons")

	// Test MongoDB connection
	var result bson.M
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	app := &application{
		logger:  logger,
		service: &service.Service{DB: db, SqsQueueID: *sqsQueueID},
	}

	mux := flow.New()
	mux.HandleFunc("/register", app.registerPersonHandler, "POST")

	portString := fmt.Sprintf(":%d", *port)

	logger.Printf("Starting server on %s\n", portString)
	err = http.ListenAndServe(portString, mux)
	logger.Fatal(err)
}
