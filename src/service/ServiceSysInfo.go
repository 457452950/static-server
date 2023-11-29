package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"static-server/config"
)

var (
	Version           = config.SysInfoVersion
	GoVersion         = runtime.Version()
	OSArch            = runtime.GOOS + "/" + runtime.GOARCH
	GitCommit         = config.SysInfoGitCommit
	Built             = config.SysInfoBuildTime
	Site              = config.SysInfoGitSite
	MaxUploadFilesize = config.SysInfoMaxUploadFilesize
)

type ServerSysInfo struct {
	Version           string `json:"version"`
	GoVersion         string `json:"goVersion"`
	OsArch            string `json:"osArch"`
	GitCommit         string `json:"gitCommit"`
	Built             string `json:"built"`
	Site              string `json:"site"`
	MaxUploadFilesize int64  `json:"maxFileSize"`
}

func versionMessage() []byte {
	info := ServerSysInfo{
		Version:           Version,
		GoVersion:         GoVersion,
		OsArch:            OSArch,
		GitCommit:         GitCommit,
		Built:             Built,
		Site:              Site,
		MaxUploadFilesize: MaxUploadFilesize,
	}
	data, err := json.Marshal(info)
	if err != nil {
		return nil
	}

	return data
}

func handleSysInfo(w http.ResponseWriter, r *http.Request) {
	data := versionMessage()
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Write([]byte(data))
}
