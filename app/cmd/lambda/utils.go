package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mr0cket/tikkie_person_service/internal/service"
	"github.com/aws/aws-lambda-go/events"
)

// Custom Response writers
func (app *Application) failedValidation(errors service.ValidationErrors) (events.APIGatewayProxyResponse, error) {
	var response events.APIGatewayProxyResponse
	errorResponse := map[string]interface{}{
		"message": "ValidationError",
		"detail":  errors,
	}

	jsonBytes, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Error converting interface (%v) to json: %s", errors, err)

		return response, err
	}
	response = events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       string(jsonBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return response, nil
}
