package files

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsDir()
}

type Path string

func (p Path) join(path string) Path {
	return Path(filepath.Join(string(p), path))
}

func (p Path) IsExist() bool {
	_, err := os.Stat(string(p))
	if err != nil {
		return false
	}
	return true
}

type FileTransformer struct {
	basicPath Path
}

func (ftf FileTransformer) GetBasePath() Path {
	return ftf.basicPath
}

func CreateNewFileTransformer(basePath string) FileTransformer {
	if !isDir(basePath) {
		panic("base path is not a dir.")
	}
	ftf := FileTransformer{basicPath: Path(basePath)}
	return ftf
}

// IsolationPath local path to string, Hide the local directory
func (ftf FileTransformer) IsolationPath(path string) string {
	return ""
}

// TransformPath transform the string to local directory
func (ftf FileTransformer) TransformPath(path string) (Path, error) {
	if path == "" {
		return ftf.basicPath, nil
	}
	path = filepath.Clean(path)

	if strings.Contains(path, "..") {
		return "", errors.New("invalid path")
	}

	localPath := ftf.basicPath.join(path)
	if !localPath.IsExist() {
		return "", errors.New("file not found")
	}

	return localPath, nil
}
