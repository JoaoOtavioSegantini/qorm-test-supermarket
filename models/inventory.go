package models

import "github.com/jinzhu/gorm"

type Inventory struct {
	gorm.Model
	Produto    Product
	Quantidade float64
}
