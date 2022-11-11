package models

import "github.com/jinzhu/gorm"

type Sale struct {
	gorm.Model
	Total     float64
	ValorPago float64
	Troco     float64
}
