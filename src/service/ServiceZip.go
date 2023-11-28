package service

import (
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"

	zzip "static-server/zip"
)

func (handler *FileServiceHandler) hZip(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	realPath, _ := handler.fileTransformer.TransformPath(path)
	zzip.CompressToZip(w, realPath.Get())
}

func (handler *FileServiceHandler) hUnzip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	zipPath, path := vars["zip_path"], vars["path"]
	ctype := mime.TypeByExtension(filepath.Ext(path))
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}
	err := zzip.ExtractFromZip(filepath.Join(handler.Config.Root, zipPath), path, w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
