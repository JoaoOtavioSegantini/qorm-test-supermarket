package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joaootav/system_supermarket/config"
	"github.com/joaootav/system_supermarket/config/routes"
	"github.com/joaootav/system_supermarket/database"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {

	gin.SetMode(gin.ReleaseMode)

	database.Connect(os.Getenv("DATABASE_URL"))
	database.Migrate()
	router := gin.Default()
	router.Static("/public", "./public")
	router.StaticFile("favicon.ico", "./public/favicon.ico")

	router.Any("/admin/*resources", gin.WrapH(database.Mux))
	router.Use(gin.WrapH(routes.Router()))

	log.Printf("Listening on: %v\n", config.Config.Port)
	router.Run(fmt.Sprintf(":%d", config.Config.Port))
}
