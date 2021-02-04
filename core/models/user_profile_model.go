package models

type UserProfile struct {
	Id          string `json:"id" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	DOB         string `json:"dob" validate:"required"`
	Height      string `json:"height" validate:"required"`
	Weight      string `json:"weight" validate:"required"`
	Race        string `json:"race"`
}
