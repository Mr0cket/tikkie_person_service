package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Mr0cket/tikkie_person_service/internal/service"
	"github.com/aws/aws-lambda-go/events"
)

func (app Application) createPersonHandler(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input service.CreatePersonInput
	var response events.APIGatewayProxyResponse

	err := json.Unmarshal([]byte(r.Body), &input)
	if err != nil {
		log.Fatalln("Unable to parse JSON body")
		return response, err
	}

	personID, err := app.service.CreatePerson(ctx, &input)

	if err != nil {
		if errors.Is(err, service.ErrFailedValidation) {
			log.Println("Validation errors in person payload:", input)
			return app.failedValidation(input.ValidationErrors)
		} else {
			return response, err
		}
	}

	response = events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       personID,
	}
	return response, nil
}

func (app *Application) listPersonsHandler(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	persons, err := app.service.ListPersons(ctx)
	var response events.APIGatewayProxyResponse
	if err != nil {

		return response, err
	}
	personsJSON, err := json.Marshal(persons)
	if err != nil {
		return response, err
	}
	response = events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(personsJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return response, nil
}

func (app *Application) healthHandler(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Test database connection
	success := app.service.HealthCheck()
	if !success {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Database connection failed",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}
