package sqlite

import (
	"context"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/repository"
)

type cellRepositorySL struct {
	*SQLite
}

func NewCellRepositorySL(SQLite *SQLite) repository.Cell {
	return &cellRepositorySL{
		SQLite,
	}
}

func (c *cellRepositorySL) Create(ctx context.Context, cell *model.Cell) error {
	query := `INSERT INTO cell (name,user_id) VALUES ($1,$2)`

	_, err := c.db.Exec(query, cell.Name, cell.UserID)
	return err
}

func (c *cellRepositorySL) DeleteByID(ctx context.Context, id int) error {
	query := `DELETE FROM cell WHERE id = $1`

	_, err := c.db.Exec(query, id)
	return err
}

// TODO cделать обработку ошибку no one result rows
func (c *cellRepositorySL) GetByUserID(ctx context.Context, id int64) ([]model.Cell, error) {
	query := `SELECT id,name,user_id FROM cell WHERE user_id = $1`

	rows, err := c.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	cells := make([]model.Cell, 0, 16)
	for rows.Next() {
		var cell model.Cell

		err = rows.Scan(&cell.ID, &cell.Name, &cell.UserID)
		if err != nil {
			return nil, err
		}

		cells = append(cells, cell)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cells, nil
}
