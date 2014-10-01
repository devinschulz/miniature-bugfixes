package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/revel/revel"
)

type Login struct {
	*revel.Controller
	App
}

func (c Login) Index() revel.Result {
	auth := c.Auth()
	return c.Render(auth)
}

func (c Login) PostLogin(email, password string, remember bool) revel.Result {
	if email == "" || password == "" {
		c.Flash.Error("Missing fields")
		return c.Redirect(Login.Index)
	}
	ValidateLogin(c.Validation, email, password)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Login.Index)
	}

	user, err := c.GetUserByEmail(c.Txn, email)
	checkERROR(err)
	if user != nil {
		err = bcrypt.CompareHashAndPassword(user.EncryptedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = email
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome Back " + user.Name)
			return c.Redirect(App.Index)
		}
		c.Flash.Error("Incorrect Password")
		return c.Redirect(Login.Index)
	}
	c.Flash.Error("Email Not Found")
	return c.Redirect(Login.Index)
}

func (c Login) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	c.Flash.Success("Logged out successfully")
	return c.Redirect(App.Index)
}

func ValidateLogin(v *revel.Validation, email, password string) {
	v.Required(email).
		Message("Email is Required")
	v.Email(email).
		Message("Valid Email Required")
	v.Required(password).
		Message("Password Required")
}
