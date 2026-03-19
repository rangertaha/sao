package ui

import (
	"net/http"

	uiassets "github.com/rangertaha/sao/web/ui"
)

// AssetHandler returns an HTTP handler for embedded UI assets.
func AssetHandler() (http.Handler, error) {
	return uiassets.FileServer()
}
