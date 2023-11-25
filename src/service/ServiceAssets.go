package service

import (
	"net/http"
)

type AssetsHandler struct {
	assets http.FileSystem
}

func (handler *AssetsHandler) Set(assets http.FileSystem) {
	handler.assets = assets
}

func (handler *AssetsHandler) Get() http.Handler {
	return http.FileServer(handler.assets)
}
