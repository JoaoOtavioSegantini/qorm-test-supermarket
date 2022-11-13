package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type OutSale struct {
	gorm.Model
	Product      Product
	ProductID    uint
	ValorDaVenda float64
	Data         time.Time
}
