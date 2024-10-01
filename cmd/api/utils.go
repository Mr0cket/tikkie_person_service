package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mr0cket/tikkie_person_service/internal/service"
)

// Custom Response writers
func (app *Application) failedValidation(w http.ResponseWriter, r *http.Request, errors service.ValidationErrors) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]interface{}{
		"message": "ValidationError",
		"detail":  errors,
	}

	jsonBytes, err := json.Marshal(errorResponse)
	if err != nil {
		log.Println("Error converting interface (%v) to json: %s", errors, err)
		app.serverError(w, r, err)
		return
	}
	w.Write(jsonBytes)
}

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	errorResponse := map[string]interface{}{
		"message": "Internal Server Error",
		"detail":  err.Error(),
	}
	jsonBytes, _ := json.Marshal(errorResponse) // This we can assume is safely json encodable
	w.Write(jsonBytes)
}

// Generic functions for the application
