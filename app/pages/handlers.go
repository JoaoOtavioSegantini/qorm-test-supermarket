package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joaootav/system_supermarket/config/utils"
	"github.com/joaootav/system_supermarket/models/blogs"
	"github.com/qor/render"
)

// Controller home controller
type Controller struct {
	View *render.Render
}

// Index home index page
func (ctrl Controller) Index(w http.ResponseWriter, req *http.Request) {
	var article blogs.Article

	database := utils.GetDB(req)

	i := database.Preload("Author").Where("title_with_slug = ?", utils.URLParam("slug", req)).Find(&article).RecordNotFound()
	if i {
		return

	}
	ctrl.View.Execute("index", gin.H{"article": article}, req, w)
}

func (ctrl Controller) List(w http.ResponseWriter, req *http.Request) {
	var articles []blogs.Article

	database := utils.GetDB(req)
	database.Find(&articles)

	ctrl.View.Execute("list", gin.H{"articles": articles}, req, w)
}
