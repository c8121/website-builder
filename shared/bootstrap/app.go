package boostrap

import (
	"flag"

	assetstoragebootstrap "github.com/c8121/asset-storage/shared/bootstrap"
	assetstorageconfig "github.com/c8121/asset-storage/shared/config"
	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_db_entities/sqlite"
	"github.com/c8121/website-builder/internal/entity"
	"github.com/c8121/website-builder/shared/config"
)

func InitWebsiteBuiderApp() {

	config.AppConfig = config.NewConfigFromFlags(flag.ExitOnError)
	assetstorageconfig.AppConfig = assetstorageconfig.NewConfigFromFlags(flag.ContinueOnError)

	assetstorageconfig.SetDirectories(assetstorageconfig.AppConfig, config.AppConfig.ProjectDir)

	assetstoragebootstrap.CreateAppDir(assetstorageconfig.AppConfig.AssetStorageBaseDir)
	assetstoragebootstrap.CreateAppDir(assetstorageconfig.AppConfig.AssetStorageTempDir)

	assetstoragebootstrap.CreateAppDir(assetstorageconfig.AppConfig.AssetMetaDataBaseDir)

	assetstoragebootstrap.CreateAppDir(assetstorageconfig.AppConfig.AssetCollectionsBaseDir)

	assetstoragebootstrap.CreateAppDir(assetstorageconfig.AppConfig.AssetFulltextIndexBaseDir)

	assetstoragebootstrap.OpenDatabase()

	entity.DbContext = go_db_entities.NewContext(
		assetstoragebootstrap.GetDatabase(),
		sqlite.NewQueryBuilder(),
		sqlite.NewFieldMapping(),
	)

	entity.CheckCreateTables()

}

func FinalizeWebsiteBuiderApp() {
	assetstoragebootstrap.FinalizeAssetStorageApp()
}
