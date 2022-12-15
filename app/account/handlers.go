package account

import (
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/joaootav/system_supermarket/config/utils"
	"github.com/joaootav/system_supermarket/models"
	"github.com/joaootav/system_supermarket/tracer"
	"github.com/qor/render"
)

// Controller products controller
type Controller struct {
	View *render.Render
}

// Profile profile show page
func (ctrl Controller) Profile(w http.ResponseWriter, req *http.Request) {
	var (
		currentUser                     = utils.GetCurrentUser(req)
		tx                              = utils.GetDB(req)
		billingAddress, shippingAddress models.Address
	)

	ctx := req.Context()

	_, profile := tracer.Tracer.Start(ctx, "render-user-profile-page")

	// TODO refactor
	tx.Model(currentUser).Related(&currentUser.Addresses, "Addresses")
	tx.Model(currentUser).Related(&billingAddress, "DefaultBillingAddress")
	tx.Model(currentUser).Related(&shippingAddress, "DefaultShippingAddress")

	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("profile", map[string]interface{}{
		"CurrentUser": currentUser, "DefaultBillingAddress": billingAddress, "DefaultShippingAddress": shippingAddress,
	}, req, w)

	profile.End()
}

// Orders orders page
func (ctrl Controller) Orders(w http.ResponseWriter, req *http.Request) {
	var (
		Orders      []models.Order
		currentUser = utils.GetCurrentUser(req)
		tx          = utils.GetDB(req)
	)

	ctx := req.Context()

	_, orders := tracer.Tracer.Start(ctx, "render-orders-page")

	tx.Preload("OrderItems").Where("state <> ? AND state != ?", models.DraftState, "").Where(&models.Order{UserID: &currentUser.ID}).Find(&Orders)
	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("orders", map[string]interface{}{"Orders": Orders}, req, w)

	orders.End()
}

// Update update profile page
func (ctrl Controller) Update(w http.ResponseWriter, req *http.Request) {
	var (
		currentUser                     = utils.GetCurrentUser(req)
		tx                              = utils.GetDB(req)
		billingAddress, shippingAddress models.Address
		user                            = models.User{}
		decoder                         = schema.NewDecoder()
	)

	ctx := req.Context()

	_, defaultAdresses := tracer.Tracer.Start(ctx, "render-default-addresses-page")

	req.ParseForm()
	decoder.Decode(&user, req.Form)

	user.AcceptLicense, _ = strconv.ParseBool(req.Form.Get("accept-license"))
	user.AcceptNews, _ = strconv.ParseBool(req.Form.Get("accept-news"))
	user.AcceptPrivate, _ = strconv.ParseBool(req.Form.Get("accept-private"))

	tx.Model(currentUser).Updates(
		models.User{
			Name:                   user.Name,
			Email:                  user.Email,
			DefaultBillingAddress:  user.DefaultBillingAddress,
			DefaultShippingAddress: user.DefaultShippingAddress,
			AcceptPrivate:          user.AcceptPrivate,
			AcceptLicense:          user.AcceptLicense,
			AcceptNews:             user.AcceptNews,
		},
	)

	// TODO refactor
	tx.Model(currentUser).Related(&currentUser.Addresses, "Addresses")
	tx.Model(currentUser).Related(&billingAddress, "DefaultBillingAddress")
	tx.Model(currentUser).Related(&shippingAddress, "DefaultShippingAddress")

	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("profile", map[string]interface{}{
		"CurrentUser": currentUser, "DefaultBillingAddress": billingAddress, "DefaultShippingAddress": shippingAddress,
	}, req, w)

	defaultAdresses.End()
}

// AddCredit add credit
func (ctrl Controller) AddCredit(w http.ResponseWriter, req *http.Request) {
	// FIXME
}

// Update update profile address page
func (ctrl Controller) UpdateAddress(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	_, address := tracer.Tracer.Start(ctx, "render-addresses-page")

	var (
		currentUser                     = utils.GetCurrentUser(req)
		tx                              = utils.GetDB(req)
		billingAddress, shippingAddress models.Address
		decoder                         = schema.NewDecoder()
	)

	req.ParseForm()
	decoder.Decode(&billingAddress, req.Form)
	billingAddress.UserID = currentUser.ID

	tx.Model(models.Address{}).Create(&billingAddress)
	tx.Model(currentUser).Where(&models.User{DefaultBillingAddress: billingAddress.ID}).Update(&currentUser)

	// TODO refactor
	tx.Model(currentUser).Related(&currentUser.Addresses, "Addresses")
	tx.Model(currentUser).Related(&billingAddress, "DefaultBillingAddress")
	tx.Model(currentUser).Related(&shippingAddress, "DefaultShippingAddress")

	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("profile", map[string]interface{}{
		"CurrentUser": currentUser, "DefaultBillingAddress": billingAddress, "DefaultShippingAddress": shippingAddress,
	}, req, w)

	address.End()
}

// Update update profile shipping page
func (ctrl Controller) UpdateShippingAddress(w http.ResponseWriter, req *http.Request) {
	var (
		currentUser                     = utils.GetCurrentUser(req)
		tx                              = utils.GetDB(req)
		billingAddress, shippingAddress models.Address
		decoder                         = schema.NewDecoder()
	)

	ctx := req.Context()

	_, shpaddress := tracer.Tracer.Start(ctx, "render-shipping-address-page")

	req.ParseForm()
	decoder.Decode(&shippingAddress, req.Form)
	shippingAddress.UserID = currentUser.ID

	tx.Model(models.Address{}).Create(&shippingAddress)
	tx.Model(currentUser).Where(&models.User{DefaultShippingAddress: shippingAddress.ID}).Update(&currentUser.DefaultShippingAddress)

	// TODO refactor
	tx.Model(currentUser).Related(&currentUser.Addresses, "Addresses")
	tx.Model(currentUser).Related(&billingAddress, "DefaultBillingAddress")
	tx.Model(currentUser).Related(&shippingAddress, "DefaultShippingAddress")

	w.WriteHeader(http.StatusOK)

	ctrl.View.Execute("profile", map[string]interface{}{
		"CurrentUser": currentUser, "DefaultBillingAddress": billingAddress, "DefaultShippingAddress": shippingAddress,
	}, req, w)

	shpaddress.End()
}
