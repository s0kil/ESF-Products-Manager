package model

import (
	"upper.io/db.v3"
)

type Product struct {
	ID    uint   `db:"id,omitempty"`
	Title string `db:"title" form:"title"`
}

func All(table db.Collection) []Product {
	var products []Product

	err := table.Find().All(&products)
	if err != nil {
		// TODO: Handle Error
	}

	return products
}

func (p Product) New(table db.Collection) error {
	_, err := table.Insert(p)
	return err
}
