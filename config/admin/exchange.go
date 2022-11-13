package admin

import (
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/exchange"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/validations"
)

// ProductExchange product exchange
var ProductExchange = exchange.NewResource(&models.Product{}, exchange.Config{PrimaryField: "Name"})

func init() {
	ProductExchange.Meta(&exchange.Meta{Name: "Name"})
	ProductExchange.Meta(&exchange.Meta{Name: "PrecoDeVenda"})
	ProductExchange.Meta(&exchange.Meta{Name: "PrecoDeCusto"})

	ProductExchange.AddValidator(&resource.Validator{
		Handler: func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			if utils.ToInt(metaValues.Get("PrecoDeCusto").Value) < 100 {
				return validations.NewError(record, "PrecoDeCusto", "price can't less than 100")
			}
			return nil
		},
	})
}
