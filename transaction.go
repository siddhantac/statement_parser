package main

import (
	"strconv"
	"strings"
)

type transaction struct {
	Amount   string
	Name     string
	Day      string
	Month    string
	Category string
}

func NewTransaction(name, amount, day, month string) *transaction {
	t := &transaction{
		Name:   name,
		Amount: amount,
		Day:    day,
		Month:  month,
	}

	t.Category = guessCategory(t.Name)
	return t
}

func guessCategory(name string) string {
	name = strings.ToLower(name)

	categories := map[string][]string{
		"Groceries":     {"fairprice", "cold storage", "perk coffee"},
		"Transport":     {"cabcharge asia", "bus/mrt", "transit", "grab* dr"},
		"Meals":         {"subway", "deliveroo", "grab* r", "starbucks"},
		"Personal":      {"lazada", "fitness first", "steam"},
		"Entertainment": {"netflix", "youtube", "spotify"},
		"Utilities":     {"liberty wireless", "sp digital", "myrepublic", "gomo mobile plan"},
		"Insurance":     {"prudential"},
	}

	for k, v := range categories {
		for _, option := range v {
			if strings.Contains(name, option) {
				return k
			}
		}
	}

	return ""
}

type ByDate []*transaction

func (d ByDate) Len() int      { return len(d) }
func (d ByDate) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d ByDate) Less(i, j int) bool {
	var iMonthIdx, jMonthIdx int
	for k, val := range months {
		if d[i].Month == val {
			iMonthIdx = k
		}
		if d[j].Month == val {
			jMonthIdx = k
		}
	}

	if iMonthIdx < jMonthIdx {
		return true
	}

	if iMonthIdx > jMonthIdx {
		return false
	}

	iDay, err := strconv.Atoi(d[i].Day)
	if err != nil {
		panic(err)
	}

	jDay, err := strconv.Atoi(d[j].Day)
	if err != nil {
		panic(err)
	}

	return iDay < jDay
}

var months = []string{
	"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC",
}
