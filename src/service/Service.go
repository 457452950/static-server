package service

import (
	"net/http"
	"static-server/config"
	"static-server/logger"

	"github.com/codeskyblue/go-accesslog"
	"github.com/goji/httpauth"
)

type Service struct {
	appConfig config.AppConfig
}

func CreateServer(conf config.AppConfig) (srv *Service) {
	srv = &Service{
		appConfig: conf,
	}

	fileServiceHandler := CreateFileStaticService(conf.StSrvConf)

	var handler http.Handler = fileServiceHandler

	// init logger
	handler = srv.InitHttpLogger(handler)
	// init auth
	handler = srv.SetAuthentication(handler)

	return
}

func (srv *Service) Run() {

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
