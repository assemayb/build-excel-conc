package main

import (
	"fmt"
	"sync"

	"github.com/tealeg/xlsx"
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
		data = append(data, []interface{}{"John", "30"})
	}
	return data
}

func main() {
	data := simulateLargeData(50_000)

	chunkSize := len(data) / numWorkers
	var chunks = make([]Data, numWorkers)
	chunks = chunkData(chunkSize, data, chunks)

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	headers := []HeaderInfo{
		{en: "Name", ar: "الاسم"},
		{en: "Age", ar: "العمر"},
	}
	addExcelFileHeaders(headers, sheet, "en")

	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go processChunk(chunk, sheet, &wg)
	}

	wg.Wait()
	fmt.Println("ASdsadsad")
	err := file.Save("example.xlsx")
	if err != nil {
		fmt.Println("Error saving file:", err)
	}
}

func processChunk(data Data, sheet *xlsx.Sheet, wg *sync.WaitGroup) {
	for _, item := range data {
		newRow := sheet.AddRow()
		newRow.Sheet.SetColWidth(0, len(data), 25)
		for _, value := range item {
			cell := newRow.AddCell()
			cell.SetValue(value)
		}
	}
	wg.Done()
}

func chunkData(chunkSize int, data Data, chunks []Data) []Data {
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

func addExcelFileHeaders(headers []HeaderInfo, sheet *xlsx.Sheet, lang string) {
	headerRow := sheet.AddRow()
	for _, header := range headers {
		cell := headerRow.AddCell()
		if lang == "en" {
			cell.Value = header.en
		} else {
			cell.Value = header.ar
		}
	}
}
