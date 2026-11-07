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

var (
	client = &http.Client{
		Timeout: time.Second * 10,
	}

	AssetStorageUrl = "http://localhost:9998"
	OutputBaseDir   = "/tmp"
)

func InsertAsset(src string) (string, error) {

	hash, err := extractAssetHash(src)
	if err != nil {
		slog.Warn("Error extracting asset hash", "error", err)
		return DownloadAsset(src)
	}

	meta, err := DownloadAssetMetaData(hash)
	if err != nil {
		slog.Warn("Error downloading asset meta", "error", err)
		return DownloadAsset(src)
	}

	name, err := ExtractNameFromMetaData(meta)
	if err != nil {
		slog.Warn("Error extracting asset name", "error", err)
		name = hash
	}

	dest, err := GetDestination(name)
	if err != nil {
		slog.Warn("Error getting destination", "error", err)
		return DownloadAsset(src)
	}

	slog.Info("Inserting asset", "hash", hash, "dest", dest)

	temp, err := DownloadAsset(src)
	if err != nil {
		return "", err
	}

	err = os.Rename(temp, dest)
	if err != nil {
		return "", err
	}

	return dest, nil
}

func GetDestination(name string) (string, error) {
	dir := OutputBaseDir

	name = NormalizeAsFileName(name, fmt.Sprintf("asset-%d", time.Now().UnixNano()))
	basename, extension := splitName(name)

	num := 0
	path := filepath.Join(dir, name)
	for checkExists(path) {
		num++
		if num > 1000 {
			return "", errors.New("too many files")
		}
		path = filepath.Join(dir, fmt.Sprintf("%s(%d).%s", basename, num, extension))
	}

	return path, nil
}

func splitName(name string) (string, string) {
	p := strings.LastIndexByte(name, '.')
	if p < 0 {
		return name, ""
	}
	return name[:p], name[p+1:]
}

func checkExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ExtractNameFromMetaData(meta map[string]any) (string, error) {
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

func extractAssetHash(src string) (string, error) {

	if !strings.HasPrefix(src, "/assets/") {
		slog.Warn("Source does not start with /assets/")
	}

	p := strings.LastIndex(src, "/")
	if p < 0 {
		return "", errors.New("cannot extract asset hash")
	}

	return src[p+1:], nil
}

func DownloadAsset(src string) (string, error) {

	src = AssetStorageUrl + src
	slog.Info("Download asset", "src", src)

	tmp, err := os.CreateTemp(OutputBaseDir, "")
	if err != nil {
		return "", err
	}
	defer util.CloseOrLog(tmp)

	resp, err := client.Get(src)
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

func DownloadAssetMetaData(hash string) (map[string]any, error) {

	src := AssetStorageUrl + "/assets/metadata/" + hash
	slog.Info("Download asset metadata", "src", src)

	resp, err := client.Get(src)
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
