package auth

import (
	"net/http"

	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/roles"
)

func init() {
	roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && currentUser.(*models.User).Role == "Admin"
	})
}

type AdminAuth struct {
}

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/auth/logout"
}

func (AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	currentUser, _ := Auth.GetCurrentUser(c.Request).(qor.CurrentUser)
	return currentUser
}
