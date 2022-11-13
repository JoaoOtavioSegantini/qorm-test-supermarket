package migrations

import (
	"fmt"

	"github.com/joaootav/system_supermarket/database"
	"github.com/joaootav/system_supermarket/models"
	"github.com/joaootav/system_supermarket/models/blogs"
	"github.com/joaootav/system_supermarket/models/settings"
	"github.com/qor/activity"
	"github.com/qor/banner_editor"
	"github.com/qor/help"
	"github.com/qor/media/asset_manager"
	"github.com/qor/transition"
)

func init() {
	AutoMigrate(&models.UserGroup{}, &models.User{})
	AutoMigrate(&models.Category{}, &models.Product{})
	AutoMigrate(
		&models.Inventory{},
		&models.OnSale{},
		&models.OutSale{},
		&models.Sale{},
		&models.Supplier{},
		&models.MySEOSetting{},
		&settings.Setting{},
		&settings.MediaLibrary{},
		&asset_manager.AssetManager{},
		&banner_editor.QorBannerEditorSetting{},
		&transition.StateChangeLog{},
		&activity.QorActivity{},
		&blogs.Page{}, &blogs.Article{},
		&help.QorHelpEntry{},
	)
	fmt.Println("--> Migrations completed.")

}

func AutoMigrate(values ...interface{}) {
	for _, value := range values {
		database.DB.AutoMigrate(value)
	}
}
