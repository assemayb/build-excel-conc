package excelPkg

import (
	"fmt"
	"sync"

	"github.com/xuri/excelize/v2"
)

var (
	numWorkers = 5
)

func BuildExcelFile(data Data, headers []HeaderInfo, lang string, sheetName string) *excelize.File {
	file := excelize.NewFile()
	index, err := file.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}

	file.SetActiveSheet(index)
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

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

	return file
}

func processChunk(dataChunk Data, wg *sync.WaitGroup, chunkIndex int, chunkSize int, file *excelize.File, sheetName string) {
	for idx, row := range dataChunk {
		rowIndex := chunkIndex*chunkSize + idx
		for i, cellValue := range row {
			cell := fmt.Sprintf("%s%d", string('A'+i), rowIndex+2)
			file.SetColWidth(sheetName, cell, cell, 50)
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
	var rightToLeft bool
	if lang == "ar" {
		rightToLeft = true
	}
	if rightToLeft {
		file.SetSheetView(sheetName, 0, &excelize.ViewOptions{
			RightToLeft: &rightToLeft,
		})
	}

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
