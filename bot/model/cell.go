package model

import (
	"regexp"
	"strconv"
	"strings"
)

// Cell - the main thing is the set of sub-cells (in repository under cells)
type Cell struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int64  `json:"user_id"`
}

type UnderCell struct {
	ID     int    `json:"id"`
	CellID int    `json:"cell_id"`
	Name   string `json:"name"`
}

func IsCell(data string) bool {
	if strings.HasPrefix(data, "cell") {
		return true
	}
	return false
}

func IsUnderCell(data string) bool {
	if strings.HasPrefix(data, "underCell") {
		return true
	}
	return false
}

// TODO обработку ошибок
func FindIntStr(data string) (int, string) {
	re := regexp.MustCompile("[0-9]+")
	digits := re.FindAllString(data, -1)
	digitsStr := strings.Join(digits, "")

	name := strings.Split(data, "_")

	digitsInt, _ := strconv.Atoi(digitsStr)

	return digitsInt, name[1]
}
