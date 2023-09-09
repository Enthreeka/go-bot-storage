package model

type Data struct {
	ID          int    `json:"id"`
	UnderCellID int    `json:"under_cells_id"`
	Name        string `json:"name"`
	Describe    string `json:"describe"`
}

type NameData []string
