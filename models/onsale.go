package models

import "github.com/jinzhu/gorm"

type OnSale struct {
	gorm.Model
	Nome        string
	Product     Product
	ProductID   uint
	Porcentagem float64
	Status      bool
}
