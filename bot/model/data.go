package model

type Data struct {
	ID          int    `json:"id"`
	UnderCellID int    `json:"under_cells_id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Describe    string `json:"describe"`
	PDFLink     string `json:"pdf_link"`
}
