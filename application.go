package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/CloudyKit/jet"    // Template Engine
	"github.com/gofiber/fiber"    // Web Framework
	"github.com/joho/godotenv"    // Load ENV Variables
	"github.com/json-iterator/go" // JSON Serialize, Deserialize
	"upper.io/db.v3/postgresql"   // Database Access Layer

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
	productEndPoint.Get("/new", func(c *fiber.Ctx) {
		// TODO
		c.Send("TODO")
	})

	// Product Create
	productEndPoint.Post("/new", func(c *fiber.Ctx) {
		// Parse Form
		var product model.Product
		err := c.BodyParser(&product)
		fault.Report(err, "Failed To Parse Form Body")

		err = product.New(productsTable)
		fault.Report(err, "Failed To Insert Product Into Database")

		// TODO: Info Flash With Status (Success, Failure)
		c.Redirect("/product/new")
	})

	err = App.Listen(ApplicationPort)
	fault.Report(err, "Failed To Start Server")
}

func renderView(views *jet.Set, templateName string) *bytes.Buffer {
	template, err := views.GetTemplate(templateName)
	fault.Report(err, fmt.Sprintf("Failed To Get %s Template", templateName))

	var writer bytes.Buffer
	err = template.Execute(&writer, nil, nil)
	fault.Report(err, fmt.Sprintf("Error When Executing %s Template", templateName))

	return &writer
}
