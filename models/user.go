package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/media"
	"github.com/qor/media/oss"
	"github.com/qor/sorting"
)

type User struct {
	gorm.Model
	sorting.Sorting
	Group                  []UserGroup
	GroupID                uint
	Email                  string `unique;form:"email"`
	Password               string
	Name                   string `form:"name"`
	Gender                 string
	Role                   string
	Birthday               *time.Time
	Balance                float32
	DefaultBillingAddress  uint `form:"default-billing-address"`
	DefaultShippingAddress uint `form:"default-shipping-address"`
	Addresses              []Address
	Avatar                 AvatarImageStorage

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry *time.Time

	// Accepts
	AcceptPrivate bool `form:"accept-private"`
	AcceptLicense bool `form:"accept-license"`
	AcceptNews    bool `form:"accept-news"`
}

type AvatarImageStorage struct{ oss.OSS }

func (AvatarImageStorage) GetSizes() map[string]*media.Size {
	return map[string]*media.Size{
		"small":  {Width: 50, Height: 50},
		"middle": {Width: 120, Height: 120},
		"big":    {Width: 320, Height: 320},
	}
}

type UserGroup struct {
	gorm.Model
	//l10n.Locale
	sorting.Sorting
	Name string //`l10n:"sync"`

}

func (user User) DisplayName() string {
	return user.Email
}

func (user User) AvailableLocales() []string {
	return []string{"en-US", "zh-CN"}
}

type Address struct {
	gorm.Model
	UserID      uint
	ContactName string `form:"contact-name"`
	Phone       string `form:"phone"`
	City        string `form:"city"`
	Address1    string `form:"address1"`
	Address2    string `form:"address2"`
}

func (address Address) Stringify() string {
	return fmt.Sprintf("%v, %v, %v", address.Address2, address.Address1, address.City)
}

type AuthIdentity struct {
	gorm.Model
	Provider          string // phone, email, wechat, github...
	UID               string
	EncryptedPassword string
	AuthInfo          AuthInfo
	UserID            string
	State             string // unconfirmed, confirmed, expired

	Password             string `gorm:"-"`
	PasswordConfirmation string `gorm:"-"`
}

type SignLog struct {
	UserAgent string
	At        *time.Time
	IP        string
}

type AuthInfo struct {
	PhoneVerificationCode       string
	PhoneVerificationCodeExpiry *time.Time
	PhoneConfirmedAt            *time.Time
	UnconfirmedPhone            string // only use when changing phone number

	EmailConfirmedAt *time.Time
	UnconfirmedEmail string // only use when changing email

	SignInCount uint
	SignLogs    []SignLog
}
