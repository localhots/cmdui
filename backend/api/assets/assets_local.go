// +build !binassets

package assets

import (
	"net/http"
)

func Handler() http.Handler {
	return http.NotFoundHandler()
}
