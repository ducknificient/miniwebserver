package main

import (
	"context"
	"fmt"
	configpackage "miniwebserver/config"
	"miniwebserver/controller"
	loggerpackage "miniwebserver/logger"
	"miniwebserver/router"
	serverpackage "miniwebserver/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configPath string
	config     *configpackage.AppConfiguration
	err        error
)

func init() {

	// reading from command line
	var args = os.Args
	configPath = args[1]

	// config file path
	config, err = configpackage.NewConfiguration(configPath)
	if err != nil {
		panic(err)
		return
	}
}

func main() {

	/*
		SETUP CONTEXT
	*/
	// https://dasarpemrogramangolang.novalagung.com/A-pipeline-context-cancellation.html
	ctx := context.Background()

	/*
		SETUP LOGGER
	*/
	logger, err := loggerpackage.NewLogger(config)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to init logger: %v", err.Error()))
	}

	/*
		SETUP CONTROLLER
	*/

	// init controller response
	httpresponse := controller.NewHTTPResponse(logger)

	// init default http controller
	handler := controller.NewHTTPController(config, logger, httpresponse)

	/*
		SETUP ROUTER
	*/

	// init router
	httprouter := router.NewRouter(config, handler)

	/*
		SETUP SERVER
	*/

	// init http server
	appIp := *config.AppIP + ":" + *config.AppPort
	httpserver := serverpackage.NewHTTPServer(&http.Server{
		Addr:    appIp,
		Handler: httprouter.Router,
	}, logger)

	// run server
	httpserver.Run()

	/*
		Graceful shutdown
	*/

	// Wait for kill signal of channel
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Shutdown server
	logger.Info("Shutting down server...")
	err = httpserver.Shutdown(ctx)
	if err != nil {
		// fmt.Errorf(fmt.Sprintf("Server forced to shutdown : %v\n", err.Error()))
		logger.Fatal(fmt.Sprintf("Server forced to shutdown : %v\n", err.Error()))
	}

}
