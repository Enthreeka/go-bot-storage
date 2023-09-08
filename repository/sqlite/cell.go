package sqlite

import (
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/repository"
)

type cellRepository struct {
	*SQLite
}

func CellRepository(SQLite *SQLite) repository.Cell {
	return &cellRepository{
		SQLite,
	}
}

func (c *cellRepository) Create(cell *model.Cell) error {
	query := `INSERT INTO cell (name,user_id) VALUES ($1,$2)`

	_, err := c.db.Exec(query, cell.Name, cell.UserID)
	return err
}

func (c *cellRepository) DeleteByName(name string) error {
	query := `DELETE FROM cell WHERE name = $1`

	_, err := c.db.Exec(query, name)
	return err
}

func (c *cellRepository) GetByUserID(id int64) ([]model.Cell, error) {
	query := `SELECT id,name,user_id FROM cell WHERE user_id = $1`

	rows, err := c.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	cells := make([]model.Cell, 0, 32)
	for rows.Next() {
		var cell model.Cell

		err = rows.Scan(&cell.ID, &cell.Name, &cell.UserID)
		if err != nil {
			return nil, err
		}

		cells = append(cells, cell)
	}

	return cells, nil
}
