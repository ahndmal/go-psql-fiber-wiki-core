package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"os"
)

type PagesRepo struct{}

func getDB() *bun.DB {
	dsn := os.Getenv("DB_DSN")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}

func (pr PagesRepo) GetPages(ctx context.Context) []Page {
	pages := make([]Page, 0)
	err := getDB().NewSelect().Model(&pages).Scan(ctx)
	if err != nil {
		fmt.Printf("Error when selcting pages from DB. Err: %v", err)
	}
	defer getDB().Close()

	return pages
}

func (pr PagesRepo) GetPageById(ctx context.Context, id int64) Page {
	var page Page
	err := getDB().NewSelect().Model(&page).Where("id = ?", id).Scan(ctx)
	if err != nil {
		fmt.Printf("Error when selcting pages from DB. Err: %v", err)
	}
	defer getDB().Close()

	return page
}

func (pr PagesRepo) GetNPages(ctx context.Context, limit int) []Page {
	pages := make([]Page, 0)
	err := getDB().NewSelect().Model(&pages).Limit(limit).Scan(ctx)
	if err != nil {
		fmt.Printf("Error when selcting pages from DB. Err: %v", err)
	}

	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Error closing the connection. Err: %v", err)
		}
	}(getDB())

	return pages
}
