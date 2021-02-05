package models

type UserProfile struct {
	Id          string `json:"id" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	DOB         uint64 `json:"dob" validate:"required"`
	Height      uint   `json:"height" validate:"required"`
	Weight      uint   `json:"weight" validate:"required"`
	Race        string `json:"race"`
}

type UserProfileUpdated struct {
	DisplayName string `json:"display_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Height      uint   `json:"height,omitempty"`
	Weight      uint   `json:"weight,omitempty"`
}
