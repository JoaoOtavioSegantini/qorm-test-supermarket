package config

import (
	"github.com/jinzhu/configor"
	amazonpay "github.com/qor/amazon-pay-sdk-go"
	"github.com/qor/auth/providers/facebook"
	"github.com/qor/auth/providers/github"
	"github.com/qor/auth/providers/google"
	"github.com/qor/auth/providers/twitter"
	"github.com/qor/gomerchant"
	"github.com/qor/location"
	"github.com/qor/mailer"
	"github.com/qor/mailer/logger"
	"github.com/qor/media/oss"
	"github.com/qor/oss/s3"
	"github.com/qor/qor/utils"
	"github.com/qor/redirect_back"
	"github.com/qor/session/manager"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

var Config = struct {
	Port uint `default:"9000" env:"PORT"`
	DB   struct {
		Name     string `env:"DBName" default:"demo.db"`
		Adapter  string `env:"DBAdapter" default:"sqlite"`
		Host     string `env:"DBHost" default:""`
		Port     string `env:"DBPort" default:""`
		User     string `env:"DBUser"`
		Password string `env:"DBPassword"`
	}
	S3 struct {
		AccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
		Region          string `env:"AWS_Region"`
		S3Bucket        string `env:"AWS_Bucket"`
	}
	AmazonPay struct {
		MerchantID   string `env:"AmazonPayMerchantID"`
		AccessKey    string `env:"AmazonPayAccessKey"`
		SecretKey    string `env:"AmazonPaySecretKey"`
		ClientID     string `env:"AmazonPayClientID"`
		ClientSecret string `env:"AmazonPayClientSecret"`
		Sandbox      bool   `env:"AmazonPaySandbox"`
		CurrencyCode string `env:"AmazonPayCurrencyCode" default:"JPY"`
	}
	SMTP         SMTPConfig
	Github       github.Config
	Google       google.Config
	Facebook     facebook.Config
	Twitter      twitter.Config
	GoogleAPIKey string `env:"GoogleAPIKey"`
	BaiduAPIKey  string `env:"BaiduAPIKey"`
}{}

var (
	Root           = utils.AppRoot
	AmazonPay      amazonpay.AmazonPayService
	PaymentGateway gomerchant.PaymentGateway
	Mailer         *mailer.Mailer
	RedirectBack   = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
)

func init() {
	if err := configor.Load(&Config, "config/database.yml", "config/smtp.yml", "config/application.yml"); err != nil {
		panic(err)
	}

	location.GoogleAPIKey = Config.GoogleAPIKey
	location.BaiduAPIKey = Config.BaiduAPIKey

	if Config.S3.AccessKeyID != "" {
		oss.Storage = s3.New(&s3.Config{
			AccessID:  Config.S3.AccessKeyID,
			AccessKey: Config.S3.SecretAccessKey,
			Region:    Config.S3.Region,
			Bucket:    Config.S3.S3Bucket,
		})
	}

	AmazonPay = amazonpay.New(&amazonpay.Config{
		MerchantID: Config.AmazonPay.MerchantID,
		AccessKey:  Config.AmazonPay.AccessKey,
		SecretKey:  Config.AmazonPay.SecretKey,
		Sandbox:    true,
		Region:     "jp",
	})

	// dialer := gomail.NewDialer(Config.SMTP.Host, Config.SMTP.Port, Config.SMTP.User, Config.SMTP.Password)
	// sender, err := dialer.Dial()

	// Mailer = mailer.New(&mailer.Config{
	// 	Sender: gomailer.New(&gomailer.Config{Sender: sender}),
	// })
	Mailer = mailer.New(&mailer.Config{
		Sender: logger.New(&logger.Config{}),
	})

}
