package main

import (
	"log/slog"

	"github.com/c8121/website-builder/internal/builder"
	"github.com/c8121/website-builder/shared/bootstrap"
	"github.com/c8121/website-builder/shared/config"
)

func main() {

	boostrap.InitWebsiteBuilderApp()
	defer boostrap.FinalizeWebsiteBuilderApp()

	if config.AppConfig.TemplateDir == "" {
		slog.Error("No template selected")
		return
	}

	pageBuilder, err := builder.NewBuilder(
		config.AppConfig.TemplateDir,
		config.AppConfig.OutputDir,
	)
	if err != nil {
		slog.Error("Failed to create builder", "error", err)
		return
	}

	pageBuilder.Build()
}
