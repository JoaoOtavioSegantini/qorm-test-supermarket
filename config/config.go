package config

import (
	"github.com/jinzhu/configor"
	"github.com/qor/qor/utils"
	"github.com/qor/redirect_back"
	"github.com/qor/session/manager"
)

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
}{}

var (
	Root         = utils.AppRoot
	RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
)

func init() {
	if err := configor.Load(&Config, "config/database.yml"); err != nil {
		panic(err)
	}
}
