package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber"             // Web Framework
	"github.com/joho/godotenv"             // Load ENV Variables
	jsoniter "github.com/json-iterator/go" // JSON Serialize, Deserialize
	"upper.io/db.v3/postgresql"            // Database Access Layer

	"github.com/s0kil/ESF-Products-Manager/fault"
	"github.com/s0kil/ESF-Products-Manager/model"
)

func init() {
	err := godotenv.Load()
	fault.Report(err, "Failed To Load .env File")
}

func main() {
	var (
		DatabaseHost   = os.Getenv("DATABASE_HOST")
		DatabaseName   = os.Getenv("DATABASE_NAME")
		DatabaseUser   = os.Getenv("DATABASE_USER")
		DatabaseSource = fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable",
			DatabaseUser, DatabaseHost, DatabaseName)

		ApplicationPort = os.Getenv("APPLICATION_PORT")

		json = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	DatabaseConnectionURL, err := postgresql.ParseURL(DatabaseSource)
	fault.Report(err, "Failed To Parse DatabaseSource")

	Database, err := postgresql.Open(DatabaseConnectionURL)
	fault.Report(err, "Failed To Open Database Session")
	defer Database.Close()

	productsTable := Database.Collection("products")

	App := fiber.New()
	// Load Editor
	App.Get("/", func(c *fiber.Ctx) {
		c.Redirect("/index.html")
	})
	App.Static("/", "./Editor/public")

	// Product
	productEndPoint := App.Group("/product")
	productEndPoint.Use(func(c *fiber.Ctx) {
		c.Set("Content-Type", "application/json")
		c.Next()
	})

	// Product List All
	productEndPoint.Get("/", func(ctx *fiber.Ctx) {
		data := model.All(productsTable)
		result, err := json.Marshal(data)
		fault.Report(err, "Failed To Serialize Data")

		ctx.Send(result)
	})

	// Product Create
	productEndPoint.Post("/new", func(c *fiber.Ctx) {
	})

	err = App.Listen(ApplicationPort)
	fault.Report(err, "Failed To Start Server")
}
