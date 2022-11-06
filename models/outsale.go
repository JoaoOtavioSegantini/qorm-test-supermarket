package models

import (
	"time"

	"gorm.io/gorm"
)

type OutSale struct {
	gorm.Model
	Producto     Product
	ValorDaVenda float64
	Data         time.Time
}
