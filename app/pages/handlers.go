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
	database.Preload("Author").Find(&article)

	ctrl.View.Execute("index", gin.H{"article": article}, req, w)
}
