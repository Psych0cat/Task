package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	utils "testTask"
	"testTask/businessLogic"
)

func main() {
	var c *gin.Context
	go businessLogic.UpdateDB(c)
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
	//default logger
	gin.DefaultWriter = io.MultiWriter(logfile, os.Stdout)
	Router.POST("/transaction", businessLogic.MakeTransaction)
	Router.GET("/checkbalance", businessLogic.CheckBalance)
	Router.GET("/nullbalance", businessLogic.NullBalance)
	Router.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"Healthckeck": "OK"}) })
	Router.Use(utils.RequestLogger())
	return Router
}
