package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/c8121/go_vue_spa/util"
)

// AssetLoader takes urls pointing to assets from asset-storage in copies them to outputBase dir with original name.
type AssetLoader struct {
	assetStorageUrl string
	outputBaseDir   string

	client *http.Client

	nameCache map[string]string

	assetUrlPrefixes []string
}

func NewAssetLoader(assetStorageUrl, outputBaseDir string) *AssetLoader {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return &AssetLoader{
		assetStorageUrl: assetStorageUrl,
		outputBaseDir:   outputBaseDir,
		client:          client,
		nameCache:       make(map[string]string),
		assetUrlPrefixes: []string{
			"/assets/",
		},
	}
}

func (a *AssetLoader) IsAssetUrl(url string) bool {
	if strings.TrimSpace(url) == "" {
		return false
	}
	for _, prefix := range a.assetUrlPrefixes {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}
	return false
}

// CopyAsset downloads the asset and its metadata and stores it into outputBaseDir, using original name (normalized)
func (a *AssetLoader) CopyAsset(url string) (string, error) {

	if !a.IsAssetUrl(url) {
		return url, nil
	}

	cached, ok := a.nameCache[url]
	if ok {
		return cached, nil
	}

	hash, err := a.extractAssetHash(url)
	if err != nil {
		slog.Warn("Error extracting asset hash", "error", err)
		return a.DownloadAsset(url)
	}

	meta, err := a.DownloadAssetMetaData(hash)
	if err != nil {
		slog.Warn("Error downloading asset meta", "error", err)
		return a.DownloadAsset(url)
	}

	name, err := a.ExtractNameFromMetaData(meta)
	if err != nil {
		slog.Warn("Error extracting asset name", "error", err)
		name = hash
	}

	dest, err := a.GetDestinationName(name)
	if err != nil {
		slog.Warn("Error getting destination", "error", err)
		return a.DownloadAsset(url)
	}

	slog.Info("Inserting asset", "hash", hash, "dest", dest)

	temp, err := a.DownloadAsset(url)
	if err != nil {
		return "", err
	}

	err = os.Rename(temp, dest)
	if err != nil {
		return "", err
	}

	a.nameCache[url] = dest

	return dest, nil
}

// GetDestinationName generates a unique name for the asset, using the given name as template
func (a *AssetLoader) GetDestinationName(name string) (string, error) {

	name = NormalizeAsFileName(name, fmt.Sprintf("asset-%d", time.Now().UnixNano()))
	basename, extension := a.splitName(name)

	num := 0
	path := filepath.Join(a.outputBaseDir, name)
	for a.checkExists(path) {
		num++
		if num > 1000 {
			return "", errors.New("too many files")
		}
		path = filepath.Join(a.outputBaseDir, fmt.Sprintf("%s(%d).%s", basename, num, extension))
	}

	return path, nil
}

func (a *AssetLoader) splitName(name string) (string, string) {
	p := strings.LastIndexByte(name, '.')
	if p < 0 {
		return name, ""
	}
	return name[:p], name[p+1:]
}

func (a *AssetLoader) checkExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (a *AssetLoader) ExtractNameFromMetaData(meta map[string]any) (string, error) {
	origins, ok := meta["Origins"]
	if !ok {
		return "", fmt.Errorf("missing Origins: %v", meta)
	}

	list, ok := origins.([]any)
	if !ok || len(list) == 0 {
		return "", fmt.Errorf("not a list in Origins: %v", meta)
	}

	name, ok := list[0].(map[string]any)["Name"].(string)
	if !ok {
		return "", fmt.Errorf("list is empty or Name is missing: %v", list)
	}

	return name, nil
}

func (a *AssetLoader) extractAssetHash(url string) (string, error) {

	if !strings.HasPrefix(url, "/assets/") {
		slog.Warn("Source does not start with /assets/")
	}

	p := strings.LastIndex(url, "/")
	if p < 0 {
		return "", errors.New("cannot extract asset hash")
	}

	return url[p+1:], nil
}

// DownloadAsset loads the asset from asset-storage and stores it as temp-file into outputBaseDir
func (a *AssetLoader) DownloadAsset(url string) (string, error) {

	url = a.assetStorageUrl + url
	slog.Info("Download asset", "src", url)

	tmp, err := os.CreateTemp(a.outputBaseDir, "")
	if err != nil {
		return "", err
	}
	defer util.CloseOrLog(tmp)

	resp, err := a.client.Get(url)
	if err != nil {
		util.LogError(os.Remove(tmp.Name()))
		return "", err
	}
	defer util.CloseOrLog(resp.Body)

	if resp.StatusCode != http.StatusOK {
		util.LogError(os.Remove(tmp.Name()))
		return "", err
	}

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		util.LogError(os.Remove(tmp.Name()))
		return "", err
	}

	return tmp.Name(), nil
}

func (a *AssetLoader) DownloadAssetMetaData(hash string) (map[string]any, error) {

	src := a.assetStorageUrl + "/assets/metadata/" + hash
	slog.Info("Download asset metadata", "src", src)

	resp, err := a.client.Get(src)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
