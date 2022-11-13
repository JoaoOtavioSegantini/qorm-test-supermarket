package models

import "github.com/jinzhu/gorm"

type Inventory struct {
	gorm.Model
	Product    Product
	ProductID  uint
	Quantidade float64
}
