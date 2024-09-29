package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrFailedValidation = errors.New("failed validation")

type Service struct {
	DB         *mongo.Database
	SqsQueueID string
}
