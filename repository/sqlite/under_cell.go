package sqlite

import (
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/repository"
)

type underCellRepository struct {
	*SQLite
}

func NewUnderCellRepository(SQLite *SQLite) repository.UnderCell {
	return &underCellRepository{
		SQLite,
	}
}

func (u *underCellRepository) Create(cell *model.UnderCell) error {
	query := `INSERT INTO under_cells (name, cell_id) VALUES ($1,$2)`

	_, err := u.db.Exec(query, cell.Name, cell.CellID)
	return err
}

func (u *underCellRepository) DeleteByName(name string) error {
	query := `DELETE FROM under_cells WHERE name = $1`

	_, err := u.db.Exec(query, name)
	return err
}

func (u *underCellRepository) GetByCellID(id int64) ([]model.UnderCell, error) {
	query := `SELECT under_cells.id,under_cells.cell_id,under_cells.name
				FROM "user"
				JOIN cell ON cell.user_id = "user".id
				JOIN under_cells ON under_cells.cell_id = cell.id
				WHERE "user".id = $1`

	rows, err := u.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	underCalls := make([]model.UnderCell, 0, 16)
	for rows.Next() {
		var underCall model.UnderCell

		err = rows.Scan(&underCall.ID, &underCall.CellID, &underCall.Name)
		if err != nil {
			return nil, err
		}

		underCalls = append(underCalls, underCall)
	}

	return underCalls, nil
}
