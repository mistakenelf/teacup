package formatter

import (
	"bytes"

	"github.com/ledongthuc/pdf"
)

// ReadPdf reads a PDF file given a name.
func ReadPdf(name string) (string, error) {
	f, r, err := pdf.Open(name)
	if err != nil {
		return "", err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	b, err := r.GetPlainText()

	if err != nil {
		return "", err
	}

	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
