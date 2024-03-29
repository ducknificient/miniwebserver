package test

import (
	"fmt"
	configpackage "miniwebserver/config"
	controllerpackage "miniwebserver/controller"
	loggerpackage "miniwebserver/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {

	var (
		configPath string
		config     *configpackage.AppConfiguration
		err        error
	)

	// reading from command line
	var args = []string{".", "config.json"}
	configPath = args[1]

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// config file path
	config, err = configpackage.NewConfiguration(configPath)
	if err != nil {
		panic(err)
		return
	}

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
	httpresponse := controllerpackage.NewHTTPResponse(logger)

	// init default http controller
	controller := controllerpackage.NewHTTPController(config, logger, httpresponse)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Root)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := *config.Version
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAboutHandler(t *testing.T) {

	var (
		configPath string
		config     *configpackage.AppConfiguration
		err        error
	)

	// reading from command line
	var args = []string{".", "config.json"}
	configPath = args[1]

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/about/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	// config file path
	config, err = configpackage.NewConfiguration(configPath)
	if err != nil {
		panic(err)
		return
	}

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
	httpresponse := controllerpackage.NewHTTPResponse(logger)

	// init default http controller
	controller := controllerpackage.NewHTTPController(config, logger, httpresponse)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.Root)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := *config.Version
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
