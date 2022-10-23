package main

import (
	"context"
	"github.com/go-pg/pg/v11"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"time"
)

const ALL_PAGES = "SELECT * FROM pages"

type Page struct {
	Id          int64 //`pg:",discard_unknown_columns"`
	Title       string
	Body        string
	SpaceKey    string    `pg:"space_key"`
	ParentId    int64     `pg:"parent_id"`
	AuthorId    int64     `pg:"author_id"`
	CreatedAt   time.Time `pg:"created_at"`
	LastUpdated time.Time `pg:"-"`
}

func main() {

	app := fiber.New()

	pages := make([]Page, 0)

	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: "wiki1",
		Addr:     os.Getenv("DB_HOST"),
	})
	defer db.Close(context.Background())

	_, err := db.Query(context.Background(), &pages, ALL_PAGES)
	if err != nil {
		log.Fatalln(err)
	}

	// REST api
	app.Get("/pages", func(ctx *fiber.Ctx) error {
		hundrPages := make([]Page, 100)
		for i := 0; i < 100; i++ {
			hundrPages[i] = pages[i]
		}
		return ctx.JSON(hundrPages)
	})

	app.Get("/all-pages", func(ctx *fiber.Ctx) error {
		return ctx.JSON(pages)
	})

	listErr := app.Listen(":4000")
	if err != nil {
		log.Fatalln(listErr)
	}
}
