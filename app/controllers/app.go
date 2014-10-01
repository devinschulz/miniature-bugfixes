package controllers

import (
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	"github.com/revel/revel"
)

type App struct {
	GormController
}

func (c App) Index() revel.Result {
	user := []models.User{}
	c.Txn.Find(&user)

	return c.RenderJson(user)
}
