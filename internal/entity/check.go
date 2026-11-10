package entity

import (
	"log/slog"

	"github.com/c8121/go_db_entities"
)

func CheckCreateTables() {

	entities := []go_db_entities.Entity{
		&Page{},
		&PageContent{},
		&SiteConfig{},
	}

	for _, entity := range entities {
		err := go_db_entities.Create(DbContext, entity)
		if err != nil {
			slog.Error("Failed to create entity", "error", err)
		}
	}
}

func CreateDefaultSiteConfig() {

	list, err := ListSiteConfig()
	if err != nil {
		slog.Error("Cannot create default config", "error", err)
		return
	}

	if len(list) > 0 {
		slog.Debug("Config exists, not creating default values")
		return
	}

	items := map[string]string{
		"Title":    "Awesome Site",
		"SubTitle": "Created staticly",
	}

	for name, value := range items {
		config := &SiteConfig{
			Name:  name,
			Value: value,
		}
		err := go_db_entities.Save(DbContext, config)
		if err != nil {
			slog.Error("Failed to save default config", "name", name, "value", value, "error", err)
		}
	}

}
