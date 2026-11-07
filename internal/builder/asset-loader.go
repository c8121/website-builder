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

type AssetLoader struct {
	assetStorageUrl string
	outputBaseDir   string

	client *http.Client

	nameCache map[string]string
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
	}
}

// CopyAsset downloads the asset and its metadata and stores it into outputBaseDir, using original name (normalized)
func (a *AssetLoader) CopyAsset(src string) (string, error) {

	if strings.TrimSpace(src) == "" {
		return "", errors.New("empty source path")
	}

	cached, ok := a.nameCache[src]
	if ok {
		return cached, nil
	}

	hash, err := a.extractAssetHash(src)
	if err != nil {
		slog.Warn("Error extracting asset hash", "error", err)
		return a.DownloadAsset(src)
	}

	meta, err := a.DownloadAssetMetaData(hash)
	if err != nil {
		slog.Warn("Error downloading asset meta", "error", err)
		return a.DownloadAsset(src)
	}

	name, err := a.ExtractNameFromMetaData(meta)
	if err != nil {
		slog.Warn("Error extracting asset name", "error", err)
		name = hash
	}

	dest, err := a.GetDestinationName(name)
	if err != nil {
		slog.Warn("Error getting destination", "error", err)
		return a.DownloadAsset(src)
	}

	slog.Info("Inserting asset", "hash", hash, "dest", dest)

	temp, err := a.DownloadAsset(src)
	if err != nil {
		return "", err
	}

	err = os.Rename(temp, dest)
	if err != nil {
		return "", err
	}

	a.nameCache[src] = dest

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

func (a *AssetLoader) extractAssetHash(src string) (string, error) {

	if !strings.HasPrefix(src, "/assets/") {
		slog.Warn("Source does not start with /assets/")
	}

	p := strings.LastIndex(src, "/")
	if p < 0 {
		return "", errors.New("cannot extract asset hash")
	}

	return src[p+1:], nil
}

// DownloadAsset loads the asset from asset-storage and stores it as temp-file into outputBaseDir
func (a *AssetLoader) DownloadAsset(src string) (string, error) {

	src = a.assetStorageUrl + src
	slog.Info("Download asset", "src", src)

	tmp, err := os.CreateTemp(a.outputBaseDir, "")
	if err != nil {
		return "", err
	}
	defer util.CloseOrLog(tmp)

	resp, err := a.client.Get(src)
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
