package entity

import (
	"database/sql"

	"github.com/c8121/go_db_entities"
)

var (
	Storage *go_db_entities.Storage
)

func NewSession() *go_db_entities.Session {
	return go_db_entities.NewSession(Storage)
}

func NewSessionTx() (*go_db_entities.Session, *sql.Tx, error) {
	return go_db_entities.NewSessionTx(Storage)
}
