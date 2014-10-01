package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

type Register struct {
	App
}

func (c Register) Index() revel.Result {
	return c.Render()
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
	UE, err := GetUserByEmail(c.Txn, user.Email)
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

	return c.Redirect(App.Index)
}

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := db.Where(&models.User{Email: email}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserById(db *gorm.DB, id int64) (*models.User, error) {
	var user models.User
	err := db.Where(&models.User{Id: id}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
