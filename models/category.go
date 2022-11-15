package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/l10n"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

type Category struct {
	gorm.Model
	Nome string

	l10n.Locale
	sorting.Sorting
	Code string
}

func (category Category) Validate(db *gorm.DB) {
	if strings.TrimSpace(category.Nome) == "" {
		db.AddError(validations.NewError(category, "Nome", "Name can not be empty"))
	}
}

func (category Category) DefaultPath() string {
	if len(category.Code) > 0 {
		return fmt.Sprintf("/category/%s", category.Code)
	}
	return "/"
}
