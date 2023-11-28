package service

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"static-server/config"
	"static-server/filesystem"
	zzip "static-server/zip"
)

type IndexFileItem struct {
	Path string
	Info os.FileInfo
}

type FileServiceHandler struct {
	Config config.FileServiceConfig

	assets          http.FileSystem
	fileTransformer *filesystem.FileTransformer
	fileTree        *filesystem.FileTree
	indexes         []IndexFileItem
	muxRouter       *mux.Router
	bufPool         sync.Pool // use sync.Pool caching buf to reduce gc ratio
}

func (s *FileServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.muxRouter.ServeHTTP(w, r)
}

func createFileStaticService(conf config.FileServiceConfig, assets http.FileSystem) (handler *FileServiceHandler) {
	fileTree := filesystem.CreateFileTree(conf.Root)
	log.Println(fileTree)

	if err := conf.CheckPrefix(); err != nil {
		panic(err)
	}

	if err := conf.CheckRoot(); err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	handler = &FileServiceHandler{
		Config:          conf,
		assets:          assets,
		fileTransformer: filesystem.CreateFileTransformer(conf.Prefix, conf.Root),
		fileTree:        fileTree,
		muxRouter:       router,
		bufPool: sync.Pool{
			New: func() interface{} { return make([]byte, 32*1024) },
		},
	}

	go func() {
		time.Sleep(1 * time.Second)
		for {
			startTime := time.Now()
			log.Println("Started making search index")
			handler.makeIndex()
			log.Printf("Completed search index in %v", time.Since(startTime))
			//time.Sleep(time.Second * 1)
			time.Sleep(time.Minute * 10)
		}
	}()

	// routers for Apple *.ipa
	// m.HandleFunc("/-/ipa/plist/{path:.*}", s.hPlist)
	// m.HandleFunc("/-/ipa/link/{path:.*}", s.hIpaLink)

	// init router
	router.HandleFunc("/{path:.*}", handler.handleIndex).Methods("GET", "HEAD")
	router.HandleFunc("/{path:.*}", handler.handleUploadOrMkdir).Methods("POST")
	router.HandleFunc("/{path:.*}", handler.handleDelete).Methods("DELETE")

	return
}

func (handler *FileServiceHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	log.Printf("request path {%s}.\n", path)

	realPath, err := handler.fileTransformer.TransformPath(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.FormValue("json") == "true" {
		handler.handleJsonList(w, r)
		return
	}

	if r.FormValue("op") == "info" {
		handler.hInfo(w, r)
		return
	}

	if r.FormValue("op") == "archive" {
		handler.hZip(w, r)
		return
	}

	log.Println("GET", path, realPath)
	if r.FormValue("raw") == "false" || filesystem.IsDir(realPath.Get()) {
		if r.Method == "HEAD" {
			return
		}
		handler.renderHTML(w, "assets/index.html")
	} else {
		if filepath.Base(path) == config.ConfigYamlFile {
			auth := handler.readAccessConf(realPath.Get())
			if !auth.Delete {
				http.Error(w, "Security warning, not allowed to read", http.StatusForbidden)
				return
			}
		}
		if r.FormValue("download") == "true" {
			w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(path)))
		}
		http.ServeFile(w, r, realPath.Get())
	}
}

// func (handler *FileServiceHandler) handleJsonList(w http.ResponseWriter, r *http.Request) {
// 	requestPath := mux.Vars(r)["path"]
// 	log.Printf("handleJsonList request path {%s}.\n", requestPath)

// 	realPath, _ := handler.fileTransformer.TransformPath(requestPath)
// 	search := r.FormValue("search")
// 	auth := handler.readAccessConf(realPath.Get())
// 	auth.Upload = auth.canUpload(r)
// 	auth.Delete = auth.canDelete(r)

// 	finfo := handler.fileTree.GetFile(realPath.Get())
// 	println(finfo)

// 	// path string -> info os.FileInfo
// 	fileInfoMap := make(map[string]os.FileInfo, 0)

// 	if search != "" {
// 		results := handler.findIndex(search)
// 		if len(results) > 50 { // max 50
// 			results = results[:50]
// 		}
// 		for _, item := range results {
// 			// fixme: search功能
// 			if filepath.HasPrefix(item.Path, requestPath) {
// 				fileInfoMap[item.Path] = item.Info
// 			}
// 		}
// 	} else {
// 		infos, err := ioutil.ReadDir(realPath.Get())
// 		if err != nil {
// 			http.Error(w, err.Error(), 500)
// 			return
// 		}
// 		for _, info := range infos {
// 			fileInfoMap[filepath.Join(requestPath, info.Name())] = info
// 		}
// 	}

// 	// turn file list -> json
// 	lrs := make([]HTTPFileInfo, 0)
// 	for path, info := range fileInfoMap {
// 		if !auth.canAccess(info.Name()) {
// 			continue
// 		}
// 		lr := HTTPFileInfo{
// 			Name:    info.Name(),
// 			Path:    path,
// 			ModTime: info.ModTime().UnixNano() / 1e6,
// 		}
// 		if search != "" {
// 			name, err := filepath.Rel(requestPath, path)
// 			if err != nil {
// 				log.Println(requestPath, path, err)
// 			}
// 			lr.Name = filepath.ToSlash(name) // fix for windows
// 		}
// 		if info.IsDir() {
// 			name := deepPath(realPath.Get(), info.Name())
// 			lr.Name = name
// 			lr.Path = filepath.Join(filepath.Dir(path), name)
// 			lr.Type = filesystem.FILE_TYPE_DIR
// 			lr.Size = handler.historyDirSize(string(realPath.Join(name)))
// 		} else {
// 			lr.Type = "file"
// 			lr.Size = info.Size() // formatSize(info)
// 		}
// 		lrs = append(lrs, lr)
// 	}

// 	data, _ := json.Marshal(map[string]interface{}{
// 		"files": lrs,
// 		"auth":  auth,
// 	})
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(data)
// }

func (handler *FileServiceHandler) handleJsonList(w http.ResponseWriter, r *http.Request) {
	requestPath := mux.Vars(r)["path"]
	search := r.FormValue("search")
	log.Printf("handleJsonList request path {%s}.\n", requestPath)

	realPath, _ := handler.fileTransformer.TransformPath(requestPath)

	auth := handler.readAccessConf(realPath.Get())
	auth.Upload = auth.canUpload(r)
	auth.Delete = auth.canDelete(r)

	var fileList []FileInfo

	if search != "" {
		handle := handler.fileTree.SearchFile(search)
		fileList = handler.GetFilesInfo(handle)
		println(fileList)

	} else {
		handle := handler.fileTree.GetFile(realPath.Get())
		fileList = handler.GetSubsFileInfo(handle)
		println(fileList)
	}

	data, _ := json.Marshal(
		ResFilesList{
			Files:  fileList,
			Access: auth,
		})
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (handler *FileServiceHandler) hInfo(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	relPath, _ := handler.fileTransformer.TransformPath(path)
	finfo := handler.fileTree.GetFile(relPath.Get())
	log.Println(finfo)

	fi, err := os.Stat(relPath.Get())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fji := &FileDetail{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Path:    path,
		ModTime: fi.ModTime().UnixNano() / 1e6,
	}
	ext := filepath.Ext(path)
	switch ext {
	case filesystem.FILE_TYPE_MARKDOWN_SUFFIX:
		fji.Type = filesystem.FILE_TYPE_MARKDOWN
	case filesystem.FILE_TYPE_APK_SUFFIX:
		fji.Type = filesystem.FILE_TYPE_APK
		fji.Extra = filesystem.GetApkInfo(relPath.Get())
	case "":
		fji.Type = filesystem.FILE_TYPE_DIR
	default:
		fji.Type = filesystem.FILE_TYPE_TEXT
	}
	data, _ := json.Marshal(fji)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (handler *FileServiceHandler) handleDelete(w http.ResponseWriter, req *http.Request) {
	path := mux.Vars(req)["path"]
	log.Printf("delete path {%s}.\n", path)

	realPath, _ := handler.fileTransformer.TransformPath(path)
	// path = filepath.Clean(path) // for safe reason, prevent path contain ..
	auth := handler.readAccessConf(realPath.Get())
	if !auth.canDelete(req) {
		http.Error(w, "Delete forbidden", http.StatusForbidden)
		return
	}

	// TODO: path safe check
	err := os.RemoveAll(realPath.Get())
	if err != nil {
		pathErr, ok := err.(*os.PathError)
		if ok {
			http.Error(w, pathErr.Op+" "+path+": "+pathErr.Err.Error(), 500)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	w.Write([]byte("Success"))
}

func (handler *FileServiceHandler) handleUploadOrMkdir(w http.ResponseWriter, req *http.Request) {
	path := mux.Vars(req)["path"]
	log.Printf("upload mkdir path {%s}.\n", path)

	dirpath, _ := handler.fileTransformer.TransformPath(path)

	// check auth
	auth := handler.readAccessConf(dirpath.Get())
	if !auth.canUpload(req) {
		http.Error(w, "Upload forbidden", http.StatusForbidden)
		return
	}

	file, header, err := req.FormFile("file")

	if _, err := os.Stat(dirpath.Get()); os.IsNotExist(err) {
		if err := os.MkdirAll(dirpath.Get(), os.ModePerm); err != nil {
			log.Println("Create directory:", err)
			http.Error(w, "Directory create "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if file == nil { // only mkdir
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"destination": dirpath,
		})
		return
	}

	if err != nil {
		log.Println("Parse form file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		file.Close()
		req.MultipartForm.RemoveAll() // Seen from go source code, req.MultipartForm not nil after call FormFile(..)
	}()

	filename := req.FormValue("filename")
	if filename == "" {
		filename = header.Filename
	}
	if err := checkFilename(filename); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	dstPath := dirpath.Join(filename)

	// Large file (>32MB) will store in tmp directory
	// The quickest operation is call os.Move instead of os.Copy
	// Note: it seems not working well, os.Rename might be failed

	var copyErr error
	// if osFile, ok := file.(*os.File); ok && fileExists(osFile.Name()) {
	// 	tmpUploadPath := osFile.Name()
	// 	osFile.Close() // Windows can not rename opened file
	// 	log.Printf("Move %s -> %s", tmpUploadPath, dstPath)
	// 	copyErr = os.Rename(tmpUploadPath, dstPath)
	// } else {
	dst, err := os.Create(dstPath.Get())
	if err != nil {
		log.Println("Create file:", err)
		http.Error(w, "File create "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Note: very large size file might cause poor performance
	// _, copyErr = io.Copy(dst, file)
	buf := handler.bufPool.Get().([]byte)
	defer handler.bufPool.Put(buf)
	_, copyErr = io.CopyBuffer(dst, file, buf)
	dst.Close()
	// }
	if copyErr != nil {
		log.Println("Handle upload file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	if req.FormValue("unzip") == "true" {
		err = zzip.UnzipFile(dstPath.Get(), dirpath.Get())
		os.Remove(dstPath.Get())
		message := "success"
		if err != nil {
			message = err.Error()
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     err == nil,
			"description": message,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"destination": dstPath,
	})
}

func checkFilename(name string) error {
	if strings.ContainsAny(name, "\\/:*<>|") {
		return errors.New("fiel name should not contains \\/:*<>| ")
	}
	return nil
}

var (
	_tmpls = make(map[string]*template.Template)
)

func (handler *FileServiceHandler) renderHTML(w http.ResponseWriter, name string) {
	if t, ok := _tmpls[name]; ok {
		t.Execute(w, handler)
		return
	}
	t := template.Must(template.New(name).Delims("[[", "]]").Parse(handler.assetsContent(name)))
	_tmpls[name] = t
	t.Execute(w, handler)
}

func (handler *FileServiceHandler) assetsContent(name string) string {
	fd, err := handler.assets.Open(name)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func deepPath(basedir, name string) string {
	// loop max 5, incase of for loop not finished
	maxDepth := 5
	for depth := 0; depth <= maxDepth; depth += 1 {
		finfos, err := ioutil.ReadDir(filepath.Join(basedir, name))
		if err != nil || len(finfos) != 1 {
			break
		}
		if finfos[0].IsDir() {
			name = filepath.ToSlash(filepath.Join(name, finfos[0].Name()))
		} else {
			break
		}
	}
	return name
}

func (handler *FileServiceHandler) findIndex(text string) []IndexFileItem {
	ret := make([]IndexFileItem, 0)
	for _, item := range handler.indexes {
		ok := true
		// search algorithm, space for AND
		for _, keyword := range strings.Fields(text) {
			needContains := true
			if strings.HasPrefix(keyword, "-") {
				needContains = false
				keyword = keyword[1:]
			}
			if keyword == "" {
				continue
			}
			ok = (needContains == strings.Contains(strings.ToLower(item.Path), strings.ToLower(keyword)))
			if !ok {
				break
			}
		}
		if ok {
			ret = append(ret, item)
		}
	}
	return ret
}

// todo:有点乱，重构文件系统
type Directory struct {
	size  map[string]int64
	mutex *sync.RWMutex
}

var dirInfoSize = Directory{
	size:  make(map[string]int64),
	mutex: &sync.RWMutex{},
}

func (handler *FileServiceHandler) historyDirSize(dir string) int64 {
	dirInfoSize.mutex.RLock()
	size, ok := dirInfoSize.size[dir]
	dirInfoSize.mutex.RUnlock()

	if ok {
		return size
	}

	// fixme:
	for _, fitem := range handler.indexes {
		p := filepath.Dir(fitem.Path)
		if strings.HasPrefix(p, dir) && fitem.Path != dir {
			size += fitem.Info.Size()
		}
	}

	dirInfoSize.mutex.Lock()
	dirInfoSize.size[dir] = size
	dirInfoSize.mutex.Unlock()

	return size
}

// todo: thread safe ???
func (handler *FileServiceHandler) makeIndex() error {
	var indexes = make([]IndexFileItem, 0)
	var err = filepath.Walk(handler.Config.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("WARN: Visit path: %s error: %v", strconv.Quote(path), err)
			return filepath.SkipDir
			// return err
		}
		if info.IsDir() {
			return nil
		}

		// path, _ = filepath.Rel(s.Config.Root, path)
		path = filepath.ToSlash(path)
		path, err = filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		indexes = append(indexes, IndexFileItem{path, info})
		return nil
	})
	handler.indexes = indexes
	return err
}
