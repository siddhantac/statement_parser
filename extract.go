package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
	//"github.com/dcu/pdf"
)

func extract() {
	content, err := readPdf(os.Args[1]) // Read local pdf file
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("unparsed.txt", []byte(content), 0777)

	output := parsePDF(content)
	ioutil.WriteFile("output.txt", []byte(output), 0777)
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

const searchPattern = `^(?P<day>\d\d)(?P<month>JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)(?P<name>.*)(SingaporeSG|SINGAPORESG)(?P<amount>\d+\.\d+)$`

func parsePDF(content string) string {
	// mc := regexp.MustCompile(`^(?P<day>[0123][0-9])(?P<month>[A-Z][A-Z][A-Z])(?P<name>.*)(?P<amount>[0-9\.]+)$`)
	mc := regexp.MustCompile(searchPattern)

	r := strings.NewReader(content)
	scanner := bufio.NewScanner(r)
	var output string

	// allTransactions := make([]*transaction, 0)
	for scanner.Scan() {
		match := mc.FindStringSubmatch(scanner.Text())

		if len(match) == 0 {
			continue
		}

		// allTransactions = append(allTransactions, NewTransaction(match[3], match[4], match[1], match[2]))
		// if err := scanner.Err(); err != nil {
		// 	fmt.Println(err)
		// }
		output = fmt.Sprintf("%s%s\n", output, match[1:])
	}

	return output
}
