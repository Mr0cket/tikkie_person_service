package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Mr0cket/tikkie_person_service/internal/service"
)

func (app *Application) createPersonHandler(w http.ResponseWriter, r *http.Request) {
	var input service.CreatePersonInput

	s := r.ContentLength
	bodyBuffer := make([]byte, s)
	r.Body.Read(bodyBuffer)
	err := json.Unmarshal(bodyBuffer, &input)
	if err != nil {
		log.Fatalln("Unable to parse JSON body")
		return
	}

	personID, err := app.service.CreatePerson(&input)

	if err != nil {
		if errors.Is(err, service.ErrFailedValidation) {
			// TODO: Return the validation errors in JSON format.
			log.Println("Validation errors in person payload:", input)
			app.failedValidation(w, r, input.ValidationErrors)
			return
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)

	// TODO: Change to return the personID in JSON format.
	w.Write([]byte(personID))
}

func (app *Application) listPersonsHandler(w http.ResponseWriter, r *http.Request) {
	persons, err := app.service.ListPersons()
	if err != nil {
		app.serverError(w, err)
		return
	}
	personsJSON, err := json.Marshal(persons)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(personsJSON)
}
