package service

import (
	"context"
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ValidationErrors map[string]string

type CreatePersonInput struct {
	FirstName        string           `bson:"firstName" json:"firstName" `
	LastName         string           `bson:"lastName" json:"lastName"`
	PhoneNumber      string           `bson:"phoneNumber" json:"phoneNumber"`
	Address          string           `bson:"address" json:"address"` // Using a simple string for simplicity. In real life, you should use a proper address object.
	ValidationErrors ValidationErrors `bson:"-" json:"-"`
}

type Person struct {
	ID          string `bson:"_id" json:"id"`
	FirstName   string `bson:"firstName" json:"firstName"`
	LastName    string `bson:"lastName" json:"lastName"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
	Address     string `bson:"address" json:"address"` // Using a simple string for simplicity. In real life, you should use a proper address object.
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

func (s *Service) CreatePerson(input *CreatePersonInput) (string, error) {

	// 1. Validate input
	input.ValidationErrors = ValidatePerson(input)

	if len(input.ValidationErrors) > 0 {
		return "", ErrFailedValidation
	}

	// 2. Create a person object in the database
	// 2.1 If the person already exists, return an error
	// Assume phone number is unique per person, so we can check if the person already exists by performing lookup on phone number.
	existingPerson := s.DB.FindDocument("persons_master", bson.M{"phoneNumber": input.PhoneNumber})

	if existingPerson.Err() == nil {
		log.Println(existingPerson.Raw())
		return "", ErrExistingPerson
	}
	// 2.2 Create a new person
	person, err := s.DB.InsertDocument("persons_master", input)

	if err != nil {
		return "", err
	}
	personID := person.InsertedID.(primitive.ObjectID).Hex()

	log.Printf("Created person ID: %v\n", personID)

	// 3. Publish a new event (Module)
	attributes := map[string]string{
		"dataType":    "person",
		"messageType": "create",
		"itemId":      personID,
		"source":      "person-service",
	}
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	log.Println(string(jsonBytes))
	messageID := s.Sqs.SendMessage(attributes, string(jsonBytes))
	log.Printf("Published createPerson event: %v\n", messageID)

	return personID, nil
}

func (s *Service) ListPersons() ([]Person, error) {

	// 1. Get all persons from the database
	cursor := s.DB.ListDocuments("persons_master")
	persons_list := make([]Person, cursor.RemainingBatchLength())

	for cursor.Next(context.TODO()) {
		var person Person
		err := cursor.Decode(&person)
		if err != nil {
			return nil, err
		}
		log.Printf("Person: %v\n", person)
		persons_list = append(persons_list, person)
	}
	return persons_list, nil
}
