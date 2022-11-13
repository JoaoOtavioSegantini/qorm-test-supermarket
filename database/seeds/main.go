package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/configor"
	"github.com/qor/banner_editor"
	i18n_database "github.com/qor/i18n/backends/database"
	"github.com/qor/media/oss"
	"github.com/qor/publish2"
	"github.com/qor/seo"
	"github.com/qor/sorting"

	"github.com/joaootav/system_supermarket/config/admin"
	"github.com/joaootav/system_supermarket/config/auth"
	"github.com/joaootav/system_supermarket/database"
	"github.com/joaootav/system_supermarket/database/migrations"
	"github.com/joaootav/system_supermarket/models"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/providers/password"
	"github.com/qor/notification"
	"github.com/qor/qor"
)

var (
	AdminUser    *models.User
	Notification = notification.New(&notification.Config{})
	DraftDB      = database.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff)
	Tables       = []interface{}{
		&auth_identity.AuthIdentity{},
		&models.UserGroup{},
		&models.User{},
		&models.Category{},
		&models.Product{},
		&i18n_database.Translation{},
		&notification.QorNotification{},
	}
)

func createAdminUsers() {
	AdminUser = &models.User{}
	AdminUser.Email = "test@test.com"
	AdminUser.Confirmed = true
	AdminUser.Name = "QOR Admin"
	AdminUser.Role = "Admin"
	DraftDB.Create(AdminUser)

	provider := auth.Auth.GetProvider("password").(*password.Provider)
	hashedPassword, _ := provider.Encryptor.Digest("testing")
	now := time.Now()

	authIdentity := &auth_identity.AuthIdentity{}
	authIdentity.Provider = "password"
	authIdentity.UID = AdminUser.Email
	authIdentity.EncryptedPassword = hashedPassword
	authIdentity.UserID = fmt.Sprint(AdminUser.ID)
	authIdentity.ConfirmedAt = &now

	DraftDB.Create(authIdentity)

	// Send welcome notification
	Notification.Send(&notification.Message{
		From:        AdminUser,
		To:          AdminUser,
		Title:       "Welcome To QOR Admin",
		Body:        "Welcome To QOR Admin",
		MessageType: "info",
	}, &qor.Context{DB: DraftDB})
}

func createSeo() {
	globalSeoSetting := models.MySEOSetting{}
	globalSetting := make(map[string]string)
	globalSetting["SiteName"] = "Qor Demo"
	globalSeoSetting.Setting = seo.Setting{GlobalSetting: globalSetting}
	globalSeoSetting.Name = "QorSeoGlobalSettings"
	globalSeoSetting.LanguageCode = "en-US"
	globalSeoSetting.QorSEOSetting.SetIsGlobalSEO(true)

	if err := database.DB.Create(&globalSeoSetting).Error; err != nil {
		log.Fatalf("create seo (%v) failure, got err %v", globalSeoSetting, err)
	}

	defaultSeo := models.MySEOSetting{}
	defaultSeo.Setting = seo.Setting{Title: "{{SiteName}}", Description: "{{SiteName}} - Default Description", Keywords: "{{SiteName}} - Default Keywords", Type: "Default Page"}
	defaultSeo.Name = "Default Page"
	defaultSeo.LanguageCode = "en-US"
	if err := database.DB.Create(&defaultSeo).Error; err != nil {
		log.Fatalf("create seo (%v) failure, got err %v", defaultSeo, err)
	}

	productSeo := models.MySEOSetting{}
	productSeo.Setting = seo.Setting{Title: "{{SiteName}}", Description: "{{SiteName}} - {{Name}} - {{Code}}", Keywords: "{{SiteName}},{{Name}},{{Code}}", Type: "Product Page"}
	productSeo.Name = "Product Page"
	productSeo.LanguageCode = "en-US"
	if err := database.DB.Create(&productSeo).Error; err != nil {
		log.Fatalf("create seo (%v) failure, got err %v", productSeo, err)
	}

	// seoSetting := models.SEOSetting{}
	// seoSetting.SiteName = Seeds.Seo.SiteName
	// seoSetting.DefaultPage = seo.Setting{Title: Seeds.Seo.DefaultPage.Title, Description: Seeds.Seo.DefaultPage.Description, Keywords: Seeds.Seo.DefaultPage.Keywords}
	// seoSetting.HomePage = seo.Setting{Title: Seeds.Seo.HomePage.Title, Description: Seeds.Seo.HomePage.Description, Keywords: Seeds.Seo.HomePage.Keywords}
	// seoSetting.ProductPage = seo.Setting{Title: Seeds.Seo.ProductPage.Title, Description: Seeds.Seo.ProductPage.Description, Keywords: Seeds.Seo.ProductPage.Keywords}

	// if err := DraftDB.Create(&seoSetting).Error; err != nil {
	// 	log.Fatalf("create seo (%v) failure, got err %v", seoSetting, err)
	// }
}

func createWidgets() {
	// home page banner
	type ImageStorage struct{ oss.OSS }
	topBannerSetting := admin.QorWidgetSetting{}
	topBannerSetting.Name = "home page banner"
	topBannerSetting.Description = "This is a top banner"
	topBannerSetting.WidgetType = "NormalBanner"
	topBannerSetting.GroupName = "Banner"
	topBannerSetting.Scope = "from_google"
	topBannerSetting.Shared = true
	topBannerValue := &struct {
		Title           string
		ButtonTitle     string
		Link            string
		BackgroundImage ImageStorage `sql:"type:varchar(4096)"`
		Logo            ImageStorage `sql:"type:varchar(4096)"`
	}{
		Title:       "Welcome Googlistas!",
		ButtonTitle: "LEARN MORE",
		Link:        "http://getqor.com",
	}
	if file, err := openFileByURL("https://venngage-wordpress-pt.s3.amazonaws.com/uploads/2019/06/33-modelos-de-slide-incr%C3%ADveis-e-dicas-de-design.png"); err == nil {
		defer file.Close()
		topBannerValue.BackgroundImage.Scan(file)
	} else {
		fmt.Printf("open file (%q) failure, got err %v", "banner", err)
	}

	if file, err := openFileByURL("https://venngage-wordpress-pt.s3.amazonaws.com/uploads/2019/06/33-modelos-de-slide-incr%C3%ADveis-e-dicas-de-design.png"); err == nil {
		defer file.Close()
		topBannerValue.Logo.Scan(file)
	} else {
		fmt.Printf("open file (%q) failure, got err %v", "qor_logo", err)
	}

	topBannerSetting.SetSerializableArgumentValue(topBannerValue)
	if err := DraftDB.Create(&topBannerSetting).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", topBannerSetting, err)
	}

	// SlideShow banner
	type slideImage struct {
		Title    string
		SubTitle string
		Button   string
		Link     string
		Image    oss.OSS
	}
	slideshowSetting := admin.QorWidgetSetting{}
	slideshowSetting.Name = "home page banner"
	slideshowSetting.GroupName = "Banner"
	slideshowSetting.WidgetType = "SlideShow"
	slideshowSetting.Scope = "default"
	slideshowValue := &struct {
		SlideImages []slideImage
	}{}

	for _, s := range Seeds.Slides {
		slide := slideImage{Title: s.Title, SubTitle: s.SubTitle, Button: s.Button, Link: s.Link}
		if file, err := openFileByURL(s.Image); err == nil {
			defer file.Close()
			slide.Image.Scan(file)
		} else {
			fmt.Printf("open file (%q) failure, got err %v", "banner", err)
		}
		slideshowValue.SlideImages = append(slideshowValue.SlideImages, slide)
	}
	slideshowSetting.SetSerializableArgumentValue(slideshowValue)
	if err := DraftDB.Create(&slideshowSetting).Error; err != nil {
		fmt.Printf("Save widget (%v) failure, got err %v", slideshowSetting, err)
	}

	// Featured Products
	featureProducts := admin.QorWidgetSetting{}
	featureProducts.Name = "featured products"
	featureProducts.Description = "featured product list"
	featureProducts.WidgetType = "Products"
	featureProducts.SetSerializableArgumentValue(&struct {
		Products       []string
		ProductsSorter sorting.SortableCollection
	}{
		Products:       []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		ProductsSorter: sorting.SortableCollection{PrimaryKeys: []string{"1", "2", "3", "4", "5", "6", "7", "8"}},
	})
	if err := DraftDB.Create(&featureProducts).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", featureProducts, err)
	}

	// Banner edit items
	for _, s := range Seeds.BannerEditorSettings {
		setting := banner_editor.QorBannerEditorSetting{}
		id, _ := strconv.Atoi(s.ID)
		setting.ID = uint(id)
		setting.Kind = s.Kind
		setting.Value.SerializedValue = s.Value
		if err := DraftDB.Create(&setting).Error; err != nil {
			log.Fatalf("Save QorBannerEditorSetting (%v) failure, got err %v", setting, err)
		}
	}

	// Men collection
	menCollectionWidget := admin.QorWidgetSetting{}
	menCollectionWidget.Name = "men collection"
	menCollectionWidget.Description = "Men collection baner"
	menCollectionWidget.WidgetType = "FullWidthBannerEditor"
	menCollectionWidget.Value.SerializedValue = `{"Value":"%3Cdiv%20class%3D%22qor-bannereditor__html%22%20style%3D%22position%3A%20relative%3B%20height%3A%20100%25%3B%22%20data-image-width%3D%221280%22%20data-image-height%3D%22480%22%3E%3Cspan%20class%3D%22qor-bannereditor-image%22%3E%3Cimg%20src%3D%22%2Fsystem%2Fmedia_libraries%2F1%2Ffile.jpg%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2212%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2010.0781%25%3B%20top%3A%2018.125%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22129%22%20data-position-top%3D%2287%22%3E%3Ch1%20class%3D%22banner-title%22%20style%3D%22color%3A%20%3B%22%3EMEN%20COLLECTION%3C%2Fh1%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2210%22%20style%3D%22position%3A%20absolute%3B%20left%3A%209.92188%25%3B%20top%3A%2029.7917%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22127%22%20data-position-top%3D%22143%22%3E%3Ch2%20class%3D%22banner-sub-title%22%20style%3D%22color%3A%20%3B%22%3ECheck%20the%20newcomming%20collection%3C%2Fh2%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2211%22%20style%3D%22position%3A%20absolute%3B%20left%3A%209.92188%25%3B%20top%3A%2047.0833%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22127%22%20data-position-top%3D%22226%22%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3Eview%20more%3C%2Fa%3E%3C%2Fspan%3E%3C%2Fdiv%3E"}`
	if err := DraftDB.Create(&menCollectionWidget).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", menCollectionWidget, err)
	}

	// Women collection
	womenCollectionWidget := admin.QorWidgetSetting{}
	womenCollectionWidget.Name = "women collection"
	womenCollectionWidget.Description = "Women collection banner"
	womenCollectionWidget.WidgetType = "FullWidthBannerEditor"
	womenCollectionWidget.Value.SerializedValue = `{"Value":"%3Cdiv%20class%3D%22qor-bannereditor__html%22%20style%3D%22position%3A%20relative%3B%20height%3A%20100%25%3B%22%20data-image-width%3D%221280%22%20data-image-height%3D%22480%22%3E%3Cspan%20class%3D%22qor-bannereditor-image%22%3E%3Cimg%20src%3D%22%2Fsystem%2Fmedia_libraries%2F2%2Ffile.jpg%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2223%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2010.0781%25%3B%20top%3A%2018.125%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22129%22%20data-position-top%3D%2287%22%3E%3Ch1%20class%3D%22banner-title%22%20style%3D%22color%3A%20%3B%22%3EWOMEN%20COLLECTION%3C%2Fh1%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2221%22%20style%3D%22position%3A%20absolute%3B%20left%3A%209.92188%25%3B%20top%3A%2029.7917%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22127%22%20data-position-top%3D%22143%22%3E%3Ch2%20class%3D%22banner-sub-title%22%20style%3D%22color%3A%20%3B%22%3ECheck%20the%20newcomming%20collection%3C%2Fh2%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2222%22%20style%3D%22position%3A%20absolute%3B%20left%3A%209.92188%25%3B%20top%3A%2047.0833%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22127%22%20data-position-top%3D%22226%22%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3Eview%20more%3C%2Fa%3E%3C%2Fspan%3E%3C%2Fdiv%3E"}`
	if err := DraftDB.Create(&womenCollectionWidget).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", womenCollectionWidget, err)
	}

	// New arrivals promotio
	newArrivalsCollectionWidget := admin.QorWidgetSetting{}
	newArrivalsCollectionWidget.Name = "new arrivals promotion"
	newArrivalsCollectionWidget.Description = "New arrivals promotion banner"
	newArrivalsCollectionWidget.WidgetType = "FullWidthBannerEditor"
	newArrivalsCollectionWidget.Value.SerializedValue = `{"Value":"%3Cdiv%20class%3D%22qor-bannereditor__html%22%20style%3D%22position%3A%20relative%3B%20height%3A%20100%25%3B%22%20data-image-width%3D%221172%22%20data-image-height%3D%22300%22%3E%3Cspan%20class%3D%22qor-bannereditor-image%22%3E%3Cimg%20src%3D%22%2Fsystem%2Fmedia_libraries%2F3%2Ffile.jpg%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2233%22%20style%3D%22position%3A%20absolute%3B%20left%3A%209.47099%25%3B%20top%3A%2030%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22111%22%20data-position-top%3D%2290%22%3E%3Ch1%20class%3D%22banner-title%22%20style%3D%22color%3A%20%3B%22%3ENew%20Arrivals%3C%2Fh1%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2231%22%20style%3D%22position%3A%20absolute%3B%20left%3A%208.61775%25%3B%20top%3A%20auto%3B%20right%3A%20auto%3B%20bottom%3A%2030.6667%25%3B%22%20data-position-left%3D%22101%22%20data-position-top%3D%22173%22%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3ESHOP%20COLLECTION%3C%2Fa%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2232%22%20style%3D%22position%3A%20absolute%3B%20left%3A%209.55631%25%3B%20top%3A%2016%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22112%22%20data-position-top%3D%2248%22%3E%3Cp%20class%3D%22banner-text%22%20style%3D%22color%3A%20%3B%22%3ETHE%20STYLE%20THAT%20FITS%20EVERYTHING%3C%2Fp%3E%3C%2Fspan%3E%3C%2Fdiv%3E"}`
	if err := DraftDB.Create(&newArrivalsCollectionWidget).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", newArrivalsCollectionWidget, err)
	}

	// Model products
	modelCollectionWidget := admin.QorWidgetSetting{}
	modelCollectionWidget.Name = "model products"
	modelCollectionWidget.Description = "Model products banner"
	modelCollectionWidget.WidgetType = "FullWidthBannerEditor"
	modelCollectionWidget.Value.SerializedValue = `{"Value":"%3Cdiv%20class%3D%22qor-bannereditor__html%22%20style%3D%22position%3A%20relative%3B%20height%3A%20100%25%3B%22%20data-image-width%3D%221100%22%20data-image-height%3D%221200%22%3E%3Cspan%20class%3D%22qor-bannereditor-image%22%3E%3Cimg%20src%3D%22%2Fsystem%2Fmedia_libraries%2F4%2Ffile.jpg%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2249%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2026.4545%25%3B%20top%3A%204.41667%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22291%22%20data-position-top%3D%2253%22%3E%3Ch1%20class%3D%22banner-title%22%20style%3D%22color%3A%20%3B%22%3EENJOY%20THE%20NEW%20FASHION%20EXPERIENCE%3C%2Fh1%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2242%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2043.2727%25%3B%20top%3A%208.41667%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22476%22%20data-position-top%3D%22101%22%3E%3Cp%20class%3D%22banner-text%22%20style%3D%22color%3A%20%3B%22%3ENew%20look%20of%202017%3C%2Fp%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2243%22%20style%3D%22position%3A%20absolute%3B%20left%3A%205.45455%25%3B%20top%3A%2044.25%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%2260%22%20data-position-top%3D%22531%22%3E%3Cdiv%20class%3D%22model-buy-block%22%3E%3Ch2%20class%3D%22banner-sub-title%22%3ETOP%3C%2Fh2%3E%3Cp%20class%3D%22banner-text%22%3E%2429.99%3C%2Fp%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3EVIEW%20DETAILS%3C%2Fa%3E%3C%2Fdiv%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2244%22%20style%3D%22position%3A%20absolute%3B%20left%3A%20auto%3B%20top%3A%2050.8333%25%3B%20right%3A%209.58527%25%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22841%22%20data-position-top%3D%22610%22%3E%3Cdiv%20class%3D%22model-buy-block%22%3E%3Ch2%20class%3D%22banner-sub-title%22%3EPINK%20JACKET%3C%2Fh2%3E%3Cp%20class%3D%22banner-text%22%3E%2469.99%3C%2Fp%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3EVIEW%20DETAILS%3C%2Fa%3E%3C%2Fdiv%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2247%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2012.3636%25%3B%20top%3A%20auto%3B%20right%3A%20auto%3B%20bottom%3A%2014.2032%25%3B%22%20data-position-left%3D%22136%22%20data-position-top%3D%22903%22%3E%3Cdiv%20class%3D%22model-buy-block%22%3E%3Ch2%20class%3D%22banner-sub-title%22%3EBOTTOM%3C%2Fh2%3E%3Cp%20class%3D%22banner-text%22%3E%2432.99%3C%2Fp%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3EVIEW%20DETAILS%3C%2Fa%3E%3C%2Fdiv%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2245%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2053.2727%25%3B%20top%3A%2048.5%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22586%22%20data-position-top%3D%22582%22%3E%3Cimg%20src%3D%22%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Farrow-left.png%22%20class%3D%22banner-image%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2246%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2015.5455%25%3B%20top%3A%2043.0833%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22171%22%20data-position-top%3D%22517%22%3E%3Cimg%20src%3D%22%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Farrow-right.png%22%20class%3D%22banner-image%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2248%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2019.2727%25%3B%20top%3A%20auto%3B%20right%3A%20auto%3B%20bottom%3A%2024.8333%25%3B%22%20data-position-left%3D%22212%22%20data-position-top%3D%22879%22%3E%3Cimg%20src%3D%22%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Farrow-right.png%22%20class%3D%22banner-image%22%3E%3C%2Fspan%3E%3C%2Fdiv%3E"}`
	if err := DraftDB.Create(&modelCollectionWidget).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", modelCollectionWidget, err)
	}
}

func main() {
	createRecords()
}

var Seeds = struct {
	Setting struct {
		ShippingFee     uint
		GiftWrappingFee uint
		CODFee          uint `gorm:"column:cod_fee"`
		TaxRate         int
		Address         string
		City            string
		Region          string
		Country         string
		Zip             string
		Latitude        float64
		Longitude       float64
	}
	Seo struct {
		SiteName    string
		DefaultPage struct {
			Title       string
			Description string
			Keywords    string
		}
		HomePage struct {
			Title       string
			Description string
			Keywords    string
		}
		ProductPage struct {
			Title       string
			Description string
			Keywords    string
		}
	}
	Slides []struct {
		Title    string
		SubTitle string
		Button   string
		Link     string
		Image    string
	}
	MediaLibraries []struct {
		Title string
		Image string
	}
	BannerEditorSettings []struct {
		ID    string
		Kind  string
		Value string
	}
}{}

func init() {
	migrations.AutoMigrate()
	filepaths, _ := filepath.Glob(filepath.Join("config", "database", "seeds", "data", "*.yml"))
	if err := configor.Load(&Seeds, filepaths...); err != nil {
		panic(err)
	}
}

func createRecords() {
	createAdminUsers()
	fmt.Println("--> Created admin users.")
	createSeo()
	fmt.Println("--> Created seo.")
	createWidgets()
	fmt.Println("--> Created widgets")
}

func openFileByURL(rawURL string) (*os.File, error) {
	if fileURL, err := url.Parse(rawURL); err != nil {
		return nil, err
	} else {
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]

		filePath := filepath.Join(os.TempDir(), fileName)

		if _, err := os.Stat(filePath); err == nil {
			return os.Open(filePath)
		}

		file, err := os.Create(filePath)
		if err != nil {
			return file, err
		}

		check := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		resp, err := check.Get(rawURL) // add a filter to check redirect
		if err != nil {
			return file, err
		}
		defer resp.Body.Close()
		fmt.Printf("----> Downloaded %v\n", rawURL)

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return file, err
		}
		return file, nil
	}
}
