package controllers

import (
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	// "github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

type Articles struct {
	App
}

func (c Articles) Index() revel.Result {
	auth := c.Auth()
	var articles []models.Article
	c.Txn.Find(&articles)
	return c.Render(auth, articles)
}

func (c Articles) New() revel.Result {
	auth := c.Auth()
	if auth {
		action := "/articles/new"
		submitButton := "Create Article"
		return c.Render(auth, action, submitButton)
	}
	c.Flash.Error("You must be logged in to create articles")
	return c.Redirect(App.Index)
}

func (c Articles) NewPost(article models.Article) revel.Result {
	article.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}

	c.Txn.Create(article)

	c.Flash.Success("Article Created")

	return c.Redirect(Articles.New)
}
