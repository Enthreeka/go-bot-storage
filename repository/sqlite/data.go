package sqlite

import (
	"github.com/Enthreeka/go-bot-storage/bot/model"
)

type dataRepository struct {
	*SQLite
}

func NewDataRepository(SQLite *SQLite) *dataRepository {
	return &dataRepository{
		SQLite,
	}
}

func (d *dataRepository) Create(data model.Data) error {
	query := `INSERT INTO data (name) VALUES ($1)`

	d.db.Exec(query)

	panic("implement me")
}

func (d *dataRepository) Delete(name string) error {
	//query := `DELETE FROM data WHERE under_cells_id = $1`
	panic("delete")
}

func (d *dataRepository) GetByUnderCellID(underCellID int) ([]model.Data, error) {
	query := `SELECT id,name FROM data WHERE under_cells_id = $1`

	rows, err := d.db.Query(query, underCellID)
	if err != nil {
		return nil, err
	}

	datas := make([]model.Data, 0, 16)
	for rows.Next() {
		var data model.Data

		err = rows.Scan(&data.ID, &data.Name)
		if err != nil {
			return nil, err
		}

		datas = append(datas, data)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return datas, nil
}

func (d *dataRepository) GetDataByName(dataName string, underCellID int) (*model.Data, error) {
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

//func (d *dataRepository) GetDescribeByName() {
//	query := `SELECT `
//
//}
