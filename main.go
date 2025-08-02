package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"objectswaterfall.com/application/handlers"
	"objectswaterfall.com/data"
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

	engine.POST("/add", handlers.Add)
	engine.POST("/start", handlers.Start)
	engine.GET("/stop", handlers.Stop)
	engine.POST("/seed", handlers.Seed)
	engine.GET("/getWorkers", handlers.GetWorkers)

	engine.Run(":8888")
}
