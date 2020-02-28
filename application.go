package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"upper.io/db.v3/postgresql"

	"github.com/s0kil/ESF-Products-Manager/model"
)

func init() {
	// Load Environment Variables From .env File
	err := godotenv.Load()
	fault(err, "Failed To Load .env File")
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
	fault(err, "Failed To Parse DBSource")

	DB, err := postgresql.Open(DBConnectionURL)
	fault(err, "Failed To Open DB Session")
	defer DB.Close()

	productsTable := DB.Collection("products")

	View := jet.NewHTMLSet("./view")
	View.SetDevelopmentMode(!AppProductionMode)

	App := fiber.New()
	App.Use(func(c *fiber.Ctx) {
		c.Set("Content-Type", "text/html")
		c.Next()
	})

	App.Get("/", func(c *fiber.Ctx) {
		c.Send(renderView(View, "home.jet"))
	})

	// Product
	productEndPoint := App.Group("/product")
	// Product List All
	productEndPoint.Get("/", func(c *fiber.Ctx) {
		templateName := "product/index.jet"
		view, e := View.GetTemplate(templateName)
		fault(e, fmt.Sprintf("Failed To Get %s Template", templateName))

		var products []model.Product
		products = model.All(productsTable)

		var writer bytes.Buffer
		vars := make(jet.VarMap).Set("products", []model.Product{})
		e = view.Execute(&writer, vars, &products)
		fault(e, fmt.Sprintf("Error When Executing %s Template", templateName))

		c.Send(&writer)
	})
	// Product Create
	productEndPoint.Get("/new", func(c *fiber.Ctx) {
		c.Send(renderView(View, "product/new.jet"))
	})
	productEndPoint.Post("/new", func(c *fiber.Ctx) {
		// Parse Form
		var product model.Product
		e := c.BodyParser(&product)
		fault(e, "Failed To Parse Form Body")

		err := product.New(productsTable)
		fault(err, "Failed To Insert Product Into DB")

		// TODO: Info Flash With Status (Success, Failure)
		c.Redirect("/product/new")
	})

	e := App.Listen(os.Getenv("APP_PORT"))
	fault(e, "Failed To Start Server")
}

func renderView(View *jet.Set, templateName string) *bytes.Buffer {
	view, e := View.GetTemplate(templateName)
	fault(e, fmt.Sprintf("Failed To Get %s Template", templateName))

	var writer bytes.Buffer
	e = view.Execute(&writer, nil, nil)
	fault(e, fmt.Sprintf("Error When Executing %s Template", templateName))

	return &writer
}

func fault(err error, reason string) {
	if err != nil {
		log.Fatal(reason, err)
	}
}
