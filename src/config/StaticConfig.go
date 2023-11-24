package config

import (
	"os"
)

var (
	ConfigDefaultLocalHost = ""
	ConfigDefaultLocalPort = 80
	ConfigDefaultRootDir   = os.TempDir()
	ConfigDefaultTheme     = "black"
	ConfigDefaultPrefix    = ""

	SysInfoVersion   = "unknown"
	SysInfoBuildTime = "unknown time"
	SysInfoGitCommit = "unknown git commit"
	SysInfoGitSite   = "https://github.com/457452950/static-server"

	PrefixSpecialSymbol = "/-/"
	PrefixSysInfo       = PrefixSpecialSymbol + "sysinfo"
	PrefixAssets        = PrefixSpecialSymbol + "assets"
)
