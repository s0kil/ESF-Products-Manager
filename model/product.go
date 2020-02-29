package model

import (
	"upper.io/db.v3"

	"github.com/s0kil/ESF-Products-Manager/fault"
)

type Product struct {
	id    uint   `db:"id,omitempty"`
	Title string `db:"title" json:"title"`
}

func All(table db.Collection) (products []Product) {
	err := table.Find().All(&products)
	if err != nil {
		fault.Report(err, "Could Not Get Products")
	}

	return
}

func (p Product) New(table db.Collection) error {
	_, err := table.Insert(p)
	return err
}
