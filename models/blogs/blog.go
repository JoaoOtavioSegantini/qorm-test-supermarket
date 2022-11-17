package blogs

import (
	"github.com/jinzhu/gorm"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/publish2"
	"github.com/qor/slug"
)

type Article struct {
	gorm.Model
	Author        models.User
	AuthorID      uint
	Title         string
	Content       string `gorm:"type:text"`
	TitleWithSlug slug.Slug
	publish2.Version
	publish2.Schedule
	publish2.Visible
}

// func (article Article) GetUrl() string {
// 	a := strings.Split(article.TitleWithSlug.Slug, "%7b")
// 	b := strings.Split(a[1], "%7d")

// 	return b[0]
// }
