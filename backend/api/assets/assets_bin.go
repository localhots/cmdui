// +build binassets

package assets

import (
	"net/http"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

func Handler() http.Handler {
	// Following functions would be defined in bindata_assetfs.go during the
	// build process
	return http.FileServer(&assetfs.AssetFS{
		Asset:     reactRouter,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "/",
	})
}

func reactRouter(path string) ([]byte, error) {
	reactPrefixes := []string{"cmd", "jobs", "users"}
	for _, prefix := range reactPrefixes {
		if strings.HasPrefix(path, prefix) {
			return Asset("index.html")
		}

	}

	return Asset(path)
}
