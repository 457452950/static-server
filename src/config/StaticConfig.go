package config

import (
	"os"
)

var (
	DefaultLocalHost = ""
	DefaultLocalPort = 80
	DefaultRootDir   = os.TempDir()
	DefaultTheme     = "black"
	DefaultPrefix    = ""
)
