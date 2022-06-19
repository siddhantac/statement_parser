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
	// content, err := readByColumn(os.Args[1])
	// if err != nil {
	// 	panic(err)
	// }
	content, err := readPdf(os.Args[1]) // Read local pdf file
	// content, err := readPlainText(os.Args[1]) // Read local pdf file
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("unparsed.txt", []byte(content), 0777)
	// content, err := ioutil.ReadFile("unparsed.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// output := parsePDF3(content)
	// ioutil.WriteFile("output3.txt", []byte(output), 0777)

	// allTransactions := parsePDF3(string(content))
	// out, err := os.OpenFile("output5.txt", os.O_CREATE|os.O_RDWR, 0755)
	// if err != nil {
	// 	panic(err)
	// }
	// defer out.Close()

	// sort.Sort(ByDate(allTransactions))

	// sum := 0.0
	// for _, tr := range allTransactions {
	// 	amt, err := strconv.ParseFloat(tr.Amount, 32)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	sum += amt
	// 	data := fmt.Sprintf("%s,%s-%s,%s,%s\n", tr.Name, tr.Day, tr.Month, tr.Amount, tr.Category)
	// 	out.Write([]byte(data))
	// }
	// fmt.Println("Total:", sum)

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

const (
	lineBegin      = `^`
	lineEnd        = `$`
	searchPattern1 = `(?P<day>\d\d)(?P<month>JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)`
	// searchPattern2 = `(?P<name>.*)(SingaporeSG|SINGAPORESG)(?P<amount>\d+\.\d+)`
	searchPattern2 = `(?P<name>.*)[A-Z][A-Z](?P<amount>\d+\.\d+)`
	searchPattern3 = searchPattern1 + searchPattern2
	refundPattern  = searchPattern1 + `(?P<name>.*)[A-Z][A-Z]\((?P<amount>\d+\.\d+)\)`
)

// const searchPattern = `^(?P<day>\d\d)(?P<month>JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)(?P<name>.*)(SingaporeSG|SINGAPORESG)(?P<amount>\d+\.\d+)$`

func parsePDF(content string) string {
	// mc := regexp.MustCompile(`^(?P<day>[0123][0-9])(?P<month>[A-Z][A-Z][A-Z])(?P<name>.*)(?P<amount>[0-9\.]+)$`)
	mc := regexp.MustCompile(searchPattern3)

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

// const searchPattern1 = `^(?P<day>\d\d)(?P<month>JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)$`
// const searchPattern2 = `^(?P<name>.*)(SingaporeSG|SINGAPORESG)(?P<amount>\d+\.\d+)$`

func parsePDF2(content string) string {
	mc1 := regexp.MustCompile(searchPattern1)
	mc2 := regexp.MustCompile(searchPattern2)

	r := strings.NewReader(content)
	scanner := bufio.NewScanner(r)
	var output string

	// allTransactions := make([]*transaction, 0)
	for scanner.Scan() {
		match1 := mc1.FindStringSubmatch(scanner.Text())

		if len(match1) == 0 {
			continue
		}

		scanner.Scan()
		match2 := mc2.FindStringSubmatch(scanner.Text())
		if len(match2) == 0 {
			continue
		}
		// allTransactions = append(allTransactions, NewTransaction(match[3], match[4], match[1], match[2]))
		// if err := scanner.Err(); err != nil {
		// 	fmt.Println(err)
		// }
		output = fmt.Sprintf("%s%s,%s\n", output, match1[1:], match2[1:])
	}

	return output
}

func parsePDF3(content string) []*transaction {
	mc1 := regexp.MustCompile(lineBegin + searchPattern1 + lineEnd)
	mc2 := regexp.MustCompile(lineBegin + searchPattern2 + lineEnd)
	// mc3 := regexp.MustCompile(lineBegin + searchPattern3 + lineEnd)
	// mc1 := regexp.MustCompile(searchPattern1)
	// mc2 := regexp.MustCompile(searchPattern2)
	mc3 := regexp.MustCompile(searchPattern3 + `.*$`)
	mcRefund := regexp.MustCompile(refundPattern + `.*$`)

	r := strings.NewReader(content)
	scanner := bufio.NewScanner(r)
	// var output string

	allTransactions := make([]*transaction, 0)
	for scanner.Scan() {
		matchRefund := mcRefund.FindStringSubmatch(scanner.Text())
		if len(matchRefund) != 0 {
			allTransactions = append(allTransactions, NewTransaction(matchRefund[3], "-"+matchRefund[4], matchRefund[1], matchRefund[2]))
			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			}
			continue
		}

		match1 := mc1.FindStringSubmatch(scanner.Text())

		if len(match1) == 0 {
			match3 := mc3.FindStringSubmatch(scanner.Text())
			if len(match3) != 0 {
				// output = fmt.Sprintf("%s%s,%s\n", output, match3[1:], match3[1:])

				fmt.Println(match3[1:])
				allTransactions = append(allTransactions, NewTransaction(match3[3], match3[4], match3[1], match3[2]))
				if err := scanner.Err(); err != nil {
					fmt.Println(err)
				}
			}
			continue
		}

		scanner.Scan()
		match2 := mc2.FindStringSubmatch(scanner.Text())
		if len(match2) == 0 {
			continue
		}

		fmt.Printf("> %s,%s\n", match1[1:], match2[1:])
		allTransactions = append(allTransactions, NewTransaction(match2[1], match2[2], match1[1], match1[2]))
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
		// output = fmt.Sprintf("%s%s,%s\n", output, match1[1:], match2[1:])
	}

	return allTransactions
}
