package models

import (
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"github.com/russross/blackfriday"
	"html/template"
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
	Meta        map[string]interface{} `sql:"-"`
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
	v.MinSize(article.Title, 3).
		Message("Title must be at least 3 characters").
		Key("article.Title")
	v.Required(article.Slug).
		Message("Slug is Required").
		Key("article.Slug")
	v.MinSize(article.Slug, 3).
		Message("Slug must be at least 3 characters").
		Key("article.Slug")
	v.Required(article.Content).
		Message("Content is Required").
		Key("article.Content")
}

const trimLength = 300

func (article *Article) AddArticleMeta(db *gorm.DB) {
	if article.Meta == nil {
		article.Meta = make(map[string]interface{})
	}
	article.Meta["markdown"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Content))))
	if len(article.Content) > trimLength {
		article.Meta["teaser"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Content[0:trimLength]))))
	} else {
		article.Meta["teaser"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Content[0:len(article.Content)]))))
	}
}
