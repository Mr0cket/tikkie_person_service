package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/Mr0cket/tikkie_person_service/external/mongo"
	"github.com/Mr0cket/tikkie_person_service/external/sqs"
	"github.com/Mr0cket/tikkie_person_service/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Environment  string `default:"dev"`
	MongoServer  string `split_words:"true"`
	SQSQueueName string `split_words:"true"`
	Database     string `default:"persons"`
	Region       string `envconfig:"AWS_REGION"`
	DbSecretArn  string `default:"abcd" split_words:"true"`
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
	err := envconfig.Process("app", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	secretsClient := secretsmanager.NewFromConfig(awsCfg)
	result, err := secretsClient.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(cfg.DbSecretArn),
	})
	if err != nil {
		log.Fatalf("error getting secret value: %v", err)
	}

	var mongoUser mongo.MongoUser
	err = json.Unmarshal([]byte(*result.SecretString), &mongoUser)
	if err != nil {
		log.Fatalf("error unmarshalling secret: %v", err)
	}

	db := mongo.NewClient(context.TODO(), cfg.MongoServer, cfg.Database, mongoUser)

	sqsClient := sqs.NewClient(context.TODO(), awsCfg, cfg.SQSQueueName, cfg.Region)

	app = &Application{
		logger:  logger,
		service: &service.Service{DB: *db, Sqs: *sqsClient},
	}
	logger.Println("Application initialized")
}

func handlers(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	app.logger.Println(request.HTTPMethod, request.Path, request.RequestContext.Identity.SourceIP)

	switch request.HTTPMethod {
	case "POST":
		return app.createPersonHandler(ctx, request)
	case "GET":
		if request.Path == "/health" {
			return app.healthHandler(ctx, request)
		}
		return app.listPersonsHandler(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	}
}

func main() {
	defer app.service.DB.Close()
	lambda.Start(handlers)
}
