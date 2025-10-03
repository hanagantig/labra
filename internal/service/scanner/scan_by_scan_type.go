package scanner

import (
	"context"
	"fmt"
)

func (s *Service) ScanByScanType(ctx context.Context, sType scannerType, bytes []byte) (string, error) {
	scanner,ok := s.scanners[sType]
	if !ok {
		return "", fmt.Errorf("scanner not found for type %v", sType)
	}

	return scanner.ScanBytes(ctx, bytes)
}
