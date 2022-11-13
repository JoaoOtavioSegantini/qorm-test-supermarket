package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/sorting"
)

type User struct {
	gorm.Model
	sorting.Sorting
	Group     []UserGroup
	GroupID   uint
	Email     string `unique;form:"email"`
	Password  string
	Name      string `form:"name"`
	Gender    string
	Role      string
	Confirmed bool
}

type UserGroup struct {
	gorm.Model
	//l10n.Locale
	sorting.Sorting
	Name string //`l10n:"sync"`

}

func (user User) DisplayName() string {
	return user.Email
}

func (user User) AvailableLocales() []string {
	return []string{"en-US", "zh-CN"}
}
