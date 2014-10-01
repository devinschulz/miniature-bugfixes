package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	"github.com/revel/revel"
)

type Register struct {
	App
}

func (c Register) Index() revel.Result {
	auth := c.Auth()
	return c.Render(auth)
}

func (c Register) RegisterPost(user *models.User, verifyPassword string) revel.Result {
	// Validate user data
	user.RegistrationValidate(c.Validation)
	c.Validation.Required(verifyPassword == user.Password).Message("Passwords don't match")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Register.Index)
	}

	// Check to see if email is not in the DB already
	UE, err := c.GetUserByEmail(c.Txn, user.Email)
	checkERROR(err)

	if UE != nil {
		c.Validation.Error("Email already taken").Key("user.Email")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Register.Index)
	}

	// Generate Hashed Password
	user.EncryptedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Create the record
	c.Txn.Create(user)

	// Set User Session
	c.Session["user"] = user.Email
	c.Flash.Success("Welcome " + user.Name)

	return c.Redirect(App.Index)
}
