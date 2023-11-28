package service

import (
	"static-server/filesystem"
)

type HTTPFileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mtime"`
}

type FileDetail struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Size    int64       `json:"size"`
	Path    string      `json:"path"`
	ModTime int64       `json:"mtime"`
	Extra   interface{} `json:"extra,omitempty"`
}

func (handler *FileServiceHandler) GetDetail(fn *filesystem.FileNode) *FileDetail {
	finfo := &FileDetail{
		Name:    fn.FileInfo.Name(),
		Type:    fn.GetFileType(),
		Size:    fn.Size,
		Path:    handler.fileTransformer.IsolationPath(fn.GetName()),
		ModTime: fn.FileInfo.ModTime().UnixNano() / 1e6,
	}
	if finfo.Type == filesystem.FILE_TYPE_APK {
		finfo.Extra = filesystem.GetApkInfo(fn.GetName())
	}
	return finfo
}

func (handler *FileServiceHandler) GetSubsFileDetail(fn *filesystem.FileNode) []FileDetail {
	res := make([]FileDetail, 0)
	for _, v := range fn.SubFiles {
		res = append(res, *handler.GetDetail(v))
	}
	return res
}

func (handler *FileServiceHandler) GetFilesDetail(fn []*filesystem.FileNode) []FileDetail {
	res := make([]FileDetail, 0)
	for _, v := range fn {
		res = append(res, *handler.GetDetail(v))
	}
	return res
}

type FileInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	Path    string `json:"path"`
	ModTime int64  `json:"mtime"`
}

func (handler *FileServiceHandler) GetFileInfo(fn *filesystem.FileNode) *FileInfo {
	finfo := &FileInfo{
		Name:    fn.FileInfo.Name(),
		Size:    fn.Size,
		Path:    handler.fileTransformer.IsolationPath(fn.GetName()),
		ModTime: fn.FileInfo.ModTime().UnixNano() / 1e6,
	}
	if fn.FileInfo.IsDir() {
		finfo.Type = filesystem.FILE_TYPE_DIR
	} else {
		finfo.Type = filesystem.FILE_TYPE_FILE
	}

	return finfo
}

func (handler *FileServiceHandler) GetSubsFileInfo(fn *filesystem.FileNode) []FileInfo {
	res := make([]FileInfo, 0)
	for _, v := range fn.SubFiles {
		res = append(res, *handler.GetFileInfo(v))
	}
	return res
}

func (handler *FileServiceHandler) GetFilesInfo(fn []*filesystem.FileNode) []FileInfo {
	res := make([]FileInfo, 0)
	for _, v := range fn {
		res = append(res, *handler.GetFileInfo(v))
	}
	return res
}

// func (handler *FileServiceHandler) handleJsonList(w http.ResponseWriter, r *http.Request) {
// 	requestPath := mux.Vars(r)["path"]
// 	search := r.FormValue("search")
// 	log.Printf("handleJsonList request path {%s}.\n", requestPath)

// 	realPath, _ := handler.fileTransformer.TransformPath(requestPath)

// 	auth := handler.readAccessConf(realPath.Get())
// 	auth.Upload = auth.canUpload(r)
// 	auth.Delete = auth.canDelete(r)

// 	// path string -> info os.FileInfo
// 	// fileInfoMap := make(map[string]os.FileInfo, 0)

// 	var fileList []FileInfo

// 	if search != "" {
// 		// results := handler.findIndex(search)
// 		// if len(results) > 50 { // max 50
// 		// 	results = results[:50]
// 		// }
// 		// for _, item := range results {
// 		// 	// fixme: search功能
// 		// 	if filepath.HasPrefix(item.Path, requestPath) {
// 		// 		fileInfoMap[item.Path] = item.Info
// 		// 	}
// 		// }

// 		handle := handler.fileTree.SearchFile(realPath.Get())
// 		fileList = handler.GetFilesInfo(handle)
// 		println(fileList)

// 	} else {
// 		// infos, err := ioutil.ReadDir(realPath.Get())
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), 500)
// 		// 	return
// 		// }
// 		// for _, info := range infos {
// 		// 	fileInfoMap[filepath.Join(requestPath, info.Name())] = info
// 		// }

// 		handle := handler.fileTree.GetFile(realPath.Get())
// 		fileList = handler.GetSubsFileInfo(handle)
// 		println(fileList)
// 	}

// 	// turn file list -> json
// 	// lrs := make([]HTTPFileInfo, 0)
// 	// for path, info := range fileInfoMap {
// 	// 	if !auth.canAccess(info.Name()) {
// 	// 		continue
// 	// 	}
// 	// 	lr := HTTPFileInfo{
// 	// 		Name:    info.Name(),
// 	// 		Path:    path,
// 	// 		ModTime: info.ModTime().UnixNano() / 1e6,
// 	// 	}
// 	// 	if search != "" {
// 	// 		name, err := filepath.Rel(requestPath, path)
// 	// 		if err != nil {
// 	// 			log.Println(requestPath, path, err)
// 	// 		}
// 	// 		lr.Name = filepath.ToSlash(name) // fix for windows
// 	// 	}
// 	// 	if info.IsDir() {
// 	// 		name := deepPath(realPath.Get(), info.Name())
// 	// 		lr.Name = name
// 	// 		lr.Path = filepath.Join(filepath.Dir(path), name)
// 	// 		lr.Type = filesystem.FILE_TYPE_DIR
// 	// 		lr.Size = handler.historyDirSize(string(realPath.Join(name)))
// 	// 	} else {
// 	// 		lr.Type = "file"
// 	// 		lr.Size = info.Size() // formatSize(info)
// 	// 	}
// 	// 	lrs = append(lrs, lr)
// 	// }

// 	data, _ := json.Marshal(
// 		ResFilesList{
// 			Files:  fileList,
// 			Access: auth,
// 		})
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(data)
// }
