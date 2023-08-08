package models

type User struct {
	Id           uint   `json:"id"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     []byte `json:"password" validate:"required"`
	IsAmbassador bool   `json:"is_ambassador"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
