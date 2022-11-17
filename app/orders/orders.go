package orders

import (
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/joaootav/system_supermarket/config/application"
	"github.com/joaootav/system_supermarket/config/utils/funcmapmaker"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/activity"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/qor/render"
	"github.com/qor/transition"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}

// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
}

// ConfigureApplication configure application
func (app App) ConfigureApplication(application *application.Application) {
	controller := &Controller{View: render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("orders")}, "app/orders/views")}

	funcmapmaker.AddFuncMapMaker(controller.View)
	app.ConfigureAdmin(application.Admin)

	application.Router.Get("/cart", controller.Cart)
	application.Router.Put("/cart", controller.UpdateCart)
	application.Router.Post("/cart", controller.UpdateCart)
	application.Router.Get("/cart/checkout", controller.Checkout)
	application.Router.Put("/cart/checkout", controller.Checkout)
	application.Router.Post("/cart/complete", controller.Complete)
	application.Router.Post("/cart/complete/creditcard", controller.CompleteCreditCard)
	application.Router.Get("/cart/success", controller.CheckoutSuccess)
	application.Router.Post("/order/callback/amazon", controller.AmazonCallback)
}

// ConfigureAdmin configure admin interface
func (App) ConfigureAdmin(Admin *admin.Admin) {
	Admin.AddMenu(&admin.Menu{Name: "Order Management", Priority: 2})
	// Add Order
	order := Admin.AddResource(&models.Order{}, &admin.Config{Menu: []string{"Order Management"}})
	order.Meta(&admin.Meta{Name: "ShippingAddress", Type: "single_edit"})
	order.Meta(&admin.Meta{Name: "BillingAddress", Type: "single_edit"})
	order.Meta(&admin.Meta{Name: "ShippedAt", Type: "date"})
	order.Meta(&admin.Meta{Name: "PaymentLog", Type: "readonly", FormattedValuer: func(record interface{}, _ *qor.Context) (result interface{}) {
		order := record.(*models.Order)
		return template.HTML(strings.Replace(strings.TrimSpace(order.PaymentLog), "\n", "<br>", -1))
	}})
	order.Meta(&admin.Meta{Name: "DeliveryMethod", Type: "select_one",
		Config: &admin.SelectOneConfig{
			Collection: func(_ interface{}, context *admin.Context) (options [][]string) {
				var methods []models.DeliveryMethod
				context.GetDB().Find(&methods)

				for _, m := range methods {
					idStr := fmt.Sprintf("%d", m.ID)
					var option = []string{idStr, fmt.Sprintf("%s (%0.2f) руб", m.Name, m.Price)}
					options = append(options, option)
				}

				return options
			},
		},
	})

	orderItemMeta := order.Meta(&admin.Meta{Name: "OrderItems"})
	orderItemMeta.Resource.Meta(&admin.Meta{Name: "SizeVariation", Config: &admin.SelectOneConfig{Collection: sizeVariationCollection}})

	// define scopes for Order
	for _, state := range []string{"pending", "processing", "cancelled", "shipped", "paid_cancelled", "returned"} {
		var state = state
		order.Scope(&admin.Scope{
			Name:  state,
			Label: cases.Title(language.Und).String((strings.Replace(state, "_", " ", -1))),
			Group: "Order Status",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(models.Order{Transition: transition.Transition{State: state}})
			},
		})
	}

	// define actions for Order
	type trackingNumberArgument struct {
		TrackingNumber string
	}

	order.Action(&admin.Action{
		Name: "Processing",
		Handler: func(argument *admin.ActionArgument) error {
			for _, order := range argument.FindSelectedRecords() {
				db := argument.Context.GetDB()
				if err := models.OrderState.Trigger("process", order.(*models.Order), db); err != nil {
					return err
				}
				db.Save(order)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*models.Order); ok {
				return order.State == "pending"
			}
			return false
		},
		Modes: []string{"show", "menu_item"},
	})

	order.Action(&admin.Action{
		Name: "Ship",
		Handler: func(argument *admin.ActionArgument) error {
			var (
				tx                     = argument.Context.GetDB().Begin()
				trackingNumberArgument = argument.Argument.(*trackingNumberArgument)
			)

			if trackingNumberArgument.TrackingNumber != "" {
				for _, record := range argument.FindSelectedRecords() {
					order := record.(*models.Order)
					order.TrackingNumber = &trackingNumberArgument.TrackingNumber
					models.OrderState.Trigger("ship", order, tx, "tracking number "+trackingNumberArgument.TrackingNumber)
					if err := tx.Save(order).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			} else {
				return errors.New("invalid shipment number")
			}

			tx.Commit()
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*models.Order); ok {
				return order.State == "processing"
			}
			return false
		},
		Resource: Admin.NewResource(&trackingNumberArgument{}),
		Modes:    []string{"show", "menu_item"},
	})

	order.Action(&admin.Action{
		Name: "Cancel",
		Handler: func(argument *admin.ActionArgument) error {
			for _, order := range argument.FindSelectedRecords() {
				db := argument.Context.GetDB()
				order := order.(*models.Order)
				if err := models.OrderState.Trigger("cancel", order, db); err != nil {
					return err
				}
				db.Save(order)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*models.Order); ok {
				for _, state := range []string{"draft", "pending", "processing", "shipped"} {
					if order.State == state {
						return true
					}
				}
			}
			return false
		},
		Modes: []string{"show", "menu_item"},
	})

	order.IndexAttrs("ID", "User", "PaymentAmount", "ShippedAt", "CancelledAt", "State", "ShippingAddress")
	order.NewAttrs("-DiscountValue", "-AbandonedReason", "-CancelledAt", "-PaymentLog", "-AmazonOrderReferenceID", "-AmazonAddressAccessToken")
	order.EditAttrs("-DiscountValue", "-AbandonedReason", "-CancelledAt", "-State", "-PaymentLog", "-AmazonOrderReferenceID", "-AmazonAddressAccessToken")
	order.ShowAttrs("-DiscountValue", "-State", "-AmazonAddressAccessToken")
	order.SearchAttrs("User.Name", "User.Email", "AmazonOrderReferenceID", "ShippingAddress.Phone", "ShippingAddress.ContactName", "ShippingAddress.Address1", "ShippingAddress.Address2")

	oldSearchHandler := order.SearchHandler
	order.SearchHandler = func(keyword string, context *qor.Context) *gorm.DB {
		return oldSearchHandler(keyword, context).Where("state <> ? AND state <> ?", "", models.DraftState)
	}

	// Add activity for order
	activity.Register(order)

	// Define another resource for same model
	abandonedOrder := Admin.AddResource(&models.Order{}, &admin.Config{Name: "Abandoned Order", Menu: []string{"Order Management"}})
	abandonedOrder.Meta(&admin.Meta{Name: "ShippingAddress", Type: "single_edit"})
	abandonedOrder.Meta(&admin.Meta{Name: "BillingAddress", Type: "single_edit"})

	// Define default scope for abandoned orders
	abandonedOrder.Scope(&admin.Scope{
		Default: true,
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("abandoned_reason IS NOT NULL AND abandoned_reason <> ?", "")
		},
	})

	// Define scopes for abandoned orders
	for _, amount := range []int{5000, 10000, 20000} {
		var amount = amount
		abandonedOrder.Scope(&admin.Scope{
			Name:  fmt.Sprint(amount),
			Group: "Amount Greater Than",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where("payment_amount > ?", amount)
			},
		})
	}

	abandonedOrder.IndexAttrs("-ShippingAddress", "-BillingAddress", "-DiscountValue", "-OrderItems")
	abandonedOrder.NewAttrs("-DiscountValue")
	abandonedOrder.EditAttrs("-DiscountValue")
	abandonedOrder.ShowAttrs("-DiscountValue")

	// Delivery Methods
	Admin.AddResource(&models.DeliveryMethod{}, &admin.Config{Menu: []string{"Site Management"}})
}

func sizeVariationCollection(resource interface{}, context *qor.Context) (results [][]string) {
	for _, sizeVariation := range models.SizeVariations() {
		results = append(results, []string{strconv.Itoa(int(sizeVariation.ID)), sizeVariation.Stringify()})
	}
	return
}
