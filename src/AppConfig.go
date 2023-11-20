package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type StaticServerConfig struct {
	Root            string `json:"root"`   // 本地根路径
	Prefix          string `json:"prefix"` // url prefix, eg /foo
	Upload          bool   `json:"upload"`
	Delete          bool   `json:"delete"`
	Title           string `json:"title"` //
	Theme           string `json:"theme"`
	Plistproxy      string `json:"plistproxy"`
	GoogleTrackerID string `json:"google-tracker-id"`
	Auth            struct {
		Type   string `json:"type"`
		HTTP   string `json:"http"`
		Openid string `json:"openid"`
		ID     string `json:"id"`     // for oauth2
		Secret string `json:"secret"` // for oauth2
	} `json:"auth"`
}

type AppConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Httpauth string `json:"httpauth"`
	SSL      struct {
		Enable bool   `json:"enable"`
		Cert   string `json:"cert"`
		Key    string `json:"key"`
	} `json:"ssl"`
	Cors      bool               `json:"cors"`
	Xheaders  bool               `json:"xheaders"`
	Debug     bool               `json:"debug"`
	StSrvConf StaticServerConfig `json:"ss"`
}

func GetDefaultConfig() AppConfig {
	config := AppConfig{}

	config.Port = 80
	config.Debug = true

	config.StSrvConf.Root = "/tmp"
	config.StSrvConf.Prefix = ""
	config.StSrvConf.Upload = false
	config.StSrvConf.Delete = false
	config.StSrvConf.Theme = "black"

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
