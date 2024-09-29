package service

import "context"

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
	return "ID", nil
}

func ValidatePerson(input CreatePersonInput) ValidationErrors {
	Errors := make(ValidationErrors)

	if len(input.FirstName) < 2 {
		Errors["first_name"] = "Must be at least 2 characters long"
	}

	if len(input.LastName) < 2 {
		Errors["last_name"] = "Must be at least 2 characters long"
	}

	return Errors
}
