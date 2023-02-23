package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type ProductCreate struct {
	Code  string `json:"Code" form:"Code"`
	Price uint   `json:"Price" form:"Price"`
}

var db *gorm.DB

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println("failed to load .env file")
	}
	for i := 0; i < 5; i++ {
		var db_conn_str string

		db_url := os.Getenv("DB_URL")
		if db_url != "" {
			db_conn_str = db_url
		} else {
			db_conn_str = fmt.Sprintf("user=postgres password=%s port=5432 dbname=postgres", os.Getenv("POSTGRES_PASSWORD"))
		}

		db, err = gorm.Open(postgres.Open(db_conn_str), &gorm.Config{})
		if err != nil {
			log.Printf("failed to connect database, reconnecting...")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	productApi := api.Group("products")

	productApi.Get("/", func(c *fiber.Ctx) error {
		var products = []Product{}

		result := db.Find(&products)
		if result.Error != nil {
			return result.Error
		}

		return c.JSON(fiber.Map{"products": products})
	})

	productApi.Post("/", func(c *fiber.Ctx) error {
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

	productApi.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		result := db.Delete(&Product{}, id)
		if result.Error != nil {
			return result.Error
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	log.Fatal(app.Listen(":3000"))
}
