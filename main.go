package main

import (
	"log"
	"share-cat/Nepeta-cataria/boot"
	"share-cat/Nepeta-cataria/global"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	server.Use(gin.Logger())
	server.Use(cors.New(corsConfig))

	app := boot.NewApplication(global.DB)
	routes.SetupRouter(server, global.DB, app.Controllers)

	server.Run()

}

func dbInit() {
	err := global.DB.AutoMigrate(
		&models.Judge{},
		&models.Permission{},
		&models.School{},
		&models.Social{},
		&models.Solution{},
		&models.SupportLang{},
		&models.SystemConfig{},
		&models.User{},
		&models.UserProfile{},
		&models.ChatSession{},
	)
	if err != nil {
		log.Fatal(err)
	}

}
