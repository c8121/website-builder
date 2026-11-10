package builder

import (
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// PostprocessorInsertAssets inserts assets
type PostprocessorInsertAssets struct {
	Config *BuildConfig
}

func NewPostprocessorInsertAssets(config *BuildConfig) *PostprocessorInsertAssets {
	return &PostprocessorInsertAssets{
		Config: config,
	}
}

func (p *PostprocessorInsertAssets) Process(doc *goquery.Document) error {
	return traverseDocumentNodes(doc, func(n *html.Node) error {
		for i, attr := range n.Attr {

			key := strings.ToLower(attr.Key)
			switch key {
			case "src", "href":
				if p.Config.AssetLoader.IsAssetUrl(attr.Val) {
					assetSrc, err := p.Config.AssetLoader.CopyAsset(attr.Val)
					if err == nil {
						n.Attr[i].Val = assetSrc
					} else {
						slog.Warn("Failed to load asset from URL", "src", attr.Val, "err", err)
					}
				}
			}
		}

		return nil
	})
}
