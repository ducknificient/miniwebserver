package server

import (
	"context"
	"fmt"
	"miniwebserver/logger"
	"net/http"
)

type DefaultHTTPServer struct {
	ctx    context.Context
	logger logger.Logger
	server *http.Server
}

func NewHTTPServer(httpserver *http.Server, logger logger.Logger) *DefaultHTTPServer {
	return &DefaultHTTPServer{
		server: httpserver,
		logger: logger,
	}
}

func (s *DefaultHTTPServer) Run() (err error) {

	go func() {
		if err = s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal(fmt.Sprintf("Failed to initialize server: %v\n", err.Error()))
		}
	}()

	s.logger.Info("App listening on " + s.server.Addr)

	return err
}

func (s *DefaultHTTPServer) Shutdown(ctx context.Context) (err error) {

	err = s.server.Shutdown(s.ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Server forced to shutdown : %v\n", err.Error()))
		return err
	}

	return err
}
