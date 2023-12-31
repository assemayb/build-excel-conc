package controller

import (
	excelPkg "excel-builder-conc/excel-builder"
)

type RequestBody struct {
	Headers   excelPkg.Headers `json:"headers"`
	Data      excelPkg.Data    `json:"data"`
	Lang      string           `json:"lang"`
	SheetName string           `json:"sheetName"`
}
