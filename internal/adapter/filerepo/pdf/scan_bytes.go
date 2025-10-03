package pdf

import (
	"bytes"
	"context"
	"fmt"
	"github.com/martoche/pdf"
)

func (r Repository) ScanBytes(ctx context.Context, file []byte) (string, error) {
	rd, err := pdf.NewReader(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		return "", err
	}

	p, err := rd.GetPlainText()
	if err != nil {
		return "", err
	}

	buf, ok := p.(*bytes.Buffer)
	if !ok {
		return "", fmt.Errorf("the library no longer uses bytes.Buffer to implement io.Reader")
	}

	return buf.String(), nil
}
