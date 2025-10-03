package scanner

import (
	"context"
	"fmt"
)

func (s *Service) ScanByFileType(ctx context.Context, fType string, bytes []byte) (string, error) {
	scanner, ok := s.scanners[supportedFileTypes[fType]]
	if !ok {
		return "", fmt.Errorf("scanner not found for file type %v", fType)
	}

	return scanner.ScanBytes(ctx, bytes)
}
