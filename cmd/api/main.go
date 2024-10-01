package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/flow"
	"github.com/kelseyhightower/envconfig"

	"github.com/Mr0cket/tikkie_person_service/external/mongo"
	"github.com/Mr0cket/tikkie_person_service/internal/service"
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
	db := mongo.NewClient(cfg.MongoURI, cfg.Database)
	defer db.Close()
	app := &Application{
		logger:  logger,
		service: &service.Service{DB: *db, SqsQueueName: cfg.SQSQueue},
	}

	mux := flow.New()
	mux.HandleFunc("/register", app.createPersonHandler, "POST")

	portString := fmt.Sprintf(":%d", cfg.Port)

	logger.Printf("Starting server on %s\n", portString)
	err = http.ListenAndServe(portString, mux)
	logger.Fatal(err)
}
