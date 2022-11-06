package models

import "gorm.io/gorm"

type Supplier struct {
	gorm.Model
	Nome     string
	Email    string
	Telefone string
}
