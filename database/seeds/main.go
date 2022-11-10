package main

import (
	"fmt"
	"time"

	i18n_database "github.com/qor/i18n/backends/database"
	"github.com/qor/publish2"

	"github.com/joaootav/system_supermarket/config/auth"
	"github.com/joaootav/system_supermarket/database"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/providers/password"
	"github.com/qor/notification"
	"github.com/qor/qor"
)

var (
	AdminUser    *models.User
	Notification = notification.New(&notification.Config{})
	DraftDB      = database.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff)
	Tables       = []interface{}{
		&auth_identity.AuthIdentity{},
		&models.UserGroup{},
		&models.User{},
		&models.Category{},
		&models.Product{},
		&i18n_database.Translation{},
		&notification.QorNotification{},
	}
)

func createAdminUsers() {
	AdminUser = &models.User{}
	AdminUser.Email = "test@test.com"
	AdminUser.Confirmed = true
	AdminUser.Name = "QOR Admin"
	AdminUser.Role = "Admin"
	DraftDB.Create(AdminUser)

	provider := auth.Auth.GetProvider("password").(*password.Provider)
	hashedPassword, _ := provider.Encryptor.Digest("testing")
	now := time.Now()

	authIdentity := &auth_identity.AuthIdentity{}
	authIdentity.Provider = "password"
	authIdentity.UID = AdminUser.Email
	authIdentity.EncryptedPassword = hashedPassword
	authIdentity.UserID = fmt.Sprint(AdminUser.ID)
	authIdentity.ConfirmedAt = &now

	DraftDB.Create(authIdentity)

	// Send welcome notification
	Notification.Send(&notification.Message{
		From:        AdminUser,
		To:          AdminUser,
		Title:       "Welcome To QOR Admin",
		Body:        "Welcome To QOR Admin",
		MessageType: "info",
	}, &qor.Context{DB: DraftDB})
}

func main() {
	createRecords()
}

func createRecords() {
	createAdminUsers()
	fmt.Println("--> Created admin users.")
}
