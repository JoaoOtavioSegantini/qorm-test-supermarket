package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/publish2"
	qor_seo "github.com/qor/seo"
)

type Product struct {
	gorm.Model
	Name         string
	Category     Category
	CategoryID   uint
	Supplier     Supplier
	SupplierID   uint
	PrecoDeCusto float64
	PrecoDeVenda float64
	Mediacao     int
	Status       bool

	publish2.Version
	publish2.Schedule
	publish2.Visible
}

func (product Product) GetSEO() *qor_seo.SEO {
	return SEOCollection.GetSEO("Product Page")
}
