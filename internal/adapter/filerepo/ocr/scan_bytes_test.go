package ocr

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_ScanBytes(t *testing.T) {
	tdata := map[string]struct {
		lang           string
		fileBytes      []byte
		expectedResult string
		expectedErr    error
	}{
		"first": {
			expectedErr: errors.New("image data cannot be empty"),
		},
	}

	for name, td := range tdata {
		t.Run(name, func(t *testing.T) {
			repo := NewRepository()

			res, err := repo.ScanBytes(context.Background(), td.fileBytes)

			assert.Equal(t, td.expectedResult, res)
			assert.Equal(t, td.expectedErr, err)
		})
	}
}
