package db

import "time"

type Page struct {
	Id          int64
	Title       string
	Body        string
	SpaceKey    string    `bun:"space_key"`
	ParentId    int64     `bun:"parent_id"`
	AuthorId    int64     `bun:"author_id"`
	CreatedAt   time.Time `bun:"created_at"`
	LastUpdated time.Time `bun:"-"`
}

type Author struct {
	Id   int64
	Name string
}
