package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaootav/system_supermarket/config/utils"
	"github.com/joaootav/system_supermarket/models/blogs"
	"github.com/joaootav/system_supermarket/tracer"
	"github.com/qor/render"
)

// Controller home controller
type Controller struct {
	View *render.Render
}

// Index home index page
func (ctrl Controller) Index(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, art := tracer.Tracer.Start(ctx, "render-article-page")

	var article blogs.Article

	database := utils.GetDB(req)

	i := database.Preload("Author").Where("title_with_slug = ?", utils.URLParam("slug", req)).Find(&article).RecordNotFound()
	if i {
		utils.AddFlashMessage(w, req, "Record not found", "error")

	}
	w.WriteHeader(http.StatusOK)
	ctrl.View.Execute("index", gin.H{"article": article}, req, w)
	art.End()

}

func (ctrl Controller) List(w http.ResponseWriter, req *http.Request) {
	var articles []blogs.Article

	ctx := req.Context()
	_, list := tracer.Tracer.Start(ctx, "render-list-article-page")

	database := utils.GetDB(req)
	database.Preload("Author").Find(&articles)
	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("list", gin.H{"articles": articles}, req, w)
	list.End()
}
