package service

// Validation could be abstracted with a library, but for the purposes of this assignment, will do manually.
// Assume all fields are required.
func ValidatePerson(input *CreatePersonInput) ValidationErrors {
	Errors := make(ValidationErrors)

	if len(input.FirstName) < 2 {
		Errors["firstName"] = "Must be at least 2 characters long"
	}

	if len(input.LastName) < 2 {
		Errors["lastName"] = "Must be at least 2 characters long"
	}

	// Check if phone number is 10 digits & starts with '06'
	if len(input.PhoneNumber) != 10 || input.PhoneNumber[0:2] != "06" {
		Errors["phoneNumber"] = "Must be 10 digits and start with '06'"
	}

	if len(input.Address) < 2 {
		Errors["address"] = "Must be at least 2 characters long"
	}
	return Errors
}
