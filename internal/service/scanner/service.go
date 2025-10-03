package scanner

import "context"

const (
	PDFType = scannerType(iota + 1)
	OCRType
)

var supportedFileTypes = map[string]scannerType{
	"image/jpg":       OCRType,
	"image/jpeg":      OCRType,
	"image/png":       OCRType,
	"image/pdf":       PDFType,
	"application/pdf": PDFType,
}

type scannerType int

type byteScanner interface {
	ScanBytes(ctx context.Context, bytes []byte) (string, error)
}

type Service struct {
	scanners map[scannerType]byteScanner
}

func NewService(pdfScn byteScanner, ocrScn byteScanner) *Service {
	return &Service{
		scanners: map[scannerType]byteScanner{
			PDFType: pdfScn,
			OCRType: ocrScn,
		},
	}
}
