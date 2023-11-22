package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"static-server/files"
)

var (
	prefix_match_reg = "^(/.+)+$"
	reg_prefix_match = regexp.MustCompile(prefix_match_reg)
)

type AuthInfo struct {
	Type   string `json:"type"`
	HTTP   string `json:"http"`
	OpenID string `json:"openid"`
	ID     string `json:"id"`     // for oauth2
	Secret string `json:"secret"` // for oauth2
}

func (Auth AuthInfo) Check() bool {
	switch Auth.Type {
	case "http":
		userAndPass := strings.SplitN(Auth.HTTP, ":", 2)
		if len(userAndPass) != 2 {
			return false
		}
	case "openid":
		panic("not implemented")
	case "oauth2-proxy":
		panic("not implemented")
	}

	return true
}

func (Auth AuthInfo) GetUserAndPass() (string, string, error) {
	userpass := strings.SplitN(Auth.HTTP, ":", 2)
	if len(userpass) == 2 {
		return userpass[0], userpass[1], nil
	}
	return "", "", errors.New("invalid config")
}

type FileServiceConfig struct {
	Root            string   `json:"root"`   // 本地根路径
	Prefix          string   `json:"prefix"` // url prefix, eg /foo
	Upload          bool     `json:"upload"`
	Delete          bool     `json:"delete"`
	Title           string   `json:"title"` //
	Theme           string   `json:"theme"`
	Plistproxy      string   `json:"plistproxy"`
	GoogleTrackerID string   `json:"google-tracker-id"`
	Auth            AuthInfo `json:"auth"`
}

type Ssl struct {
	Enable bool   `json:"enable"`
	Cert   string `json:"cert"`
	Key    string `json:"key"`
}

type AppConfig struct {
	Host      string            `json:"host"`
	Port      int               `json:"port"`
	Httpauth  string            `json:"httpauth"`
	SSL       Ssl               `json:"ssl"`
	Cors      bool              `json:"cors"`
	XHeaders  bool              `json:"xheaders"`
	Debug     bool              `json:"debug"`
	StSrvConf FileServiceConfig `json:"ss"`
}

func GetDefaultConfig() (config AppConfig) {
	config = AppConfig{
		Host:     DefaultLocalHost,
		Port:     DefaultLocalPort,
		Httpauth: "",
		SSL: Ssl{
			Enable: false,
			Cert:   "",
			Key:    "",
		},
		Cors:     false,
		XHeaders: false,
		Debug:    false,
		StSrvConf: FileServiceConfig{
			Root:            DefaultRootDir,
			Prefix:          DefaultPrefix,
			Upload:          false,
			Delete:          false,
			Title:           "",
			Theme:           DefaultTheme,
			Plistproxy:      "",
			GoogleTrackerID: "",
			Auth: AuthInfo{
				Type:   "",
				HTTP:   "",
				OpenID: "",
				ID:     "",
				Secret: "",
			},
		},
	}

	return config
}

func LoadFromFile(file string) AppConfig {
	conf := GetDefaultConfig()

	json_file, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return conf
	}
	defer json_file.Close()

	file_data, err := ioutil.ReadAll(json_file)
	if err != nil {
		log.Println(err)
		return conf
	}

	err = json.Unmarshal(file_data, &conf)
	if err != nil {
		log.Println(err)
		return conf
	}

	return conf
}

func DumpConfig(conf AppConfig) {
	data, err := json.Marshal(conf)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(data))
}

func (conf FileServiceConfig) CheckPrefix() error {
	if conf.Prefix == "" {
		return nil
	}

	ok := reg_prefix_match.MatchString(conf.Prefix)
	if !ok {
		return errors.New("prefix match fail. usage: '/' '/app' '/a/b' ")
	}
	return nil
}

func (conf *FileServiceConfig) CheckRoot() error {
	conf.Root = filepath.ToSlash(filepath.Clean(conf.Root))
	if !strings.HasSuffix(conf.Root, "/") {
		conf.Root = conf.Root + "/"
	}
	log.Printf("local root path: %s\n", conf.Root)

	ok := files.IsDir(conf.Root)
	if !ok {
		return errors.New("file not exist. ")
	}

	return nil
}
