package controllers

import (
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

type Categories struct {
	*revel.Controller
	App
}

func (c Categories) Index() revel.Result {
	auth := c.Auth()
	if auth {
		categories := c.GetCategories(c.Txn)
		return c.Render(auth, categories)
	}
	c.Flash.Error("You must be logged in to create categories")
	return c.Redirect(App.Index)
}

func (c Categories) New() revel.Result {
	auth := c.Auth()
	if auth {
		action := "/categories/new"
		submitButton := "Create Category"

		return c.Render(auth, action, submitButton)
	}
	c.Flash.Error("You must be logged in to create categories")
	return c.Redirect(App.Index)
}

func (c Categories) NewCategory(category models.Category) revel.Result {
	auth := c.Auth()
	if auth {
		category.Validate(c.Validation)
		if c.Validation.HasErrors() {
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(Categories.New)
		}

		c.Txn.Create(category)
		c.Flash.Success("Category Created")

		return c.Redirect(Categories.Index)
	}
	c.Flash.Error("You must be logged in to create categories")
	return c.Redirect(App.Index)
}

func (c Categories) GetCategories(db *gorm.DB) []models.Category {
	var categories []models.Category
	db.Find(&categories)
	return categories
}
