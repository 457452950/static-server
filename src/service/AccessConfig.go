package service

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"static-server/config"
	"static-server/filesystem"

	"gopkg.in/yaml.v2"
)

type AccessTable struct {
	Regex string `yaml:"regex"`
	Allow bool   `yaml:"allow"`
}

type UserControl struct {
	Email string
	// Access bool
	Upload bool
	Delete bool
	Token  string
}

type AccessConf struct {
	Upload       bool          `yaml:"upload" json:"upload"`
	Delete       bool          `yaml:"delete" json:"delete"`
	Users        []UserControl `yaml:"users" json:"users"`
	AccessTables []AccessTable `yaml:"accessTables"`
}

// var reCache = make(map[string]*regexp.Regexp)

// func (c *AccessConf) canAccess(fileName string) bool {
// 	for _, table := range c.AccessTables {
// 		pattern, ok := reCache[table.Regex]
// 		if !ok {
// 			pattern, _ = regexp.Compile(table.Regex)
// 			reCache[table.Regex] = pattern
// 		}
// 		// skip wrong format regex
// 		if pattern == nil {
// 			continue
// 		}
// 		if pattern.MatchString(fileName) {
// 			return table.Allow
// 		}
// 	}
// 	return true
// }

// func (c *AccessConf) canDelete(r *http.Request) bool {
// 	session, err := store.Get(r, defaultSessionName)
// 	if err != nil {
// 		return c.Delete
// 	}
// 	val := session.Values["user"]
// 	if val == nil {
// 		return c.Delete
// 	}
// 	userInfo := val.(*UserInfo)
// 	for _, rule := range c.Users {
// 		if rule.Email == userInfo.Email {
// 			return rule.Delete
// 		}
// 	}
// 	return c.Delete
// }

// func (c *AccessConf) canUploadByToken(token string) bool {
// 	for _, rule := range c.Users {
// 		if rule.Token == token {
// 			return rule.Upload
// 		}
// 	}
// 	return c.Upload
// }

// todo: how to work?
// func (c *AccessConf) canUpload(r *http.Request) bool {
// 	token := r.FormValue("token")
// 	if token != "" {
// 		return c.canUploadByToken(token)
// 	}
// 	session, err := store.Get(r, defaultSessionName)
// 	if err != nil {
// 		return c.Upload
// 	}
// 	val := session.Values["user"]
// 	if val == nil {
// 		return c.Upload
// 	}
// 	userInfo := val.(*UserInfo)

// 	for _, rule := range c.Users {
// 		if rule.Email == userInfo.Email {
// 			return rule.Upload
// 		}
// 	}
// 	return c.Upload
// }

func (handler *FileServiceHandler) defaultAccessConf() AccessConf {
	return AccessConf{
		Upload: handler.Config.Upload,
		Delete: handler.Config.Delete,
	}
}

// localPath 绝对路径
func (handler *FileServiceHandler) readAccessConf(localPath string) AccessConf {
	if localPath == "" {
		panic("load access config error, local path is ''")
	}

	relativePath, err := filepath.Rel(handler.fileTransformer.GetBasePath().Get(), localPath)

	if relativePath == ".." || err != nil {
		return handler.defaultAccessConf()
	}

	if filesystem.IsFile(localPath) {
		localPath = filepath.Dir(localPath)
	}

	cfgFile := filesystem.Path(localPath).Join(config.ConfigYamlFile)
	if cfgFile.IsExist() {
		data, err := ioutil.ReadFile(cfgFile.Get())
		if err != nil {
			if os.IsNotExist(err) {
				log.Println(err)
				return handler.defaultAccessConf()
			}
			log.Printf("Err read .ghs.yml: %v", err)
		}
		var ac AccessConf
		err = yaml.Unmarshal(data, &ac)
		if err != nil {
			log.Printf("Err format .ghs.yml: %v", err)
		}
		return ac
	} else {
		return handler.readAccessConf(filepath.Dir(localPath))
	}
}
