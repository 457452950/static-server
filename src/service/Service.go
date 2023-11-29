package service

import (
	"fmt"
	"log"
	"net/http"
	"static-server/config"
	"static-server/logger"

	"github.com/codeskyblue/go-accesslog"
	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Service struct {
	appConfig  config.AppConfig
	service    *http.Server
	tlsService *http.Server
	assets     AssetsHandler
}

func CreateServer(conf config.AppConfig) (srv *Service) {
	// init sysinfo config
	MaxUploadFilesize = int64(conf.MaxUpload * 1024)

	srv = &Service{
		appConfig: conf,
	}

	return
}

func (srv *Service) Init(assets http.FileSystem) {
	srv.assets.Set(assets)

	fileServiceHandler := createFileStaticService(srv.appConfig.StSrvConf, assets)

	var handler http.Handler = fileServiceHandler

	// init logger
	handler = srv.InitHttpLogger(handler)
	// init auth
	handler = srv.SetAuthentication(handler)

	handler = srv.SetCORS(handler)

	handler = srv.SetXHeaders(handler)

	handler = srv.SetPrefixAndHandler(handler)

	localBinded := srv.GetLocalBinded()
	log.Printf("%s\n", localBinded)

	var local = srv.appConfig.Host
	if local == "" {
		local = getLocalIP()
	}
	log.Printf("address http://%s:%d\n", local, srv.appConfig.Port)

	srv.service = &http.Server{
		Handler: handler,
		Addr:    localBinded,
	}

	if srv.appConfig.EnableSsl() {
		localBinded := srv.GetTLSLocalBinded()
		log.Printf("%s\n", localBinded)
		log.Printf("address https://%s:%d\n", local, srv.appConfig.SSL.Port)
		srv.tlsService = &http.Server{
			Handler: handler,
			Addr:    localBinded,
		}
	}
}

func (srv *Service) InitHttpLogger(handler http.Handler) http.Handler {
	return accesslog.NewLoggingHandler(handler, logger.GetLogger())
}

func (srv *Service) SetAuthentication(handler http.Handler) http.Handler {
	if !srv.appConfig.StSrvConf.Auth.Check() {
		panic("invalid config")
	}

	switch srv.appConfig.StSrvConf.Auth.Type {
	case "http":
		user, pass, err := srv.appConfig.StSrvConf.Auth.GetUserAndPass()
		if err != nil {
			panic("invalid config")
		}
		handler = httpauth.SimpleBasicAuth(user, pass)(handler)
		// case "openid":
		// 	handleOpenID(app_config.StSrvConf.Auth.OpenID, false) // FIXME(ssx): set secure default to false
		// 	// case "github":
		// 	// 	handleOAuth2ID(gcfg.Auth.Type, gcfg.Auth.ID, gcfg.Auth.Secret) // FIXME(ssx): set secure default to false
		// case "oauth2-proxy":
		// 	handleOauth2()
	}

	return handler
}

func (srv *Service) SetCORS(handler http.Handler) http.Handler {
	if srv.appConfig.Cors {
		return handlers.CORS()(handler)
	}
	return handler
}

func (srv *Service) SetXHeaders(handler http.Handler) http.Handler {
	if srv.appConfig.XHeaders {
		return handlers.ProxyHeaders(handler)
	}
	return handler
}

func (srv *Service) SetPrefixAndHandler(handler http.Handler) http.Handler {
	mainRouter := mux.NewRouter()

	// set sysinfo router
	mainRouter.HandleFunc(config.PrefixSysInfo, handleSysInfo) // 路径不带有隐式通配符
	mainRouter.PathPrefix(config.PrefixAssets).Handler(http.StripPrefix(config.PrefixSpecialSymbol, srv.assets.Get()))

	if srv.appConfig.StSrvConf.Prefix != "" {
		mainRouter.PathPrefix(srv.appConfig.StSrvConf.Prefix).Handler(handler)
		// 重定向
		mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, srv.appConfig.StSrvConf.Prefix, http.StatusTemporaryRedirect)
		})
	} else {
		mainRouter.PathPrefix("/").Handler(handler)
	}

	return mainRouter
}

func (srv *Service) GetLocalBinded() string {
	return fmt.Sprintf("%s:%d", srv.appConfig.Host, srv.appConfig.Port)
}

func (srv *Service) GetTLSLocalBinded() string {
	return fmt.Sprintf("%s:%d", srv.appConfig.Host, srv.appConfig.SSL.Port)
}

func (srv *Service) runSsl() {
	err := srv.tlsService.ListenAndServeTLS(srv.appConfig.SSL.Cert, srv.appConfig.SSL.Key)
	log.Println(err)
}

func (srv *Service) Run() error {
	if srv.appConfig.EnableSsl() {
		go srv.runSsl()
	}
	return srv.service.ListenAndServe()
}
