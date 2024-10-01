package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/flow"
	"github.com/kelseyhightower/envconfig"

	"github.com/Mr0cket/tikkie_person_service/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Application struct {
	logger  *log.Logger
	service *service.Service
}

type Config struct {
	MongoURI string `default:"mongodb://root:example@localhost:27017"`
	SQSQueue string `default:"persons"`
	Database string `default:"persons"`
	Port     int    `default:"6666"`
}

func main() {
	var cfg Config
	err := envconfig.Process("APP", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	logger := log.New(os.Stdout, "app", log.LstdFlags|log.Llongfile)

	// Setup MongoDB connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.MongoURI).SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := mongoClient.Database(cfg.Database)

	// Test MongoDB connection
	var result bson.M
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	logger.Printf("Connected to MongoDB: %+v\n", result["ok"])

	app := &Application{
		logger:  logger,
		service: &service.Service{DB: db, SqsQueueName: cfg.SQSQueue},
	}

	mux := flow.New()
	mux.HandleFunc("/register", app.createPersonHandler, "POST")

	portString := fmt.Sprintf(":%d", cfg.Port)

	logger.Printf("Starting server on %s\n", portString)
	err = http.ListenAndServe(portString, mux)
	logger.Fatal(err)
}
