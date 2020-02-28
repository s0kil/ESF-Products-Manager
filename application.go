package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"upper.io/db.v3/postgresql"

	"github.com/s0kil/ESF-Products-Manager/fault"
	"github.com/s0kil/ESF-Products-Manager/model"
	"github.com/s0kil/ESF-Products-Manager/view"
)

func init() {
	// Load Environment Variables From .env File
	err := godotenv.Load()
	fault.Report(err, "Failed To Load .env File")
}

func main() {
	var (
		DbHost   = os.Getenv("DB_HOST")
		DbName   = os.Getenv("DB_NAME")
		DbUser   = os.Getenv("DB_USER")
		DbSource = fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable", DbUser, DbHost, DbName)

		AppProductionMode, _ = strconv.ParseBool(os.Getenv("APP_PRODUCTION_MODE"))
	)

	DBConnectionURL, err := postgresql.ParseURL(DbSource)
	fault.Report(err, "Failed To Parse DBSource")

	Database, err := postgresql.Open(DBConnectionURL)
	fault.Report(err, "Failed To Open DB Session")
	defer Database.Close()

	productsTable := Database.Collection("products")

	Views := jet.NewHTMLSet("./view")
	Views.SetDevelopmentMode(!AppProductionMode)

	App := fiber.New()
	App.Use(func(c *fiber.Ctx) {
		c.Set("Content-Type", "text/html")
		c.Next()
	})

	App.Get("/", func(c *fiber.Ctx) {
		c.Send(renderView(Views, "home.jet"))
	})

	// Product
	productEndPoint := App.Group("/product")
	// Product List All
	productEndPoint.Get("/", func(ctx *fiber.Ctx) {
		result := view.ProductIndex(Views, model.All(productsTable))
		ctx.Send(&result)
	})
	// Product Create
	productEndPoint.Get("/new", func(c *fiber.Ctx) {
		c.Send(renderView(Views, "product/new.jet"))
	})
	productEndPoint.Post("/new", func(c *fiber.Ctx) {
		// Parse Form
		var product model.Product
		e := c.BodyParser(&product)
		fault.Report(e, "Failed To Parse Form Body")

		err := product.New(productsTable)
		fault.Report(err, "Failed To Insert Product Into DB")

		// TODO: Info Flash With Status (Success, Failure)
		c.Redirect("/product/new")
	})

	e := App.Listen(os.Getenv("APP_PORT"))
	fault.Report(e, "Failed To Start Server")
}

func renderView(views *jet.Set, templateName string) *bytes.Buffer {
	template, err := views.GetTemplate(templateName)
	fault.Report(err, fmt.Sprintf("Failed To Get %s Template", templateName))

	var writer bytes.Buffer
	err = template.Execute(&writer, nil, nil)
	fault.Report(err, fmt.Sprintf("Error When Executing %s Template", templateName))

	return &writer
}
