package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	utils "testTask"
	"testTask/handlers"
)

func main() {
	var c *gin.Context
	go handlers.UpdateDB(c)
	r := SetupRouter()
	err := r.Run(fmt.Sprintf(":%v", utils.AppPort))
	if err != nil {
		fmt.Println("Unable to start web app")
	}
}

func SetupRouter() *gin.Engine {
	Router := gin.New()
	Router.Use(utils.RequestLogger())
	Router.Use(utils.RespLogger())
	logfile, err := os.Create(utils.LogfileName)
	if err != nil {
		log.Println("Unable to create log file")
	}
	gin.DefaultWriter = io.MultiWriter(logfile, os.Stdout)
	Router.POST("/transaction", handlers.MakeTransaction)
	Router.GET("/checkbalance", handlers.CheckBalance)
	Router.GET("/nullbalance", handlers.NullBalance)
	Router.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"Healthckeck": "OK"}) })
	Router.Use(utils.RequestLogger())
	return Router
}
