package auth

import (
	"time"

	"github.com/joaootav/system_supermarket/config"
	"github.com/joaootav/system_supermarket/database"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/authority"
	"github.com/qor/auth/providers/password"
	"github.com/qor/auth_themes/clean"
)

var (
	// Auth initialize Auth for Authentication
	Auth = clean.New(&auth.Config{
		DB:         database.DB,
		UserModel:  models.User{},
		Redirector: auth.Redirector{RedirectBack: config.RedirectBack},
	})

	// Authority initialize Authority for Authorization
	Authority = authority.New(&authority.Config{
		Auth: Auth,
	})
)

func init() {
	database.DB.AutoMigrate(&auth_identity.AuthIdentity{})
	Auth.Render.DefaultLayout = "login_layout"

	//	Auth.RegisterProvider(github.New(&config.Config.Github))
	//	Auth.RegisterProvider(google.New(&config.Config.Google))
	//	Auth.RegisterProvider(facebook.New(&config.Config.Facebook))
	//	Auth.RegisterProvider(twitter.New(&config.Config.Twitter))
	Auth.RegisterProvider(password.New(&password.Config{}))
	Authority.Register("logged_in_half_hour", authority.Rule{TimeoutSinceLastLogin: time.Minute * 30})
}
