package database

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jinzhu/gorm"
	"github.com/joaootav/system_supermarket/config"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/admin"
	"github.com/qor/l10n"
	"github.com/qor/media"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

var DB *gorm.DB
var dbError error
var Mux *http.ServeMux

func Connect(connectionString string) {

	configuration()

	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}

	Admin := admin.New(&admin.AdminConfig{DB: DB})

	Admin.AddResource(&models.Supplier{})
	Admin.AddResource(&models.Product{})

	// initialize an HTTP request multiplexer
	Mux = http.NewServeMux()

	// Mount admin interface to mux
	Admin.MountTo("/admin", Mux)

	log.Println("Connected to Database!")
}

func Migrate() {
	DB.AutoMigrate(
		&models.Category{},
		&models.Inventory{},
		&models.Product{},
		&models.OnSale{},
		&models.OutSale{},
		&models.Sale{},
		&models.Supplier{},
	)
	log.Println("Database Migration Completed!")
}

func configuration() {
	dbConfig := config.Config.DB
	if config.Config.DB.Adapter == "mysql" {
		DB, dbError = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		DB = DB.Set("gorm:table_options", "CHARSET=utf8")
	} else if config.Config.DB.Adapter == "postgres" {
		DB, dbError = gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name))
	} else if config.Config.DB.Adapter == "sqlite" {
		dbUri := ":memory:"
		DB, dbError = gorm.Open("sqlite3", dbUri)
	} else {
		panic(errors.New("not supported database adapter"))
	}

	if dbError == nil {
		if os.Getenv("DEBUG") != "" {
			DB.LogMode(true)
		}
		l10n.RegisterCallbacks(DB)
		sorting.RegisterCallbacks(DB)
		validations.RegisterCallbacks(DB)
		media.RegisterCallbacks(DB)
	} else {
		panic(dbError)
	}
}
