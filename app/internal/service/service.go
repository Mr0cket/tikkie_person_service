package service

import (
	"errors"

	"github.com/Mr0cket/tikkie_person_service/external/mongo"
	"github.com/Mr0cket/tikkie_person_service/external/sqs"
)

var ErrFailedValidation = errors.New("failed validation")
var ErrExistingPerson = errors.New("person with this phone number already exists")

type Service struct {
	DB  mongo.Mongo
	Sqs sqs.Sqs
}

func (s *Service) HealthCheck() bool {
	// Test database connection
	return s.DB.Test()
}
