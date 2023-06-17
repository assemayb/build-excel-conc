package excelPkg

type HeaderInfo struct {
	en string `json:"en"`
	ar string `json:"ar"`
}

type Row []any
type Data []Row
type ChunkOfData []Data
