package main

import (
	"ambassador/database"
	"ambassador/models"
	"math/rand"

	"github.com/bxcodec/faker/v4"
)

func main() {
	database.Connect()

	for i := 0; i < 30; i++ {
		product := models.Product{
			Title:       faker.Username(),
			Description: faker.Username(),
			Image:       faker.URL() + ".png",
			Price:       float64(rand.Intn(90) + 10),
		}

		database.DB.Create(&product)
	}
}
