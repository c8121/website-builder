package entity

import (
	"log/slog"

	"github.com/c8121/go_db_entities"
)

func CheckCreateTables() {

	entities := []go_db_entities.Entity{
		&Page{},
		&PageContent{},
	}

	for _, entity := range entities {
		err := go_db_entities.Create(DbContext, entity)
		if err != nil {
			slog.Error("Failed to create entity", "error", err)
		}
	}
}
