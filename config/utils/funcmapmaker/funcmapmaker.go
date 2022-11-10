package funcmapmaker

import (
	"html/template"
	"net/http"

	"github.com/joaootav/system_supermarket/config/admin"
	"github.com/joaootav/system_supermarket/config/i18n"
	"github.com/joaootav/system_supermarket/config/utils"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/action_bar"
	"github.com/qor/i18n/inline_edit"
	"github.com/qor/qor"
	"github.com/qor/render"
	"github.com/qor/session"
	"github.com/qor/session/manager"
)

// AddFuncMapMaker add FuncMapMaker to view
func AddFuncMapMaker(view *render.Render) *render.Render {
	oldFuncMapMaker := view.FuncMapMaker
	view.FuncMapMaker = func(render *render.Render, req *http.Request, w http.ResponseWriter) template.FuncMap {
		funcMap := template.FuncMap{}
		if oldFuncMapMaker != nil {
			funcMap = oldFuncMapMaker(render, req, w)
		}

		// Add `t` method
		for key, fc := range inline_edit.FuncMap(i18n.I18n, utils.GetCurrentLocale(req), utils.GetEditMode(w, req)) {
			funcMap[key] = fc
		}

		for key, value := range admin.ActionBar.FuncMap(w, req) {
			funcMap[key] = value
		}

		// Add `action_bar` method
		funcMap["render_action_bar"] = func() template.HTML {
			return admin.ActionBar.Actions(action_bar.Action{Name: "Edit SEO", Link: models.SEOCollection.SEOSettingURL("/help")}).Render(w, req)
		}

		funcMap["render_seo_tag"] = func() template.HTML {
			return models.SEOCollection.Render(&qor.Context{DB: utils.GetDB(req)}, "Default Page")
		}

		funcMap["flashes"] = func() []session.Message {
			return manager.SessionManager.Flashes(w, req)
		}

		funcMap["current_locale"] = func() string {
			return utils.GetCurrentLocale(req)
		}

		funcMap["current_user"] = func() *models.User {
			return utils.GetCurrentUser(req)
		}

		return funcMap
	}

	return view
}
