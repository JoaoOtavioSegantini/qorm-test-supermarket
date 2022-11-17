package utils

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/joaootav/system_supermarket/config/admin"
	"github.com/joaootav/system_supermarket/config/auth"
	"github.com/joaootav/system_supermarket/database"
	"github.com/joaootav/system_supermarket/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/qor/l10n"
	"github.com/qor/qor/utils"
	"github.com/qor/session"
	"github.com/qor/session/manager"
)

// GetCurrentUser get current user from request
func GetCurrentUser(req *http.Request) *models.User {
	if currentUser, ok := auth.Auth.GetCurrentUser(req).(*models.User); ok {
		return currentUser
	}
	return nil
}

// GetCurrentLocale get current locale from request
func GetCurrentLocale(req *http.Request) string {
	locale := l10n.Global
	if cookie, err := req.Cookie("locale"); err == nil {
		locale = cookie.Value
	}
	//	log.Printf(locale)
	return locale
}

// GetEditMode get edit mode
func GetEditMode(w http.ResponseWriter, req *http.Request) bool {
	return admin.ActionBar.EditMode(w, req)
}

// GetDB get DB from request
func GetDB(req *http.Request) *gorm.DB {
	if db := utils.GetDBFromRequest(req); db != nil {
		return db
	}
	return database.DB
}

// URLParam get url params from request
func URLParam(name string, req *http.Request) string {
	return chi.URLParam(req, name)
}

// AddFlashMessage helper
func AddFlashMessage(w http.ResponseWriter, req *http.Request, message string, mtype string) error {
	return manager.SessionManager.Flash(w, req, session.Message{Message: template.HTML(message), Type: mtype})
}

// HTMLSanitizer HTML sanitizer
var HTMLSanitizer = bluemonday.UGCPolicy()

func FormatPrice(price interface{}) string {
	switch price.(type) {
	case float32, float64:
		return fmt.Sprintf("%0.2f", price)
	case int, uint, int32, int64, uint32, uint64:
		return fmt.Sprintf("%d.00", price)
	}
	return ""
}
