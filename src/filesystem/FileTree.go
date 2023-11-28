package filesystem

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ftree  *FileTree
	fTrans *FileTransformer
)

type FileNode struct {
	Path     string // abs, parent path
	Size     int64  // file size or direction total size
	FileInfo fs.FileInfo
	Parent   *FileNode
	SubFiles map[string]*FileNode
}

func (fn FileNode) GetName() string {
	return fn.Path + "/" + fn.FileInfo.Name()
}

func (fn FileNode) GetFileType() string {
	if fn.FileInfo.IsDir() {
		return FILE_TYPE_DIR
	}

	ext := filepath.Ext(fn.FileInfo.Name())
	switch ext {
	case FILE_TYPE_MARKDOWN_SUFFIX:
		return FILE_TYPE_MARKDOWN
	case FILE_TYPE_APK_SUFFIX:
		return FILE_TYPE_APK
	default:
		return FILE_TYPE_TEXT
	}
}

// 从根向子叶递归
func (fn *FileNode) addDir(path []string) {
	if len(path) == 0 {
		return
	}

	fi := GetFileInfo(fn.Path + "/" + path[0])
	if fi == nil {
		panic("dir not exist.")
	}

	fn.SubFiles[path[0]] = &FileNode{
		Path:     fn.Path,
		FileInfo: fi,
		SubFiles: map[string]*FileNode{},
	}
	// append to cache
	ftree.Files[path[0]] = append(ftree.Files[path[0]], fn.SubFiles[path[0]])

	// 递归
	fn.SubFiles[path[0]].addDir(path[1:])
}

// 从根向子叶递归
func (fn *FileNode) addFile(path []string, fileName string) int64 {
	if len(path) == 0 {
		finfo := GetFileInfo(fn.Path + "/" + fileName)

		// make file node
		fnn := &FileNode{
			Path:     fn.Path,
			FileInfo: finfo,
			Size:     finfo.Size(),
		}

		// add to cache
		ftree.Files[fileName] = append(ftree.Files[fileName], fn)
		// add to this dir
		fn.SubFiles[fn.FileInfo.Name()] = fnn
		fn.Size += fnn.Size
		return fnn.Size
	} else {
		if fn.SubFiles[path[0]] == nil {
			fn.addDir(path[0:1])
		}
		increase := fn.SubFiles[path[0]].addFile(path[1:], fileName)
		fn.Size += increase
		return increase
	}
}

// 从子叶向根
func (fn *FileNode) rmFile(localFileName string) int64 {
	size := fn.Size

	cur := fn
	for cur.Parent != nil {
		cur.Parent.Size -= size
		cur = cur.Parent
	}

	fn.Parent.SubFiles[localFileName] = nil

	return size
}

type FileTree struct {
	Path  string                 // abs path
	Root  *FileNode              // for tree
	Files map[string][]*FileNode // for search , key is file name
	mutex sync.RWMutex
}

func CreateFileTree(localPath string) *FileTree {
	if !filepath.IsAbs(localPath) {
		var err error
		localPath, err = filepath.Abs(localPath)
		if err != nil {
			panic(err)
		}
	}

	tree := &FileTree{
		Path: localPath,
		Root: &FileNode{
			Path:     filepath.Dir(localPath),
			FileInfo: GetFileInfo(localPath),
		},
		Files: make(map[string][]*FileNode),
	}
	ftree = tree
	tree.Root.walk()
	return tree
}

// 从根向子叶递归
func (fn *FileNode) walk() int64 {
	// check file info
	if fn.FileInfo == nil {
		panic("file info get err.")
	}

	// add file to cache
	ftree.Files[fn.FileInfo.Name()] = append(ftree.Files[fn.FileInfo.Name()], fn)

	if fn.FileInfo.IsDir() {
		fn.SubFiles = map[string]*FileNode{}
		// get sub dirs
		for _, v := range getSubsFileList(fn.GetName()) {
			sub_fn := &FileNode{
				Path:     fn.GetName(),
				FileInfo: v,
				Parent:   fn,
			}
			fn.Size += sub_fn.walk()
			fn.SubFiles[v.Name()] = sub_fn
		}

	} else {
		fn.Size = fn.FileInfo.Size()
	}
	return fn.Size
}

// 从根向子叶递归
func (ft *FileTree) AddDir(localPath string) {
	ft.mutex.Lock()
	defer ft.mutex.Unlock()

	if Path(localPath).IsExist() {
		log.Printf("file %s not exist", localPath)
		return
	}

	// trans to abs path
	localPath = GetAbsPath(localPath)

	file, err := filepath.Rel(ft.Path, localPath)
	if err != nil {
		log.Printf("file %s trans to abs failed %s", file, err)
		return
	}

	// get path
	paths := strings.Split(file, "/")
	print(paths)

	ft.Root.addDir(paths)
}

// 从根向子叶递归
func (ft *FileTree) AddFile(localFile string) {
	ft.mutex.Lock()
	defer ft.mutex.Unlock()

	if Path(localFile).IsExist() {
		log.Printf("file %s not exist", localFile)
		return
	}

	// trans to abs path
	localFile = GetAbsPath(localFile)

	localFile, err := filepath.Rel(ft.Path, localFile)
	if err != nil {
		log.Printf("file %s trans to abs failed %s", localFile, err)
		return
	}

	// get path
	paths := strings.Split(localFile, "/")
	print(paths)

	ft.Root.addFile(paths[:len(paths)-1], paths[len(paths)-1])
}

func (ft *FileTree) RmFile(localFile string) {
	ft.mutex.Lock()
	defer ft.mutex.Unlock()

	// trans to abs path
	localFile = GetAbsPath(localFile)

	// get path
	fileName := filepath.Base(localFile)

	for _, v := range ft.Files[fileName] {
		if v.GetName() == localFile {
			log.Printf("%s\n", localFile)
			v.rmFile(fileName)
		}
	}
}

//
func (ft *FileTree) SearchFile(localFileName string) []*FileNode {
	var files = make([]*FileNode, 0)

	ft.mutex.RLock()
	defer ft.mutex.RUnlock()

	files = append(files, ft.Files[localFileName]...)

	return files
}

func (ft *FileTree) GetFile(localFileName string) *FileNode {
	ft.mutex.RLock()
	defer ft.mutex.RUnlock()

	for _, v := range ft.Files[filepath.Base(localFileName)] {
		if v.GetName() == localFileName {
			return v
		}
	}
	return nil
}
