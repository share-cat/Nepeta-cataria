package main

import (
	"log"
	"share-cat/Nepeta-cataria/database"
	"share-cat/Nepeta-cataria/global"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	server.Use(gin.Logger())
	server.Use(cors.New(corsConfig))

	server.Run()

}

func dbInit() {
	err := global.DB.AutoMigrate(
		&database.Judge{},
		&database.Permission{},
		&database.School{},
		&database.Social{},
		&database.Solution{},
		&database.SupportLang{},
		&database.SystemConfig{},
		&database.User{},
		&database.UserProfile{},
		&database.ChatSession{},
	)
	if err != nil {
		log.Fatal(err)
	}

}
