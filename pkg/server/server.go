package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func Get() *Server {
	return &Server{
		srv: &http.Server{
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) WithAddr(addr string) *Server {
	s.srv.Addr = fmt.Sprintf("0.0.0.0:%s", addr)
	return s
}

func (s *Server) WithErrLogger(l *log.Logger) *Server {
	s.srv.ErrorLog = l
	return s
}

func (s *Server) WithRouter(router http.Handler) *Server {
	s.srv.Handler = router
	return s
}

func (s *Server) Start() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("server missing address")
	}

	if s.srv.Handler == nil {
		return errors.New("server missing handler")
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Close() error {
	return s.srv.Close()
}
