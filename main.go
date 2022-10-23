package main

import (
	"context"
	"github.com/go-pg/pg/v11"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"time"
)

type Page struct {
	Id          int64 //`pg:",discard_unknown_columns"`
	Title       string
	Body        string
	SpaceKey    string    `pg:"-"`
	ParentId    int64     `pg:"parent_id"`
	AuthorId    int64     `pg:"author_id"`
	CreatedAt   time.Time `pg:"-"`
	LastUpdated time.Time `pg:"-"`
}

func main() {

	app := fiber.New()

	pages := make([]Page, 0)

	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: "pages",
		Addr:     os.Getenv("DB_HOST"),
	})
	defer db.Close(context.Background())

	_, err := db.Query(context.Background(), &pages, "select * from pages")
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < 100; i++ {
		page := pages[i]
		log.Println(page)
	}

	// REST api
	app.Get("/pages", func(ctx *fiber.Ctx) error {
		return ctx.JSON(pages)
	})

	app.Listen(":4000")
}
