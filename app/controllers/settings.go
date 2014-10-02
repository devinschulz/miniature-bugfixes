package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	"github.com/revel/revel"
)

type Settings struct {
	*revel.Controller
	App
}

func (c Settings) Index() revel.Result {
	auth := c.Auth()
	user, err := c.connected()
	if err != nil {
		c.Flash.Error("You must be logged in to edit your settings")
		return c.Redirect(App.Index)
	}
	return c.Render(auth, user)
}

func (c Settings) SettingsPost(user *models.User, verifyPassword string) revel.Result {
	s, err := c.connected()
	if err != nil {
		c.Flash.Error("You must be logged in to edit your settings")
		return c.Redirect(App.Index)
	}

	// Grab current user session id
	user.Id = s.Id
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Settings.Index)
	}

	// If user has a value, rehash
	if user.Password != "" || verifyPassword != "" {
		user.RegistrationValidate(c.Validation)
		c.Validation.Required(verifyPassword)
		c.Validation.Required(user.Password == verifyPassword).Message("Passwords don't match")
		if c.Validation.HasErrors() {
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(Settings.Index)
		}
		bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.EncryptedPassword = bcryptPassword
	} else {
		existingPassword := s.EncryptedPassword
		if existingPassword != nil {
			user.EncryptedPassword = existingPassword
		}
	}

	// Check if email has changed. If so, validate that
	// new email doesn't already exist
	if user.Email != s.Email {
		UE, err := c.GetUserByEmail(c.Txn, user.Email)
		checkERROR(err)
		if UE != nil {
			c.Validation.Error("Email already taken").Key("user.Email")
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(Settings.Index)
		}
	}

	// Update the Record
	c.Txn.Save(user)
	c.Flash.Success("Settings Updated")

	return c.Redirect(Settings.Index)
}
