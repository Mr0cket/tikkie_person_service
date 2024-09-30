package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrFailedValidation = errors.New("failed validation")
var ErrExistingPerson = errors.New("Person with this phone number already exists")

type Service struct {
	DB         *mongo.Database
	SqsQueueID string
}
