package controllers

import (
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

type App struct {
	GormController
}

func (c App) Index() revel.Result {
	auth := c.Auth()

	user := []models.User{}
	c.Txn.Find(&user)

	return c.Render(user, auth)
}

func (c App) Auth() bool {
	auth := false
	if user, err := c.connected(); user != nil {
		checkINFO(err)
		auth = true
	}
	return auth
}

func (c App) connected() (*models.User, error) {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User), nil
	}
	if email, ok := c.Session["user"]; ok {
		return c.GetUserByEmail(c.Txn, email)
	}
	return nil, nil
}

func (c App) GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := db.Where(&models.User{Email: email}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c App) GetUserById(db *gorm.DB, id int64) (*models.User, error) {
	var user models.User
	err := db.Where(&models.User{Id: id}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
