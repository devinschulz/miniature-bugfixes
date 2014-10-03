package controllers

import (
	"fmt"
	"github.com/iDevSchulz/miniature-bugfixes/app/models"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"time"
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

func (c Articles) Edit(id int64) revel.Result {
	// TODO: Check for empty id

	auth := c.Auth()

	if auth {
		action := fmt.Sprintf("/articles/update/%v", id)
		submitButton := "Update Article"

		article, err := GetArticleById(c.Txn, id)
		checkERROR(err)

		return c.Render(auth, action, submitButton, article)
	}
	c.Flash.Error("You must be logged in to create articles")
	return c.Redirect(App.Index)
}

func (c Articles) NewArticle(article models.Article) revel.Result {
	article.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}

	// Check if title already exists
	AT, err := GetArticleByTitle(c.Txn, article.Title)
	checkERROR(err)
	if AT != nil {
		c.Validation.Error(article.Title + " Already Exists as a Title").Key("article.Title")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}

	// Check if slug already exists
	AS, err := GetArticleBySlug(c.Txn, article.Slug)
	checkERROR(err)
	if AS != nil {
		c.Validation.Error(article.Slug + "Already Exists as a Slug").Key("article.Slug")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}

	article.UpdatedAt = time.Now()
	article.CreatedAt = time.Now()

	c.Txn.Create(article)

	c.Flash.Success("Article Created")

	return c.Redirect(Articles.New)
}

func (c Articles) UpdateArticle(article models.Article, id int64) revel.Result {
	article.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Articles.New)
	}

	var (
		a, AT, AS *models.Article
		err       error
	)

	currentPath := fmt.Sprintf("/articles/edit/%v", id)

	// Get the state of the article before the update
	a, err = GetArticleById(c.Txn, id)
	checkERROR(err)

	// Check if title already exists
	if a.Title != article.Title {
		AT, err = GetArticleByTitle(c.Txn, article.Title)
		checkERROR(err)
		if AT != nil {
			c.Validation.Error(article.Title + " Already Exists as a Title").Key("article.Title")
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(currentPath)
		}
	}

	// Check if slug already exists
	if a.Slug != article.Slug {
		AS, err = GetArticleBySlug(c.Txn, article.Slug)
		checkERROR(err)
		if AS != nil {
			c.Validation.Error(article.Slug + "Already Exists as a Slug").Key("article.Slug")
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(currentPath)
		}
	}

	article.Id = id
	article.UpdatedAt = time.Now()

	c.Txn.Save(article)
	c.Flash.Success("Article Updated")

	return c.Redirect(currentPath)
}

func (c Articles) Show(slug string) revel.Result {
	auth := c.Auth()
	article, err := GetArticleBySlug(c.Txn, slug)
	checkERROR(err)
	return c.Render(auth, article)
}

func GetArticleByTitle(db *gorm.DB, title string) (*models.Article, error) {
	var article models.Article
	err := db.Where(&models.Article{Title: title}).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func GetArticleBySlug(db *gorm.DB, slug string) (*models.Article, error) {
	var article models.Article
	err := db.Where(&models.Article{Slug: slug}).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func GetArticleById(db *gorm.DB, id int64) (*models.Article, error) {
	var article models.Article
	err := db.Where(&models.Article{Id: id}).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}
