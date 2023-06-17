package controller

import (
	excelPkg "excel-builder-conc/excel-pkg"
)

type RequestBody struct {
	Headers   []excelPkg.HeaderInfo `json:"headers"`
	Data      excelPkg.Data         `json:"data"`
	Lang      string                `json:"lang"`
	SheetName string                `json:"sheetName"`
}
