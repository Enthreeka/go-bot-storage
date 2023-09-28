package model

import (
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

func FindIdName(data string) (int, string) {
	parts := strings.Split(data, "_")
	if len(parts) < 2 {
		return 0, ""
	}

	cellID, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, ""
	}

	return cellID, parts[1]
}
