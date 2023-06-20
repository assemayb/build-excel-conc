package main

import (
	controller "excel-builder-conc/controller"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
)

func init() {
	server = gin.New()
	server.Use(gin.Logger(), gin.Recovery())

}

func main() {
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	server.POST("/api/excel/build", controller.CreateExcelFile)
	log.Fatal(server.Run(":9007"))

	//  TODO: change this in production
	// gin.SetMode(gin.TestMode)
}
