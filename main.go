package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var dsn = os.Getenv("DB_URL")

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type ProductCreate struct {
	Code  string `json:"Code" form:"Code"`
	Price uint   `json:"Price" form:"Price"`
}

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")

	api.Get("/product", func(c *fiber.Ctx) error {
		var products = []Product{}

		result := db.Find(&products)
		if result.Error != nil {
			return result.Error
		}

		return c.JSON(fiber.Map{"Products": products})
	})

	api.Post("/product", func(c *fiber.Ctx) error {
		np := new(Product)

		if err := c.BodyParser(np); err != nil {
			return err
		}

		p := Product{Code: np.Code, Price: np.Price}

		result := db.Create(&p)
		if result.Error != nil {
			return result.Error
		}

		return c.JSON(p)
	})

	log.Fatal(app.Listen(":3000"))
}
