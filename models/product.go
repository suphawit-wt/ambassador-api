package models

type Product struct {
	Id          uint    `json:"id"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Image       string  `json:"image"`
	Price       float64 `json:"price" validate:"required"`
}
