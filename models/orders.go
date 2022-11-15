package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/qor/transition"
)

var db *gorm.DB

type PaymentMethod = string

const (
	COD        PaymentMethod = "COD"
	AmazonPay  PaymentMethod = "AmazonPay"
	CreditCard PaymentMethod = "CreditCard"
)

type Order struct {
	gorm.Model
	UserID                   *uint
	User                     User
	PaymentAmount            float32
	PaymentTotal             float32
	AbandonedReason          string
	DiscountValue            uint
	DeliveryMethodID         uint `form:"delivery-method"`
	DeliveryMethod           DeliveryMethod
	PaymentMethod            string
	TrackingNumber           *string
	ShippedAt                *time.Time
	ReturnedAt               *time.Time
	CancelledAt              *time.Time
	ShippingAddressID        uint `form:"shippingaddress"`
	ShippingAddress          Address
	BillingAddressID         uint `form:"billingaddress"`
	BillingAddress           Address
	OrderItems               []OrderItem
	AmazonAddressAccessToken string
	AmazonOrderReferenceID   string
	AmazonAuthorizationID    string
	AmazonCaptureID          string
	AmazonRefundID           string
	PaymentLog               string `gorm:"size:655250"`
	transition.Transition
}

func (order Order) ExternalID() string {
	return fmt.Sprint(order.ID)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (order Order) UniqueExternalID() string {
	return fmt.Sprint(order.ID) + "_" + randomString(6)
}

func (order Order) IsCart() bool {
	return order.State == DraftState || order.State == ""
}

func (order Order) Amount() (amount float32) {
	for _, orderItem := range order.OrderItems {
		amount += orderItem.Amount()
	}
	return
}

// DeliveryFee delivery fee
func (order Order) DeliveryFee() (amount float32) {
	return order.DeliveryMethod.Price
}

func (order Order) Total() (total float32) {
	total = order.Amount() - float32(order.DiscountValue)
	total = order.Amount() + float32(order.DeliveryMethod.Price)
	return
}

type OrderItem struct {
	gorm.Model
	OrderID         uint
	SizeVariationID uint `cartitem:"SizeVariationID"`
	SizeVariation   *SizeVariation
	Quantity        uint `cartitem:"Quantity"`
	Price           float32
	DiscountRate    uint
	transition.Transition
}

// IsCart order item's state is cart
func (item OrderItem) IsCart() bool {
	return item.State == DraftState || item.State == ""
}

func (item *OrderItem) loadSizeVariation() {
	if item.SizeVariation == nil {
		item.SizeVariation = &SizeVariation{}
		db.Model(item).Preload("Size").Preload("ColorVariation.Product").Preload("ColorVariation.Color").Association("SizeVariation").Find(item.SizeVariation)
	}
}

// ProductImageURL get product image
func (item *OrderItem) ProductImageURL() string {
	item.loadSizeVariation()
	return item.SizeVariation.ColorVariation.MainImageURL()
}

// SellingPrice order item's selling price
func (item *OrderItem) SellingPrice() float32 {
	if item.IsCart() {
		item.loadSizeVariation()
		return item.SizeVariation.ColorVariation.Product.Price
	}
	return item.Price
}

// ProductName order item's color name
func (item *OrderItem) ProductName() string {
	item.loadSizeVariation()
	return item.SizeVariation.ColorVariation.Product.Name
}

// ColorName order item's color name
func (item *OrderItem) ColorName() string {
	item.loadSizeVariation()
	return item.SizeVariation.ColorVariation.Color.Name
}

// SizeName order item's size name
func (item *OrderItem) SizeName() string {
	item.loadSizeVariation()
	return item.SizeVariation.Size.Name
}

// ProductPath order item's product name
func (item *OrderItem) ProductPath() string {
	item.loadSizeVariation()
	return item.SizeVariation.ColorVariation.ViewPath()
}

// Amount order item's amount
func (item OrderItem) Amount() float32 {
	amount := item.SellingPrice() * float32(item.Quantity)
	if item.DiscountRate > 0 && item.DiscountRate <= 100 {
		amount = amount * float32(100-item.DiscountRate) / 100
	}
	return amount
}

type DeliveryMethod struct {
	gorm.Model

	Name  string
	Price float32
}
