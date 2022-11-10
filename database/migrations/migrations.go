package migrations

import (
	"github.com/joaootav/system_supermarket/database"
	"github.com/joaootav/system_supermarket/models"
)

func init() {
	database.DB.AutoMigrate(&models.UserGroup{}, &models.User{})
	database.DB.AutoMigrate(&models.Category{}, &models.Product{})
}

func AutoMigrate(values ...interface{}) {
	for _, value := range values {
		database.DB.AutoMigrate(value)
	}
}
