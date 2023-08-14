package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id           uint   `json:"id"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email" gorm:"unique"`
	Password     string `json:"-" validate:"required"`
	IsAmbassador bool   `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (user *User) SetPassword(password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user.Password = string(passwordHash)

	return nil
}

func (user *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
