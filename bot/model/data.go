package model

import "strings"

type Data struct {
	ID          int    `json:"id"`
	UnderCellID int    `json:"under_cells_id"`
	Name        string `json:"name"`
	Describe    string `json:"describe"`
}

func IsData(data string) bool {
	if strings.HasPrefix(data, data) {
		return true
	}
	return false
}

func IsFile(data string) (string, bool) {
	prefix := "file-"

	if strings.HasPrefix(data, prefix) {
		data = strings.TrimPrefix(data, prefix)
		return data, true
	}
	return "", false
}
