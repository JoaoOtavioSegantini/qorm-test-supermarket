package orders

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/joaootav/system_supermarket/config/utils"
	"github.com/joaootav/system_supermarket/models"
	amazonpay "github.com/qor/amazon-pay-sdk-go"
	"github.com/qor/gomerchant"
	"github.com/unrolled/render"

	qorrender "github.com/qor/render"
	"github.com/qor/responder"
	"github.com/qor/session/manager"
)

var Render = render.New()

// Controller products controller
type Controller struct {
	View *qorrender.Render
}

var decoder = schema.NewDecoder()

// Cart shopping cart
func (ctrl Controller) Cart(w http.ResponseWriter, req *http.Request) {
	order := getCurrentOrderWithItems(w, req)
	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("cart", map[string]interface{}{"Order": order}, req, w)
}

// Checkout checkout shopping cart
func (ctrl Controller) Checkout(w http.ResponseWriter, req *http.Request) {
	hasAmazon := req.URL.Query().Get("access_token")
	order := getCurrentOrderWithItems(w, req)
	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("checkout", map[string]interface{}{"Order": order, "HasAmazon": hasAmazon}, req, w)
}

// Complete complete shopping cart
func (ctrl Controller) Complete(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	order := getCurrentOrder(w, req)
	if order.AmazonOrderReferenceID = req.Form.Get("amazon_order_reference_id"); order.AmazonOrderReferenceID != "" {
		order.AmazonAddressAccessToken = req.Form.Get("amazon_address_access_token")
		tx := utils.GetDB(req)
		err := models.OrderState.Trigger("checkout", order, tx, "")

		if err == nil {
			tx.Save(order)
			http.Redirect(w, req, "/cart/success", http.StatusFound)
			return
		}
		utils.AddFlashMessage(w, req, err.Error(), "error")
	} else {
		utils.AddFlashMessage(w, req, "Order Reference ID not Found", "error")
	}

	http.Redirect(w, req, "/cart", http.StatusFound)
}

// CompleteCreditCard complete shopping cart with credit card
func (ctrl Controller) CompleteCreditCard(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	order := getCurrentOrder(w, req)

	expMonth, _ := strconv.Atoi(req.Form.Get("exp_month"))
	expYear, _ := strconv.Atoi(req.Form.Get("exp_year"))

	creditCard := gomerchant.CreditCard{
		Name:     req.Form.Get("name"),
		Number:   req.Form.Get("creditcard"),
		CVC:      req.Form.Get("cvv"),
		ExpYear:  uint(expYear),
		ExpMonth: uint(expMonth),
	}

	if creditCard.ValidNumber() {
		// TODO integrate with https://github.com/qor/gomerchant to handle those information
		tx := utils.GetDB(req)
		err := models.OrderState.Trigger("checkout", order, tx, "")

		if err == nil {
			tx.Save(order)
			http.Redirect(w, req, "/cart/success", http.StatusFound)
			return
		}
	}

	utils.AddFlashMessage(w, req, "Invalid Credit Card", "error")
	http.Redirect(w, req, "/cart", http.StatusFound)
}

// CheckoutSuccess checkout success page
func (ctrl Controller) CheckoutSuccess(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("success", map[string]interface{}{}, req, w)
}

type updateCartInput struct {
	SizeVariationID  uint `schema:"size_variation_id"`
	Quantity         uint `schema:"quantity"`
	ProductID        uint `schema:"product_id"`
	ColorVariationID uint `schema:"color_variation_id"`
}

// UpdateCart update shopping cart
func (ctrl Controller) UpdateCart(w http.ResponseWriter, req *http.Request) {
	var (
		input updateCartInput
		tx    = utils.GetDB(req)
	)

	req.ParseForm()
	decoder.Decode(&input, req.PostForm)

	order := getCurrentOrder(w, req)

	if input.Quantity > 0 {
		tx.Where(&models.OrderItem{OrderID: order.ID, SizeVariationID: input.SizeVariationID}).
			Assign(&models.OrderItem{Quantity: input.Quantity}).
			FirstOrCreate(&models.OrderItem{})
	} else {
		tx.Where(&models.OrderItem{OrderID: order.ID, SizeVariationID: input.SizeVariationID}).
			Delete(&models.OrderItem{})
	}

	responder.With("html", func() {
		http.Redirect(w, req, "/cart", http.StatusFound)
	}).With([]string{"json", "xml"}, func() {
		Render.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Respond(req)
}

// AmazonCallback amazon callback
func (ctrl Controller) AmazonCallback(w http.ResponseWriter, req *http.Request) {
	ipn, ok := amazonpay.VerifyIPNRequest(req)
	fmt.Printf("%#v\n", ipn)
	fmt.Printf("%#v\n", ok)
}

func getCurrentOrder(w http.ResponseWriter, req *http.Request) *models.Order {
	var (
		order       = models.Order{}
		cartID      = manager.SessionManager.Get(req, "cart_id")
		currentUser = utils.GetCurrentUser(req)
		tx          = utils.GetDB(req)
	)

	if cartID != "" {
		// log.Print("caiu ak")

		if currentUser != nil && !tx.NewRecord(currentUser) {
			// log.Print("ak")
			// log.Printf("%v", order)
			// log.Printf("cart id: %v", cartID)

			if !tx.First(&order, "ID = ? AND (user_id = ? OR user_id IS NULL)", cartID, currentUser.ID).RecordNotFound() && order.UserID == nil {
				//	log.Print("nessa função")

				tx.Model(&order).Update("UserID", currentUser.ID)
			}
		} else {
			//log.Print("caiu aqui")
			tx.First(&order, "ID = ? AND user_id IS NULL", cartID)
		}
	}

	// only create new shopping cart if updating
	if tx.NewRecord(order) || !order.IsCart() {
		order = models.Order{}
		if req.Method != "GET" {
			if currentUser != nil && !tx.NewRecord(currentUser) {
				order.UserID = &currentUser.ID
			}

			tx.Create(&order)
		}
	}
	manager.SessionManager.Add(w, req, "cart_id", order.ID)

	return &order
}

func getCurrentOrderWithItems(w http.ResponseWriter, req *http.Request) *models.Order {
	order := getCurrentOrder(w, req)

	if tx := utils.GetDB(req); !tx.NewRecord(order) {
		tx.Model(order).Association("OrderItems").Find(&order.OrderItems)

	}
	//	res2B, _ := json.Marshal(order)
	//	fmt.Println(string(res2B))
	return order
}
