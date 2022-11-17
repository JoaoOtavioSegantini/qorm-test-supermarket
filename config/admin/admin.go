package admin

import (
	"github.com/joaootav/system_supermarket/config/auth"
	"github.com/joaootav/system_supermarket/config/i18n"
	"github.com/joaootav/system_supermarket/database"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/action_bar"
	"github.com/qor/admin"
	"github.com/qor/help"
	"github.com/qor/media/asset_manager"
	"github.com/qor/media/media_library"
	"github.com/qor/publish2"
	"github.com/qor/qor"
)

var Admin *admin.Admin
var ActionBar *action_bar.ActionBar
var AssetManager *admin.Resource

func init() {
	Admin = admin.New(&admin.AdminConfig{DB: database.DB, Auth: auth.AdminAuth{}, I18n: i18n.I18n})
	admin.New(&qor.Config{DB: database.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff)})
	AssetManager = Admin.AddResource(&asset_manager.AssetManager{}, &admin.Config{Invisible: true})

	Admin.AddResource(&asset_manager.AssetManager{}, &admin.Config{Invisible: true})

	Admin.AddResource(&models.UserGroup{}, &admin.Config{Menu: []string{"User Management"}, Name: "UserGroup"})
	//	Admin.AddResource(&models.User{}, &admin.Config{Menu: []string{"User Management"}, Name: "User"})
	//Admin.AddResource(&models.Category{}, &admin.Config{Menu: []string{"Product Management"}})
	//Admin.AddResource(&models.Product{}, &admin.Config{Menu: []string{"Product Management"}})
	//Admin.AddResource(&models.Order{}, &admin.Config{Menu: []string{"Order Management"}})

	// Add action bar
	ActionBar = action_bar.New(Admin)
	ActionBar.RegisterAction(&action_bar.Action{Name: "Admin Dashboard", Link: "/admin"})

	// Add Media Library
	Admin.AddResource(&media_library.MediaLibrary{}, &admin.Config{Menu: []string{"Site Management"}})

	// Add Translations
	Admin.AddResource(i18n.I18n)

	database.DB.AutoMigrate(QorWidgetSetting{})

	// Add Help
	Help := Admin.NewResource(&help.QorHelpEntry{})
	Help.Meta(&admin.Meta{Name: "Body", Config: &admin.RichEditorConfig{AssetManager: AssetManager}})

	SetupNotification(Admin)
	SetupWorker(Admin)
	SetupSEO(Admin)
	SetupWidget(Admin)
	SetupDashboard(Admin)

}
