package models

import (
	"github.com/revel/revel"
)

type Category struct {
	Id   int64
	Name string
	Slug string
}

func (category *Category) Validate(v *revel.Validation) {
	v.Required(category.Name).
		Message("Name is Required").
		Key("category.Name")
	v.MinSize(category.Name, 3).
		Message("Name must be at least 3 characters").
		Key("category.Name")
	v.Required(category.Slug).
		Message("Slug is Required").
		Key("category.Slug")
	v.MinSize(category.Slug, 3).
		Message("Slug must be at least 3 characters").
		Key("category.Slug")
}
