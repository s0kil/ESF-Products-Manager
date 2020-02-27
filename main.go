//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

func init() {
	e := godotenv.Load()
	fault(e, "Failed To Load .env File")
}

type Product struct {
	Title string `form:"title"`
}

func main() {
	var (
		DbHost   = os.Getenv("DB_HOST")
		DbName   = os.Getenv("DB_NAME")
		DbUser   = os.Getenv("DB_USER")
		DbSource = fmt.Sprintf("postgresql://%s@%s/%s?sslmode=disable", DbUser, DbHost, DbName)

		AppProductionMode, _ = strconv.ParseBool(os.Getenv("APP_PRODUCTION_MODE"))
	)

	View := jet.NewHTMLSet("./views")
	View.SetDevelopmentMode(!AppProductionMode)

	DB := func() (r *sql.DB) {
		r, e := sql.Open("postgres", DbSource)
		fault(e, "Database Connection Error")
		return
	}()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) {
		c.Set("Content-Type", "text/html")
		c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) {
		c.Send(renderView(View, "home.jet"))
	})

	// Product
	productEndPoint := app.Group("/product")
	// Product List All
	productEndPoint.Get("/", func(c *fiber.Ctx) {
		templateName := "product/index.jet"
		view, e := View.GetTemplate(templateName)
		fault(e, fmt.Sprintf("Failed To Get %s Template", templateName))

		var products []Product
		rows, e := DB.Query("SELECT title FROM products")
		fault(e, "Failed To Select Products From DB")

		for rows.Next() {
			var productTitle string
			e = rows.Scan(&productTitle)
			fault(e, "Failed To Scan Product")
			products = append(products, Product{Title: productTitle})
		}

		var writer bytes.Buffer
		vars := make(jet.VarMap)
		vars.Set("products", &[]Product{})
		e = view.Execute(&writer, vars, products)
		fault(e, fmt.Sprintf("Error When Executing %s Template", templateName))

		c.Send(&writer)
	})
	// Product Create
	productEndPoint.Get("/new", func(c *fiber.Ctx) {
		c.Send(renderView(View, "product/new.jet"))
	})
	productEndPoint.Post("/new", func(c *fiber.Ctx) {
		// Parse Form
		var product Product
		e := c.BodyParser(&product)
		fault(e, "Failed To Parse Form Body")

		_, e = DB.Exec(fmt.Sprintf("INSERT INTO products (title) VALUES ('%s')", product.Title))
		fault(e, "Failed To Insert Product Into DB")

		// TODO: Info Flash With Status (Success, Failure)
		c.Redirect("/product/new")
	})

	fmt.Println("Launching Server")
	e := app.Listen(os.Getenv("APP_PORT"))
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
