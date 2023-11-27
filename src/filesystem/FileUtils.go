package filesystem

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GetFileInfo(fileName string) fs.FileInfo {
	fi, err := os.Stat(fileName)
	if err != nil {
		log.Printf("%s\n", err)
		return nil
	}
	return fi
}

func getSubsFileList(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("%s", err)
		return nil
	}
	return files
}

func GetAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			log.Printf("path %s trans to abs failed %s", path, err)
			return ""
		}
	}
	return path
}
