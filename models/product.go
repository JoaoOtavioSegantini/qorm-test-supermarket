package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name         string
	Categoria    Category
	Fornecedor   Supplier
	PrecoDeCusto float64
	PrecoDeVenda float64
	Mediacao     int
	Status       bool
}
