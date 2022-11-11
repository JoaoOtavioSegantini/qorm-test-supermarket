package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joaootav/system_supermarket/app/home"
	"github.com/joaootav/system_supermarket/config/admin"
	"github.com/joaootav/system_supermarket/config/application"
	"github.com/joaootav/system_supermarket/config/auth"
	"github.com/joaootav/system_supermarket/config/utils/funcmapmaker"
	"github.com/joaootav/system_supermarket/database"
	"github.com/qor/publish2"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"
)

var rootMux *http.ServeMux

func Router() *http.ServeMux {
	if rootMux == nil {
		router := chi.NewRouter()

		var (
			Admin       = admin.Admin
			Application = application.New(&application.Config{
				Router: router,
				Admin:  Admin,
				DB:     database.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
			})
		)
		funcmapmaker.AddFuncMapMaker(auth.Auth.Config.Render)

		router.Use(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				// for demo, don't use this for your production site
				w.Header().Add("Access-Control-Allow-Origin", "*")
				handler.ServeHTTP(w, req)
			})
		})

		router.Use(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				req.Header.Del("Authorization")
				handler.ServeHTTP(w, req)
			})
		})

		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				var (
					tx         = database.DB
					qorContext = &qor.Context{Request: req, Writer: w}
				)

				if locale := utils.GetLocale(qorContext); locale != "" {
					tx = tx.Set("l10n:locale", locale)
				}

				ctx := context.WithValue(req.Context(), utils.ContextDBName, publish2.PreviewByDB(tx, qorContext))
				next.ServeHTTP(w, req.WithContext(ctx))
			})
		})

		rootMux = http.NewServeMux()

		rootMux.Handle("/auth/", auth.Auth.NewServeMux())
		rootMux.Handle("/", Application.NewServeMux())
		Application.Use(home.New(&home.Config{}))

		// Mount admin interface to mux
		admin.Admin.MountTo("/admin", database.Mux)

	}

	return rootMux
}
