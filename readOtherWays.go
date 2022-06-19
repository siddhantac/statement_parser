package main

import (
	"bytes"

	"github.com/ledongthuc/pdf"
)

func readByColumn(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}

	var output string

	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		cols, err := p.GetTextByColumn()
		if err != nil {
			return "", err
		}

		for _, c := range cols {
			for _, w := range c.Content {
				output += w.S
			}
			output += "\n"
		}
	}

	return output, nil
}

func readPlainText(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
