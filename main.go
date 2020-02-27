//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"

	"github.com/CloudyKit/jet"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

func fault(err error, reason string) {
	if err != nil {
		log.Fatal(reason, err)
	}
}

func init() {
	e := godotenv.Load()
	fault(e, "Failed To Load .env File")
}

type Product struct {
	ID    int `gorm:"primary_key"`
	Title string
}

func main() {
	var (
		DbHost   = os.Getenv("DB_HOST")
		DbName   = os.Getenv("DB_NAME")
		DbUser   = os.Getenv("DB_USER")
		DbSource = fmt.Sprintf("postgresql://%s@%s/%s?sslmode=disable", DbUser, DbHost, DbName)
	)

	View := jet.NewHTMLSet("./views")

	DB := func() (r *sql.DB) {
		r, e := sql.Open("postgres", DbSource)
		fault(e, "Database Connection Error")
		return
	}()

	rows := func() (r *sql.Rows) {
		r, e := DB.Query("SELECT id, title FROM products")
		fault(e, "Database Query Failed")
		defer r.Close()
		return
	}()

	for rows.Next() {
		var id int
		var title string

		e := rows.Scan(&id, &title)
		fault(e, "DB rows.Scan Failed")

		fmt.Printf("%d %s\n", id, title)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) {
		c.Set("Content-Type", "text/html")

		view, e := View.GetTemplate("home.jet")
		fault(e, "Failed To Get home.jet Template")

		var writer bytes.Buffer
		vars := make(jet.VarMap)
		e = view.Execute(&writer, vars, nil)
		fault(e, "Error when executing home.jet template")

		c.Send(&writer)
	})

	fmt.Println("Launching Server")
	app.Listen(os.Getenv("APP_PORT"))
}
