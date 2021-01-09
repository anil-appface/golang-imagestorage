package store

import "database/sql"

//Prepper an abstraction around the database to prepare the statement
type Prepper interface {
	Prepare(query string) (*sql.Stmt, error)
}
