package excelPkg

type HeaderInfo struct {
	En string `json:"en"`
	Ar string `json:"ar"`
}
type Headers []HeaderInfo
type Row []any
type Data []Row
type ChunkOfData []Data
