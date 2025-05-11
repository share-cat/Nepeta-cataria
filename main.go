package main

import (
	"log"
	"share-cat/Nepeta-cataria/boot"
	"share-cat/Nepeta-cataria/config"
	"share-cat/Nepeta-cataria/global"
	"share-cat/Nepeta-cataria/models"
	"share-cat/Nepeta-cataria/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	initdb()
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	server.Use(gin.Logger())
	server.Use(cors.New(corsConfig))

	app := boot.NewApplication(global.DB)
	routes.SetupRouter(server, global.DB, app.Controllers)

	server.Run()

}

func initdb() {
	var err error
	DBsetting := config.GetConfig()
	global.DB, err = boot.NewDBEngine(&DBsetting.DatabaseSetting)

	if err != nil {
		log.Fatal("failed to connect database")
		log.Fatal(err)
	}

	dbInit()
}

func dbInit() {
	err := global.DB.AutoMigrate(
		&models.Achievement{},
		&models.ChatSession{},
		&models.Discussion{},
		&models.Judge{},
		&models.Permission{},
		&models.School{},
		&models.Social{},
		&models.Solution{},
		&models.SubmitEffect{},
		&models.SupportLang{},
		&models.SystemConfig{},
		&models.User{},
		&models.UserProfile{},
	)
	if err != nil {
		log.Fatal(err)
	}

}
