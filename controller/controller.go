package controller

import (
	configpackage "miniwebserver/config"
	loggerpackage "miniwebserver/logger"
	"net/http"
)

type Controller interface {
	Root(w http.ResponseWriter, r *http.Request)
	Options(w http.ResponseWriter, r *http.Request)
	About(w http.ResponseWriter, r *http.Request)
}

type DefaultController struct {
	config   configpackage.Configuration
	logger   loggerpackage.Logger
	response HTTPResponse
}

func NewHTTPController(config configpackage.Configuration, logger loggerpackage.Logger, res HTTPResponse) (defaultController *DefaultController) {

	defaultController = &DefaultController{
		config:   config,
		logger:   logger,
		response: res,
	}

	return defaultController
}

func (u *DefaultController) Root(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)

	u.response.DefaultText(w, http.StatusOK, true, *u.config.GetConfiguration().Version)
	return
}

func (u *DefaultController) Options(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)

	return
}

func (u *DefaultController) About(w http.ResponseWriter, r *http.Request) {
	defer u.response.Panic(w, r)
	u.response.Default(w, http.StatusOK, true, *u.config.GetConfiguration().Version)
	return
}
