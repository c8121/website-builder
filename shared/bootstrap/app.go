package boostrap

import (
	assetstoragebootstrap "github.com/c8121/asset-storage/shared/bootstrap"
	"github.com/c8121/website-builder/shared/config"
)

func InitWebsiteBuiderApp() {

	config.AppConfig = config.NewConfigFromFlags()

	assetstoragebootstrap.InitAssetStorageApp() //Includes DB-Connect

}

func FinalizeWebsiteBuiderApp() {
	assetstoragebootstrap.FinalizeAssetStorageApp()
}
