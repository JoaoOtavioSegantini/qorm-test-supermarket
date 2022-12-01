package home

import (
	"net/http"

	"github.com/joaootav/system_supermarket/tracer"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"
	"github.com/qor/render"
)

// Controller home controller
type Controller struct {
	View *render.Render
}

// Index home index page
func (ctrl Controller) Index(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	_, home := tracer.Tracer.Start(ctx, "render-home-page")
	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("index", map[string]interface{}{}, req, w)
	home.End()
}

// SwitchLocale switch locale
func (ctrl Controller) SwitchLocale(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, swlocale := tracer.Tracer.Start(ctx, "switch-locale-language")

	utils.SetCookie(http.Cookie{Name: "locale", Value: req.URL.Query().Get("locale")}, &qor.Context{Request: req, Writer: w})
	http.Redirect(w, req, req.Referer(), http.StatusSeeOther)
	swlocale.End()
}
