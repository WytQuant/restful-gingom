package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"restful-gingorm/config"
	"restful-gingorm/migratinos"
	"restful-gingorm/routes"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalln("Failed to load .env file")
	//}

	config.InitDB()
	defer config.CloseDB()
	migratinos.Migrate()
	//seed.Load() //run once

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.New(corsConfig))
	r.Static("/uploads", "./uploads")

	uploadDirs := [...]string{"articles", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	routes.Serve(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalln("error:", err.Error())
	}
}
