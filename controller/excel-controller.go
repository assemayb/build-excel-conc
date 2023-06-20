package controller

import (
	"bytes"
	excelPkg "excel-builder-conc/excel-builder"
	"fmt"
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

	var fileBuffer = new(bytes.Buffer)
	err = file.Write(fileBuffer)
	if err != nil {
		fmt.Println("Error Writing file data to a buffer", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		panic(err)
	}
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBuffer.Bytes())
}
