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

type Row []string
type Data []Row
type ChunkOfData []Data

func simulateLargeData(num int) Data {
	var data Data
	for i := 0; i < num; i++ {
		randomString := fmt.Sprintf("Random String %d", i)
		randomNumber := fmt.Sprintf("%d", i)
		row := Row{randomString, randomNumber}
		data = append(data, row)
	}
	return data
}

func main() {
	data := simulateLargeData(5)
	file := excelize.NewFile()
	sheetName := "Sheet1"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}
	file.SetActiveSheet(index)
	defer file.Close()
	lang := "en"
	headers := []HeaderInfo{
		{en: "Name", ar: "الاسم"},
		{en: "Age", ar: "العمر"},
	}

	for _, rowsList := range data {
		for row := 0; row < len(rowsList); row++ {
			file.InsertRows(sheetName, row, 1)
		}
	}

	addExcelFileHeaders(headers, file, sheetName, lang)

	chunkSize := len(data) / numWorkers
	var listOfRowsChunks = make(ChunkOfData, numWorkers)

	listOfRowsChunks = chunkIncomingData(chunkSize, data, listOfRowsChunks)
	fmt.Println(listOfRowsChunks)

	var wg sync.WaitGroup

	for chunkIndex, chunkOfRows := range listOfRowsChunks {

		wg.Add(1)
		go func(dataChunk Data) {

			rowIndex := chunkIndex*chunkSize + 1
			for _, row := range dataChunk {
				for i, cellValue := range row {
					cell := fmt.Sprintf("%s%d", string('A'+i), rowIndex)
					file.SetCellValue(sheetName, cell, cellValue)
				}
				rowIndex++
			}
			wg.Done()

		}(chunkOfRows)
	}

	wg.Wait()
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

func chunkIncomingData(chunkSize int, data Data, chunks []Data) ChunkOfData {
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
		headerItem := ""
		if lang == "ar" {
			headerItem = header.ar
		} else {
			headerItem = header.en
		}

		cell := fmt.Sprintf("%s%d", string('A'+i), 1)
		file.SetCellValue(sheetName, cell, headerItem)
	}
}
