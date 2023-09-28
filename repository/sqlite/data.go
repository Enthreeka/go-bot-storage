package sqlite

import (
	"context"
	"database/sql"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/repository"
)

type dataRepository struct {
	*SQLite
}

func NewDataRepository(SQLite *SQLite) repository.Data {
	return &dataRepository{
		SQLite,
	}
}

func (d *dataRepository) Create(ctx context.Context, data *model.Data) error {
	query := `INSERT INTO data (describe,under_cells_id) VALUES ($1,$2)`

	_, err := d.db.Exec(query, data.Describe, data.UnderCellID)
	return err
}

func (d *dataRepository) Update(ctx context.Context, data *model.Data) error {
	query := `UPDATE data SET describe = $1 WHERE under_cells_id = $2`

	_, err := d.db.Exec(query, data.Describe, data.UnderCellID)
	return err
}

func (d *dataRepository) GetByUnderCellID(ctx context.Context, underCellID int) (*model.Data, error) {
	query := `SELECT  describe FROM data WHERE under_cells_id = $1`

	data := &model.Data{}
	err := d.db.QueryRow(query, underCellID).Scan(&data.Describe)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return data, nil
}

func (d *dataRepository) GetDataByName(ctx context.Context, dataName string, underCellID int) (*model.Data, error) {
	query := `SELECT data.describe 
					FROM data
					JOIN under_cells ON data.under_cells_id = under_cells.id
					WHERE data.name = $1 AND data.under_cells_id=$2`
	data := &model.Data{}

	err := d.db.QueryRow(query, dataName, underCellID).Scan(&data.Describe)
	if err != nil {
		return nil, err
	}

	return data, nil
}
