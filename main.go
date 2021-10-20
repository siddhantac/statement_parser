package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

var pattern *regexp.Regexp

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var infile, outfile string
	var parseRefund bool

	flag.StringVar(&infile, "in", "", "input file")
	flag.StringVar(&outfile, "out", "", "output file")
	flag.BoolVar(&parseRefund, "refund", false, "parse refunds")
	flag.Parse()

	file, err := os.Open(infile)
	if err != nil {
		return err
	}

	pattern = regexp.MustCompile(`^(?P<day>[0123][0-9]) (?P<month>[A-Z]+) (?P<name>.*) (?P<amount>[0-9\.]+)$`)

	if parseRefund {
		pattern = regexp.MustCompile(`^(?P<day>[0123][0-9]) (?P<month>[A-Z]+) (?P<name>.*) \((?P<amount>[0-9\.]+)\)$`)
	}

	scanner := bufio.NewScanner(file)

	allTransactions := make([]*transaction, 0)
	for scanner.Scan() {
		match := pattern.FindStringSubmatch(scanner.Text())

		if len(match) == 0 {
			continue
		}

		allTransactions = append(allTransactions, NewTransaction(match[3], match[4], match[1], match[2]))
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	}

	// sort.Sort(ByDate(allTransactions))

	out, err := os.OpenFile(outfile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, tr := range allTransactions {
		data := fmt.Sprintf("%s,%s-%s,%s,%s\n", tr.Name, tr.Day, tr.Month, tr.Amount, tr.Category)
		out.Write([]byte(data))
	}

	return nil
}

/*
func process(pattern *regexp.Regexp) {
	file, err := os.Open("citi4.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		match := pattern.FindStringSubmatch(scanner.Text())

		if len(match) == 0 {
			continue
		}

		data := fmt.Sprintf("%s,%s,%s\n", match[1], match[2], match[3])
		out.Write([]byte(data))

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	}
}
*/
