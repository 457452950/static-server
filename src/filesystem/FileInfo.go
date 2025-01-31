package filesystem

import (
	"log"

	"github.com/shogo82148/androidbinary/apk"
)

var (
	FILE_TYPE_MARKDOWN        string = "markdown"
	FILE_TYPE_MARKDOWN_SUFFIX string = ".md"
	FILE_TYPE_APK             string = "apk"
	FILE_TYPE_APK_SUFFIX      string = ".apk"
	FILE_TYPE_TEXT            string = "text"
	FILE_TYPE_DIR             string = "dir"
	FILE_TYPE_FILE            string = "file"
)

type ApkInfo struct {
	PackageName  string `json:"packageName"`
	MainActivity string `json:"mainActivity"`
	Version      struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"version"`
}

// GetApkInfo path should be absolute
func GetApkInfo(path string) (ai *ApkInfo) {
	// catch panic
	defer func() {
		if err := recover(); err != nil {
			log.Println("parse-apk-info panic:", err)
		}
	}()

	// load apk file
	apkf, err := apk.OpenFile(path)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	ai = &ApkInfo{}
	ai.MainActivity, err = apkf.MainActivity()
	if err != nil {
		log.Println(err)
	}

	ai.PackageName = apkf.PackageName()
	ai.Version.Code = apkf.Manifest().VersionCode
	ai.Version.Name = apkf.Manifest().VersionName
	return
}
