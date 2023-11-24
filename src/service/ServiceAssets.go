package service

import (
	"embed"
	"net/http"
)

type AssetsHandler struct {
	assets http.FileSystem
}

func (handler *AssetsHandler) Set(assets embed.FS) {
	handler.assets = http.FS(assets)
}

func (handler *AssetsHandler) Get() http.Handler {
	return http.FileServer(handler.assets)
}
