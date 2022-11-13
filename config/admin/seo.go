package admin

import (
	"fmt"

	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/admin"
	qor_seo "github.com/qor/seo"
)

// SetupSEO add seo
func SetupSEO(Admin *admin.Admin) {
	models.SEOCollection = qor_seo.New("Common SEO")
	models.SEOCollection.RegisterGlobalVaribles(&models.SEOGlobalSetting{SiteName: "Qor Shop"})
	models.SEOCollection.SettingResource = Admin.AddResource(&models.MySEOSetting{}, &admin.Config{Invisible: true})
	models.SEOCollection.RegisterSEO(&qor_seo.SEO{
		Name: "Default Page",
	})
	models.SEOCollection.RegisterSEO(&qor_seo.SEO{
		Name:     "Product Page",
		Varibles: []string{"Name", "Code", "Fornecedor"},
		Context: func(objects ...interface{}) map[string]string {
			product := objects[0].(models.Product)
			context := make(map[string]string)
			context["Name"] = product.Name
			context["Code"] = fmt.Sprint(product.ID)
			context["Fornecedor"] = product.Supplier.Nome
			return context
		},
	})
	models.SEOCollection.RegisterSEO(&qor_seo.SEO{
		Name:     "Category Page",
		Varibles: []string{"Name", "Code"},
		Context: func(objects ...interface{}) map[string]string {
			category := objects[0].(models.Category)
			context := make(map[string]string)
			context["Name"] = category.Nome
			context["Code"] = fmt.Sprint(category.ID)
			return context
		},
	})
	Admin.AddResource(models.SEOCollection, &admin.Config{Name: "SEO Setting", Menu: []string{"Site Management"}, Singleton: true, Priority: 2})
}
