package controller

import (
	excelPkg "excel-builder-conc/excel-pkg"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateExcelFile(ctx *gin.Context) {
	var body RequestBody
	err := ctx.BindJSON(&body)
	if err != nil {
		log.Println(err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	file := excelPkg.BuildExcelFile(body.Data, body.Headers, body.Lang, body.SheetName)
	err = file.SaveAs("test.xlsx")

	if err != nil {
		log.Println(err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "file Created Successfully"})
}
