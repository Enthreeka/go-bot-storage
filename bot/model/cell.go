package model

type Cell struct {
	ID     int
	Name   string
	UserID int
}

type UnderCell struct {
	ID     int
	CellID int
	Name   string
}
