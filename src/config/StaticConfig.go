package config

import (
	"os"
)

var (
	ConfigDefaultLocalHost          = ""
	ConfigDefaultLocalPort    int16 = 80
	ConfigDefaultLocalTLSPort int16 = 443
	ConfigDefaultRootDir            = os.TempDir()
	ConfigDefaultTheme              = "black"
	ConfigDefaultPrefix             = ""

	SysInfoVersion   = "unknown"
	SysInfoBuildTime = "unknown time"
	SysInfoGitCommit = "unknown git commit"
	SysInfoGitSite   = "https://github.com/457452950/static-server"

	PrefixSpecialSymbol = "/-/"
	PrefixSysInfo       = PrefixSpecialSymbol + "sysinfo"
	PrefixAssets        = PrefixSpecialSymbol + "assets"

	ConfigYamlFile = ".ghs.yml"
)
