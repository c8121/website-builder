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

	session := NewSession()

	for _, entity := range entities {
		err := go_db_entities.Create(session, entity)
		if err != nil {
			slog.Error("Failed to create entity", "error", err)
		}
		err = go_db_entities.Validate(session, entity)
		if err != nil {
			slog.Error("Failed to validate entity", "error", err)
		}
	}
}

func CreateDefaultSiteConfig() {

	session := NewSession()

	list, err := ListSiteConfig(session)
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
		err := go_db_entities.Save(session, config)
		if err != nil {
			slog.Error("Failed to save default config", "name", name, "value", value, "error", err)
		}
	}

}
