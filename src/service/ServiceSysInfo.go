package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime"

	"static-server/config"
)

func versionMessage() string {
	t := template.Must(
		template.New("version").Parse(`GoHTTPServer
  			Version:        {{.Version}}
  			Go version:     {{.GoVersion}}
  			OS/Arch:        {{.OSArch}}
  			Git commit:     {{.GitCommit}}
  			Built:          {{.Built}}
  			Site:           {{.Site}}`))
	buf := bytes.NewBuffer(nil)
	t.Execute(buf, map[string]interface{}{
		"Version":   config.SysInfoVersion,
		"GoVersion": runtime.Version(),
		"OSArch":    runtime.GOOS + "/" + runtime.GOARCH,
		"GitCommit": config.SysInfoGitCommit,
		"Built":     config.SysInfoBuildTime,
		"Site":      config.SysInfoGitSite,
	})
	return buf.String()
}

func handleSysInfo(w http.ResponseWriter, r *http.Request) {
	data := versionMessage()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Write([]byte(data))
}
