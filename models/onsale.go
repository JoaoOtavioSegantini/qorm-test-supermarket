package models

import "gorm.io/gorm"

type OnSale struct {
	gorm.Model
	Nome        string
	Produto     Product
	Porcentagem float64
	Status      bool
}
