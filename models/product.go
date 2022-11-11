package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/publish2"
)

type Product struct {
	gorm.Model
	Name         string
	Categoria    Category
	Fornecedor   Supplier
	PrecoDeCusto float64
	PrecoDeVenda float64
	Mediacao     int
	Status       bool

	publish2.Version
	publish2.Schedule
	publish2.Visible
}
