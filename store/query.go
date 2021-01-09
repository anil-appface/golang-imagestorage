package store

import (
	"database/sql"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS 
metafest(
	qid CHAR(32) PRIMARY KEY,
	parent_qid CHAR(32) NOT NULL,
	mode INT NOT NULL,
	name CHAR(200),
	size BIGINT,
	timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

//InitialiseTables it intialises tables in the database.
func InitialiseTables(db *sql.DB) error {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
