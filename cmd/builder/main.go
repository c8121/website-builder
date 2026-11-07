package main

import (
	"log/slog"

	builder2 "github.com/c8121/website-builder/internal/builder"
	boostrap "github.com/c8121/website-builder/shared/bootstrap"
	"github.com/c8121/website-builder/shared/config"
)

func main() {

	boostrap.InitWebsiteBuiderApp()
	defer boostrap.FinalizeWebsiteBuilderApp()

	if config.AppConfig.TemplateDir == "" {
		slog.Error("No template selected")
		return
	}

	builder, err := builder2.NewBuilder(
		config.AppConfig.TemplateDir,
		config.AppConfig.OutputDir,
	)
	if err != nil {
		slog.Error("Failed to create builder", "error", err)
		return
	}

	err = builder.Build()
	if err != nil {
		slog.Error("Failed to build", "error", err)
		return
	}

}
