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
	"github.com/qor/l10n"
	"github.com/qor/media"
	"github.com/qor/publish2"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

var DB *gorm.DB
var dbError error
var Mux *http.ServeMux

// connect with chosen database
func Connect() {

	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}

	// initialize an HTTP request multiplexer
	Mux = http.NewServeMux()

	log.Println("Connected to Database!")
}

func init() {
	dbConfig := config.Config.DB
	if config.Config.DB.Adapter == "mysql" {

		DB, dbError = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		DB = DB.Set("gorm:table_options", "CHARSET=utf8")

	} else if config.Config.DB.Adapter == "postgres" {

		DB, dbError = gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name))

	} else if config.Config.DB.Adapter == "sqlite" {

		DB, dbError = gorm.Open("sqlite3", fmt.Sprintf("%v/%v", "./database/", dbConfig.Name))

	} else if config.Config.DB.Adapter == "TEST" {
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
		publish2.RegisterCallbacks(DB)

	} else {
		panic(dbError)
	}
}
