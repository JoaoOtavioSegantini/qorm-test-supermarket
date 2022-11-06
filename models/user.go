package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `unique;form:"email"`
	Password  string
	Name      string `form:"name"`
	Gender    string
	Role      string
	Confirmed bool
}

func (user User) DisplayName() string {
	return user.Email
}
