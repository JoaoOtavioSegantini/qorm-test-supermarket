package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OutSale struct {
	gorm.Model
	Producto     Product
	ValorDaVenda float64
	Data         time.Time
}
