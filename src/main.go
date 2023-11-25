package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"text/template"

	"github.com/codeskyblue/go-accesslog"
	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"static-server/config"
	"static-server/logger"
)

// type Configure struct {
// 	Conf            *os.File `yaml:"-"`
// 	Addr            string   `yaml:"addr"`
// 	Port            int      `yaml:"port"`
// 	Root            string   `yaml:"root"`
// 	Prefix          string   `yaml:"prefix"`
// 	HTTPAuth        string   `yaml:"httpauth"`
// 	Cert            string   `yaml:"cert"`
// 	Key             string   `yaml:"key"`
// 	Cors            bool     `yaml:"cors"`
// 	Theme           string   `yaml:"theme"`
// 	XHeaders        bool     `yaml:"xheaders"`
// 	Upload          bool     `yaml:"upload"`
// 	Delete          bool     `yaml:"delete"`
// 	PlistProxy      string   `yaml:"plistproxy"`
// 	Title           string   `yaml:"title"`
// 	Debug           bool     `yaml:"debug"`
// 	GoogleTrackerID string   `yaml:"google-tracker-id"`
// 	Auth            struct {
// 		Type   string `yaml:"type"` // openid|http|github
// 		OpenID string `yaml:"openid"`
// 		HTTP   string `yaml:"http"`
// 		ID     string `yaml:"id"`     // for oauth2
// 		Secret string `yaml:"secret"` // for oauth2
// 	} `yaml:"auth"`
// }

var (
	// defaultPlistProxy = "https://plistproxy.herokuapp.com/plist"
	// defaultOpenID     = "https://login.netease.com/openid"
	// gcfg              = Configure{}

	VERSION   = "unknown"
	BUILDTIME = "unknown time"
	GITCOMMIT = "unknown git commit"
	SITE      = "https://github.com/457452950/static-server"
)

func versionMessage() string {
	t := template.Must(template.New("version").Parse(`GoHTTPServer
  Version:        {{.Version}}
  Go version:     {{.GoVersion}}
  OS/Arch:        {{.OSArch}}
  Git commit:     {{.GitCommit}}
  Built:          {{.Built}}
  Site:           {{.Site}}`))
	buf := bytes.NewBuffer(nil)
	t.Execute(buf, map[string]interface{}{
		"Version":   VERSION,
		"GoVersion": runtime.Version(),
		"OSArch":    runtime.GOOS + "/" + runtime.GOARCH,
		"GitCommit": GITCOMMIT,
		"Built":     BUILDTIME,
		"Site":      SITE,
	})
	return buf.String()
}

// func parseFlags() error {
// 	// initial default conf
// 	gcfg.Root = "./"
// 	gcfg.Port = 8000
// 	gcfg.Addr = ""
// 	gcfg.Theme = "black"
// 	gcfg.PlistProxy = defaultPlistProxy
// 	gcfg.Auth.OpenID = defaultOpenID
// 	gcfg.GoogleTrackerID = "UA-81205425-2"
// 	gcfg.Title = "Go HTTP File Server"

// 	kingpin.HelpFlag.Short('h')
// 	kingpin.Version(versionMessage())
// 	kingpin.Flag("conf", "config file path, yaml format").FileVar(&gcfg.Conf)
// 	kingpin.Flag("root", "root directory, default ./").Short('r').StringVar(&gcfg.Root)
// 	kingpin.Flag("prefix", "url prefix, eg /foo").StringVar(&gcfg.Prefix)
// 	kingpin.Flag("port", "listen port, default 8000").IntVar(&gcfg.Port)
// 	kingpin.Flag("addr", "listen address, eg 127.0.0.1:8000").Short('a').StringVar(&gcfg.Addr)
// 	kingpin.Flag("cert", "tls cert.pem path").StringVar(&gcfg.Cert)
// 	kingpin.Flag("key", "tls key.pem path").StringVar(&gcfg.Key)
// 	kingpin.Flag("auth-type", "Auth type <http|openid>").StringVar(&gcfg.Auth.Type)
// 	kingpin.Flag("auth-http", "HTTP basic auth (ex: user:pass)").StringVar(&gcfg.Auth.HTTP)
// 	kingpin.Flag("auth-openid", "OpenID auth identity url").StringVar(&gcfg.Auth.OpenID)
// 	kingpin.Flag("theme", "web theme, one of <black|green>").StringVar(&gcfg.Theme)
// 	kingpin.Flag("upload", "enable upload support").BoolVar(&gcfg.Upload)
// 	kingpin.Flag("delete", "enable delete support").BoolVar(&gcfg.Delete)
// 	kingpin.Flag("xheaders", "used when behide nginx").BoolVar(&gcfg.XHeaders)
// 	kingpin.Flag("cors", "enable cross-site HTTP request").BoolVar(&gcfg.Cors)
// 	kingpin.Flag("debug", "enable debug mode").BoolVar(&gcfg.Debug)
// 	kingpin.Flag("plistproxy", "plist proxy when server is not https").Short('p').StringVar(&gcfg.PlistProxy)
// 	kingpin.Flag("title", "server title").StringVar(&gcfg.Title)
// 	kingpin.Flag("google-tracker-id", "set to empty to disable it").StringVar(&gcfg.GoogleTrackerID)

// 	kingpin.Parse() // first parse conf

// 	if gcfg.Conf != nil {
// 		defer func() {
// 			kingpin.Parse() // command line priority high than conf
// 		}()
// 		ymlData, err := ioutil.ReadAll(gcfg.Conf)
// 		if err != nil {
// 			return err
// 		}
// 		return yaml.Unmarshal(ymlData, &gcfg)
// 	}
// 	return nil
// }

func main() {
	// init log seting
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// get config from file
	app_config := config.LoadFromFile("config.json")
	config.DumpConfig(app_config)

	ss := NewHTTPStaticServer(app_config.StSrvConf)

	// plist config
	// if gcfg.PlistProxy != "" {
	// 	u, err := url.Parse(gcfg.PlistProxy)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	u.Scheme = "https"
	// 	ss.PlistProxy = u.String()
	// }
	// if ss.PlistProxy != "" {
	// 	log.Printf("plistproxy: %s", strconv.Quote(ss.PlistProxy))
	// }

	var hdlr http.Handler = ss

	hdlr = accesslog.NewLoggingHandler(hdlr, logger.GetLogger())

	switch app_config.StSrvConf.Auth.Type {
	case "http":
		// HTTP Basic Authentication
		userpass := strings.SplitN(app_config.StSrvConf.Auth.HTTP, ":", 2)
		if len(userpass) == 2 {
			user, pass := userpass[0], userpass[1]
			hdlr = httpauth.SimpleBasicAuth(user, pass)(hdlr)
		}
	case "openid":
		handleOpenID(app_config.StSrvConf.Auth.OpenID, false) // FIXME(ssx): set secure default to false
		// case "github":
		// 	handleOAuth2ID(gcfg.Auth.Type, gcfg.Auth.ID, gcfg.Auth.Secret) // FIXME(ssx): set secure default to false
	case "oauth2-proxy":
		handleOauth2()
	}

	// CORS
	if app_config.Cors {
		hdlr = handlers.CORS()(hdlr)
	}
	if app_config.XHeaders {
		hdlr = handlers.ProxyHeaders(hdlr)
	}

	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/-/sysinfo", func(w http.ResponseWriter, r *http.Request) {
		data := versionMessage()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write([]byte(data))
	})
	mainRouter.PathPrefix("/-/assets/").Handler(http.StripPrefix("/-/", http.FileServer(Assets)))
	if app_config.StSrvConf.Prefix != "" {
		mainRouter.PathPrefix(app_config.StSrvConf.Prefix).Subrouter()
		// mainRouter.Handle(app_config.StSrvConf.Prefix, hdlr)
		mainRouter.PathPrefix(app_config.StSrvConf.Prefix).Handler(hdlr)
		mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, app_config.StSrvConf.Prefix, http.StatusTemporaryRedirect)
		})
	} else {
		mainRouter.PathPrefix("/").Handler(hdlr) // 坑：一般路由放后面
	}

	// get local bind address
	var local = app_config.Host
	if local == "" {
		local = getLocalIP()
	}
	localBind := fmt.Sprintf("%s:%d", app_config.Host, app_config.Port)
	log.Printf("local %s, address http://%s:%d\n", localBind, local, app_config.Port)

	srv := &http.Server{
		Handler: mainRouter,
		Addr:    localBind,
	}

	var err error
	if app_config.EnableSsl() {
		err = srv.ListenAndServeTLS(app_config.SSL.Cert, app_config.SSL.Key)
	} else {
		err = srv.ListenAndServe()
	}
	log.Fatal(err)
}
