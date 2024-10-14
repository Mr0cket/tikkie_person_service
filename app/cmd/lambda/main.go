package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Mr0cket/tikkie_person_service/external/mongo"
	"github.com/Mr0cket/tikkie_person_service/external/sqs"
	"github.com/Mr0cket/tikkie_person_service/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoServer  string `default:"mongodb://root:example@localhost:27017" split_words:"true"`
	SQSQueueName string `default:"persons" split_words:"true"`
	Database     string `default:"persons"`
	Port         int    `default:"6666"`
	Region       string `default:"europe-west-2"`
}

type Application struct {
	logger  *log.Logger
	service *service.Service
}

var app *Application

func init() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Println("Initialising application")

	var cfg Config
	err := envconfig.Process("APP", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ssmClient := ssm.NewFromConfig(awsCfg)
	result, err := ssmClient.GetParameter(context.TODO(), &ssm.GetParameterInput{Name: aws.String(fmt.Sprintf("/aws/reference/secretsmanager/%s", "mongoUser"))})
	if err != nil {
		log.Fatal(err.Error())
	}

	var mongoUser mongo.MongoUser
	err = json.Unmarshal([]byte(*result.Parameter.Value), &mongoUser)
	if err != nil {
		log.Fatal(err.Error())
	}

	db := mongo.NewClient(context.TODO(), cfg.MongoServer, cfg.Database, mongoUser)
	defer db.Close()

	sqsClient := sqs.NewClient(context.TODO(), awsCfg, cfg.SQSQueueName, cfg.Region)

	app = &Application{
		logger:  logger,
		service: &service.Service{DB: *db, Sqs: *sqsClient},
	}
	logger.Println("Application initialized")
}

func handlers(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return app.createPersonHandler(ctx, request)
	case "GET":
		return app.listPersonsHandler(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	}
}

func main() {
	lambda.Start(handlers)
}
