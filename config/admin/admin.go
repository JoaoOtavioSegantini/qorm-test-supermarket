package admin

import (
	"github.com/joaootav/system_supermarket/config/auth"
	"github.com/joaootav/system_supermarket/config/i18n"
	"github.com/joaootav/system_supermarket/database"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/action_bar"
	"github.com/qor/admin"
)

var Admin *admin.Admin
var ActionBar *action_bar.ActionBar

func init() {

	Admin = admin.New(&admin.AdminConfig{DB: database.DB, Auth: auth.AdminAuth{}, SiteName: "qor demo", I18n: i18n.I18n})

	Admin.AddResource(&models.UserGroup{}, &admin.Config{Menu: []string{"User Management"}, Name: "UserGroup"})
	Admin.AddResource(&models.User{}, &admin.Config{Menu: []string{"User Management"}, Name: "User"})
	Admin.AddResource(&models.Category{}, &admin.Config{Menu: []string{"Product Management"}})
	Admin.AddResource(&models.Product{}, &admin.Config{Menu: []string{"Product Management"}})

	Admin.AddResource(i18n.I18n)
	ActionBar = action_bar.New(Admin)
	ActionBar.RegisterAction(&action_bar.Action{Name: "Admin Dashboard", Link: "/admin"})

}
