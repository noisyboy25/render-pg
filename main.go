package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

var dsn = os.Getenv("DB_URL")

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	db.Create(&Product{Code: "D42", Price: 100})

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		var products = []Product{}
		result := db.Find(&products)
		if result.Error != nil {
			log.Println(result.Error)
		}
		return c.JSON(products)
	})

	log.Fatal(app.Listen(":3000"))
}
