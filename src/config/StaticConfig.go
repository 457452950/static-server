package config

import (
	"os"
)

var (
	ConfigDefaultLocalHost            = ""
	ConfigDefaultLocalPort    int16   = 80
	ConfigDefaultLocalTLSPort int16   = 443
	ConfigDefaultRootDir              = os.TempDir()
	ConfigDefaultTheme                = "black"
	ConfigDefaultPrefix               = ""
	ConfigMaxUploadFilesize   float64 = 10

	SysInfoVersion                 = "unknown"
	SysInfoBuildTime               = "unknown time"
	SysInfoGitCommit               = "unknown git commit"
	SysInfoGitSite                 = "https://github.com/457452950/static-server"
	SysInfoMaxUploadFilesize int64 = int64(ConfigMaxUploadFilesize * 1024)

	PrefixSpecialSymbol = "/-/"
	PrefixSysInfo       = PrefixSpecialSymbol + "sysinfo"
	PrefixAssets        = PrefixSpecialSymbol + "assets"

	ConfigYamlFile = ".ghs.yml"
)
