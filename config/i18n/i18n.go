package i18n

import (
	"path/filepath"

	db "github.com/joaootav/system_supermarket/database"

	"github.com/joaootav/system_supermarket/config"
	"github.com/qor/i18n"
	"github.com/qor/i18n/backends/database"
	"github.com/qor/i18n/backends/yaml"
)

var I18n *i18n.I18n

func init() {
	I18n = i18n.New(database.New(db.DB), yaml.New(filepath.Join(config.Root, "config/locales")))
}
