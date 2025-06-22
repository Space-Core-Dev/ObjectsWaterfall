package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"object-shooter.com/data"
	"object-shooter.com/handlers"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		panic(err)
	}
	err = data.InitDbConnection()
	if err != nil {
		panic(err)
	}
	engine := gin.Default()

	engine.POST("/start", handlers.Start)
	engine.GET("/stop", handlers.Stop)
	engine.POST("/seed", handlers.Seed)

	engine.Run(":8888")
}
