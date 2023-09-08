package model

import "strings"

// Cell - the main thing is the set of sub-cells (in repository under cells)
type Cell struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int64  `json:"user_id"`
}

// Cells - slice for all getting cells by UserID
//type Cells []Cell

type UnderCell struct {
	ID     int    `json:"id"`
	CellID int    `json:"cell_id"`
	Name   string `json:"name"`
}

// UnderCells - slice for all getting cells by CellID
//type UnderCells []UnderCell

func IsCell(data string) bool {
	if strings.HasPrefix(data, "cell") {
		return true
	}
	return false
}
