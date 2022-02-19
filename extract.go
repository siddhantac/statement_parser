package main

import (
	"io/ioutil"
	"os"

	"github.com/ledongthuc/pdf"
	//"github.com/dcu/pdf"
)

func extract() {
	content, err := readPdf(os.Args[1]) // Read local pdf file
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("rows.txt", []byte(content), 0777)
	return
}

// WORKS!!!
// I can parse the result from this function.
func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()
	var output string

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			//println(">>>> row: ", row.Position)
			for _, word := range row.Content {
				//fmt.Println(word.S)
				output += word.S
			}
			output += "\n"
		}
	}
	return output, nil
}
