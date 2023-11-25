package files

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

type Path string

func (p Path) Join(path string) Path {
	return Path(filepath.Join(string(p), path))
}

func (p Path) IsExist() bool {
	_, err := os.Stat(string(p))
	return err == nil
}

func (p Path) Get() string {
	return string(p)
}

type FileTransformer struct {
	prefix    string
	basicPath Path
}

func (ftf FileTransformer) GetBasePath() Path {
	return ftf.basicPath
}

func CreateFileTransformer(prefix string, basePath string) FileTransformer {
	if !IsDir(basePath) {
		panic("base path is not a dir.")
	}
	var err error

	if !filepath.IsAbs(basePath) {
		basePath, err = filepath.Abs(basePath)
		if err != nil {
			panic(err)
		}
	}

	ftf := FileTransformer{
		prefix:    prefix,
		basicPath: Path(basePath),
	}
	return ftf
}

// IsolationPath local path to string, Hide the local directory
func (ftf FileTransformer) IsolationPath(path string) (string, error) {
	if strings.HasPrefix(path, "..") {
		panic("invalid path")
	}

	var err error
	if filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	path, err = filepath.Rel(string(ftf.basicPath), path)
	if err != nil {
		return "", err
	}

	path = filepath.Join(ftf.prefix, path)

	path = strings.TrimPrefix(path, "/")

	return path, nil
}

// TransformPath transform the string to local directory
func (ftf FileTransformer) TransformPath(path string) (Path, error) {
	if path == ftf.prefix {
		return ftf.basicPath, nil
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = filepath.Clean(path)

	path, err := filepath.Rel(ftf.prefix, path)
	if err != nil {
		return "", err
	}

	if strings.Contains(path, "..") {
		return "", errors.New("invalid path")
	}

	localPath := ftf.basicPath.Join(path)
	if !localPath.IsExist() {
		return "", errors.New("file not found")
	}

	return localPath, nil
}
