package models

import "gorm.io/gorm"

type Sale struct {
	gorm.Model
	Total     float64
	ValorPago float64
	Troco     float64
}
