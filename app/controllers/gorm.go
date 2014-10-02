package controllers

import (
	"database/sql"
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	r "github.com/revel/revel"
)

type GormController struct {
	*r.Controller
	Txn *gorm.DB
}

var Gdb gorm.DB

func InitDB() {

	var (
		driver, spec string
		found        bool
		err          error
	)

	if driver, found = r.Config.String("db.driver"); !found {
		r.ERROR.Fatal("No db.driver found.")
	}
	if spec, found = r.Config.String("db.spec"); !found {
		r.ERROR.Fatal("No db.spec found.")
	}

	Gdb, err = gorm.Open(driver, spec)
	checkPANIC(err)

	Gdb.SetLogger(gorm.Logger{r.INFO})
	Gdb.LogMode(true)

	Gdb.AutoMigrate(&models.User{}, &models.Article{}, &models.Category{}, &models.Tag{})
	Gdb.Model(&models.User{}).AddUniqueIndex("idx_user_email", "email")
	Gdb.Model(&models.Article{}).AddUniqueIndex("idx_title", "title")
	Gdb.Model(&models.Article{}).AddUniqueIndex("idx_slug", "slug")

	r.INFO.Println("Connection made to DB")
}

func (c *GormController) Begin() r.Result {
	txn := Gdb.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	c.Txn = txn
	return nil
}

func (c *GormController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	c.Txn.Commit()
	if err := c.Txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GormController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	c.Txn.Rollback()
	if err := c.Txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
