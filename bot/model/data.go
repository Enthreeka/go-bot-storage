package model

import "strings"

type Data struct {
	ID          int    `json:"id"`
	UnderCellID int    `json:"under_cells_id"`
	Describe    string `json:"describe"`
}

func IsFile(data string) (string, bool) {
	prefix := "file-"

	if strings.HasPrefix(data, prefix) {
		data = strings.TrimPrefix(data, prefix)
		return data, true
	}
	return "", false
}
