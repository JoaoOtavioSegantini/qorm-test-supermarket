package products

import (
	"net/http"
	"strings"

	"github.com/joaootav/system_supermarket/config/utils"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/render"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Controller products controller
type Controller struct {
	View *render.Render
}

// Index products index page
func (ctrl Controller) Index(w http.ResponseWriter, req *http.Request) {
	var (
		Products []models.Product
		tx       = utils.GetDB(req)
	)

	tx.Preload("ColorVariations").Find(&Products)

	ctrl.View.Execute("index", map[string]interface{}{"Products": Products}, req, w)
}

// Gender products gender page
func (ctrl Controller) Gender(w http.ResponseWriter, req *http.Request) {
	var (
		Products []models.Product
		tx       = utils.GetDB(req)
	)

	param := utils.URLParam("gender", req)
	genre := cases.Title(language.Und).String(param)

	tx.Where(&models.Product{Gender: genre}).Preload("ColorVariations").Find(&Products)

	ctrl.View.Execute("gender", map[string]interface{}{"Products": Products}, req, w)
}

// Show product show page
func (ctrl Controller) Show(w http.ResponseWriter, req *http.Request) {
	var (
		product        models.Product
		colorVariation models.ColorVariation
		codes          = strings.Split(utils.URLParam("code", req), "_")
		productCode    = codes[0]
		colorCode      string
		tx             = utils.GetDB(req)
	)

	if len(codes) > 1 {
		colorCode = codes[1]
	}

	if tx.Preload("Category").Where(&models.Product{Code: productCode}).First(&product).RecordNotFound() {
		http.Redirect(w, req, "/", http.StatusFound)
	}

	tx.Preload("Product").Preload("Color").Preload("SizeVariations.Size").Where(&models.ColorVariation{ProductID: product.ID, ColorCode: colorCode}).First(&colorVariation)

	ctrl.View.Execute("show", map[string]interface{}{"CurrentColorVariation": colorVariation}, req, w)
}

// Category category show page
func (ctrl Controller) Category(w http.ResponseWriter, req *http.Request) {
	var (
		category models.Category
		Products []models.Product
		tx       = utils.GetDB(req)
	)

	if tx.Where("code = ?", utils.URLParam("code", req)).First(&category).RecordNotFound() {
		http.Redirect(w, req, "/", http.StatusFound)
	}

	tx.Where(&models.Product{CategoryID: category.ID}).Preload("ColorVariations").Find(&Products)

	ctrl.View.Execute("category", map[string]interface{}{"CategoryName": category.Nome, "Products": Products}, req, w)
}
