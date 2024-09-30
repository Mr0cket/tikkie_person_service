package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Mr0cket/tikkie_person_service/internal/service"
)

func (app *Application) registerPersonHandler(w http.ResponseWriter, r *http.Request) {
	var input service.CreatePersonInput

	s := r.ContentLength
	bodyBuffer := make([]byte, s)
	r.Body.Read(bodyBuffer)
	err := json.Unmarshal(bodyBuffer, &s)
	if err != nil {
		log.Fatalln("Unable to parse JSON body")
		return
	}

	personID, err := app.service.CreatePerson(input)

	if err != nil {
		if errors.Is(err, service.ErrFailedValidation) {
			// TODO: Return the validation errors in JSON format.
			log.Fatalln("Validation errors in person payload")
		} else {
			// TODO: Return a generic error message in JSON format. (400 status code)
			log.Fatalln(err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)

	// TODO: Change to return the personID in JSON format.
	w.Write([]byte(personID))
}
