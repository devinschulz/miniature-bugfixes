package models

import (
	"github.com/revel/revel"
	"time"
)

type Article struct {
	Id          int64
	UserId      int64
	Title       string
	Slug        string
	Published   bool
	Content     string `sql:"type:text"`
	Categories  []Category
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt time.Time
	Tags        []Tag
}

type Category struct {
	Id   int64
	Name string
	Slug string
}

type Tag struct {
	Id   int64
	Name string
	Slug string
}

func (article *Article) Validate(v *revel.Validation) {
	v.Required(article.Title).
		Message("Title is Required").
		Key("article.Title")
	v.Required(article.Content).
		Message("Content is Required").
		Key("article.Content")
}
