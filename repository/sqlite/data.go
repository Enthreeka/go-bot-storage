package sqlite

import (
	"github.com/Enthreeka/go-bot-storage/bot/model"
)

type dataRepository struct {
	*SQLite
}

func DataRepository(SQLite *SQLite) *dataRepository {
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
	//TODO implement me
	panic("implement me")
}

func (d *dataRepository) GetByUnderCellID() (*model.Data, error) {
	//TODO implement me
	panic("implement me")
}
