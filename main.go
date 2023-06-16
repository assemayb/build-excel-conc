package main

import (
	"fmt"
	"sync"

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
		randomString := fmt.Sprintf("randomString%d", i)
		randomNumber := fmt.Sprintf("%d", i)
		row := Row{randomString, randomNumber}
		data = append(data, row)
	}
	return data
}

func main() {
	data := simulateLargeData(500)
	file := excelize.NewFile()
	sheetName := "Sheet1"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}
	file.SetActiveSheet(index)
	defer file.Close()

	lang := "ar"
	headers := []HeaderInfo{
		{en: "Name", ar: "الاسم"},
		{en: "Age", ar: "العمر"},
	}

	addExcelFileHeaders(headers, file, sheetName, lang)
	chunkSize := len(data) / numWorkers
	listOfRowsChunks := make(ChunkOfData, numWorkers)
	listOfRowsChunks = chunkIncomingData(chunkSize, data, listOfRowsChunks)

	var wg sync.WaitGroup
	for chunkIndex, chunkOfRows := range listOfRowsChunks {
		wg.Add(1)
		go processChunk(chunkOfRows, &wg, chunkIndex, chunkSize, file, sheetName)
	}
	wg.Wait()

	err = file.SaveAs("test.xlsx")
	if err != nil {
		panic(err)
	}
}

func processChunk(dataChunk Data, wg *sync.WaitGroup, chunkIndex int, chunkSize int, file *excelize.File, sheetName string) {
	for idx, row := range dataChunk {
		rowIndex := chunkIndex*chunkSize + idx
		for i, cellValue := range row {
			cell := fmt.Sprintf("%s%d", string('A'+i), rowIndex)
			file.SetCellValue(sheetName, cell, cellValue)
		}
		idx++
	}
	wg.Done()
}

func chunkIncomingData(chunkSize int, data Data, chunks []Data) ChunkOfData {
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers {
			end = len(data)
		}
		chunks[i] = data[start:end]
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
