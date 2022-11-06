package models

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	Produto    Product
	Quantidade float64
}
