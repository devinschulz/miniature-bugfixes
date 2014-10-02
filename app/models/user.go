package models

import (
	"github.com/revel/revel"
	"time"
)

type User struct {
	Id                int64
	Name              string
	Email             string
	EncryptedPassword []byte
	Password          string `sql:"-"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (u *User) Validate(v *revel.Validation) {
	v.Required(u.Email).
		Message("Email is Required").
		Key("user.Email")
	v.Email(u.Email).
		Message("Valid Email Required").
		Key("user.Email")
	v.Required(u.Name).
		Message("Name Required").
		Key("user.Name")
}

func (u *User) RegistrationValidate(v *revel.Validation) {
	v.Required(u.Email).
		Message("Email is Required").
		Key("user.Email")
	v.Email(u.Email).
		Message("Valid Email Required").
		Key("user.Email")
	v.Required(u.Name).
		Message("Name Required").
		Key("user.Name")
	v.Required(u.Password).
		Message("Password Required").
		Key("user.Password")
	v.MinSize(u.Password, 6).
		Message("Password must be at least 6 characters").
		Key("user.Password")
}
