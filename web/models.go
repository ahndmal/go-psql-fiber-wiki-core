package web

import "time"

type Page struct {
	Id          int64
	Title       string
	Body        string
	SpaceKey    string
	ParentId    int64
	AuthorId    int64
	CreatedAt   time.Time
	LastUpdated time.Time
}

type Author struct {
	Id   int64
	Name string
}
