package service

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.mongodb.org/mongo-driver/bson"
)

type ValidationErrors map[string]string

type CreatePersonInput struct {
	FirstName        string `bson:"firstName" json:"firstName" `
	LastName         string `bson:"lastName" json:"lastName"`
	PhoneNumber      string `bson:"phoneNumber" json:"phoneNumber"`
	Address          string `bson:"address" json:"address"` // Using a simple string for simplicity. In real life, you should use a proper address object.
	ValidationErrors ValidationErrors
}

type CreatePersonOutput struct {
	ID string `json:"id"`
}

// High-level steps (Business logic).
// 1. Validate input
// 2. Create a person object in the database
// 2.1 If the person already exists, return an error
// 3. Publish a new event (Module)
// 4. Return ID of the new person

func (s *Service) CreatePerson(input CreatePersonInput) (string, error) {

	// 1. Validate input
	input.ValidationErrors = ValidatePerson(input)

	if len(input.ValidationErrors) > 0 {
		return "", ErrFailedValidation
	}

	// 2. Create a person object in the database
	s.DB.Collection("persons").InsertOne(context.TODO(), input) // TODO: Add support for ctx (Context)

	// 2.1 If the person already exists, return an error
	// Assume phone number is unique per person, so we can check if the person already exists by performing lookup on phone number.
	existingPerson := s.DB.Collection("persons").FindOne(context.TODO(), bson.M{"phoneNumber": input.PhoneNumber})

	if existingPerson.Err() == nil {
		return "", ErrExistingPerson
	}
	// 2.2 Create a new person
	person, err := s.DB.Collection("persons").InsertOne(context.TODO(), input)

	if err != nil {
		return "", err
	}
	log.Printf("Created person ID: %v\n", person.InsertedID)

	// 3. Publish a new event to the SQS queue
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := sqs.NewFromConfig(cfg)

	result, err := svc.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(s.SqsQueueName),
	})
	if err != nil {
		log.Fatalf("Unable to fetch queue URL: %v", err)
	}

	_, err = svc.SendMessage(context.TODO(), &sqs.SendMessageInput{
		DelaySeconds: *aws.Int32(10),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"type": {DataType: aws.String("String"), StringValue: aws.String("personCreated")},
		},
		MessageBody: aws.String("New person created"), // TODO: use a proper JSON object
		QueueUrl:    result.QueueUrl,
	})

	if err != nil {
		log.Fatalf("failed to send message, %v", err)
	}

	return "ID", nil
}

// Validation could be abstracted with a library, but for the purposes of this assignment, will do manually.
// Assume all fields are required.
func ValidatePerson(input CreatePersonInput) ValidationErrors {
	Errors := make(ValidationErrors)

	if len(input.FirstName) < 2 {
		Errors["first_name"] = "Must be at least 2 characters long"
	}

	if len(input.LastName) < 2 {
		Errors["last_name"] = "Must be at least 2 characters long"
	}

	// Check if phone number is 10 digits & starts with '06'
	if len(input.PhoneNumber) != 10 || input.PhoneNumber[0:2] != "06" {
		Errors["phone_number"] = "Must be 10 digits and start with '06'"
	}

	if len(input.Address) < 2 {
		Errors["address"] = "Must be at least 2 characters long"
	}
	return Errors
}
