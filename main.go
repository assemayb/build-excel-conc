package main

import (
	"fmt"
	"sync"

	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
)

var (
	numWorkers = 5
)

type HeaderInfo struct {
	en string `json:"en"`
	ar string `json:"ar"`
}

type Data [][]interface{}

func simulateLargeData(num int) Data {
	var data Data
	for i := 0; i < num; i++ {
		randomString := fmt.Sprintf("Random String %d", i)
		randomNumber := fmt.Sprintf("%d", i)
		data = append(data, []interface{}{randomString, randomNumber})
	}
	return data
}

func main() {
	data := simulateLargeData(50_000)

	file := excelize.NewFile()
	sheetName := "Sheet1"
	index, err := file.NewSheet(sheetName)
	file.SetActiveSheet(index)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	headers := []HeaderInfo{
		{
			en: "Name",
			ar: "الاسم",
		},
		{
			en: "Age",
			ar: "العمر",
		},
	}
	addExcelFileHeaders(headers, file, sheetName, "en")

	chunkSize := len(data) / numWorkers
	var chunks = make([]Data, numWorkers)
	chunks = chunkIncomingData(chunkSize, data, chunks)

	// Process each chunk of data in a separate goroutine
	var wg sync.WaitGroup
	for _, chunk := range chunks {

		wg.Add(1)
		go func(data Data) {
			rowIndex, err := file.GetSheetIndex(sheetName)
			if err != nil {
				panic(err)
			}

			for _, item := range data {
				for i, value := range item {
					cell := fmt.Sprintf("%s%d", string('A'+i), rowIndex)
					file.SetCellValue(sheetName, cell, value)
				}
				rowIndex++
			}
			wg.Done()

		}(chunk)

	}

	// Wait for all the goroutines to complete
	wg.Wait()

	fmt.Println("Done processing chunks")
	err = file.SaveAs("test.xlsx")

	if err != nil {
		panic(err)
	}
}

func processChunk(data Data, sheet *xlsx.Sheet, wg *sync.WaitGroup) {
	for _, item := range data {

		newRow := sheet.AddRow()
		error := newRow.Sheet.SetColWidth(0, len(data), 25)
		if error != nil {
			panic(error)
		}

		for _, value := range item {
			cell := newRow.AddCell()
			cell.SetValue(value)
		}

	}
	wg.Done()
}

func chunkIncomingData(chunkSize int, data Data, chunks []Data) []Data {
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == 3 {
			end = len(data)
		}
		chunks[i] = data[start:end]
		chunks = append(chunks, data[start:end])
	}
	return chunks
}

func addExcelFileHeaders(headers []HeaderInfo, file *excelize.File, sheetName string, lang string) {
	for i, header := range headers {
		var headerValue string

		if lang == "ar" {
			headerValue = header.ar
		} else {
			headerValue = header.en
		}

		cell := fmt.Sprintf("%s%d", string('A'+i), 1)
		file.SetCellValue(sheetName, cell, headerValue)
	}
}
