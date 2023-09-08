package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	tableUser = `CREATE TABLE IF NOT EXISTS  "user"(
    id INTEGER primary key  AUTOINCREMENT,
    nickname varchar(100) unique not null,
    first_name varchar(100) not null,
    last_name varchar(100) not null,
    role varchar(5) default 'user'
														)`

	tableCell = `CREATE TABLE IF NOT EXISTS cell(
    id INTEGER primary key AUTOINCREMENT,
    name varchar(200) not null,
    user_id int,
    foreign key (user_id)
    references "user" (id) on delete cascade
												)`

	tableUnderCell = `CREATE TABLE IF NOT EXISTS under_cells(
   id INTEGER primary key AUTOINCREMENT,
   cell_id int,
   name varchar(200) not null,
   foreign key (cell_id)
    references cell (id) on delete cascade
															)`

	tableData = `CREATE TABLE IF NOT EXISTS data(
  id INTEGER primary key AUTOINCREMENT,
  under_cells_id int,
  link text null,
  describe text null,
  pdf_link text null,
  foreign key (under_cells_id)
      references under_cells (id) on delete cascade
												)`
)

type SQLite struct {
	db *sql.DB
}

func New() (*SQLite, error) {
	database, err := sql.Open("sqlite3", "./bot_lite.db")
	if err != nil {
		return nil, err
	}

	statement, err := database.Prepare(tableUser)
	if err != nil {
		return nil, err
	}
	statement.Exec()

	statement, err = database.Prepare(tableCell)
	if err != nil {
		return nil, err
	}
	statement.Exec()

	statement, err = database.Prepare(tableUnderCell)
	if err != nil {
		return nil, err
	}
	statement.Exec()

	statement, err = database.Prepare(tableData)
	if err != nil {
		return nil, err
	}
	statement.Exec()

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	db := &SQLite{
		db: database,
	}

	return db, nil
}
