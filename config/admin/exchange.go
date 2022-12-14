package admin

import (
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/exchange"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/validations"
)

func SetupProductExchange() *exchange.Resource {

	// ProductExchange product exchange
	ProductExchange := exchange.NewResource(&models.Product{}, exchange.Config{PrimaryField: "Code"})

	// Define columns are exportable/importable
	ProductExchange.Meta(&exchange.Meta{Name: "Code"})
	ProductExchange.Meta(&exchange.Meta{Name: "Name"})
	ProductExchange.Meta(&exchange.Meta{Name: "Price"})

	ProductExchange.AddValidator(&resource.Validator{
		Handler: func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			if utils.ToInt(metaValues.Get("Price").Value) < 100 {
				return validations.NewError(record, "Price", "price can't less than 100")
			}
			return nil
		},
	})

	return ProductExchange
}
