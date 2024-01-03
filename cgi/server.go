package cgi

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var defaultServer *Server

func Init(opts ...Option) error {
	svc, err := NewServer(opts...)
	if err != nil {
		return err
	}
	defaultServer = svc
	return nil
}

func MustInit(opts ...Option) {
	if err := Init(opts...); err != nil {
		panic(err)
	}
}

func Run() error {
	return defaultServer.Run()
}

type Server struct {
	c *config
}

func NewServer(opts ...Option) (*Server, error) {
	c := &config{
		attach: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(c)
	}
	s := &Server{c: c}
	if err := s.initServer(); err != nil {
		return nil, fmt.Errorf("init server fail, err:%w", err)
	}
	return s, nil
}

func (s *Server) initServer() error {
	if len(s.c.addresses) == 0 {
		return fmt.Errorf("no bind address found")
	}
	return nil
}

func (s *Server) Run() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	s.registDefault(engine)
	s.c.registerFn(engine)
	if err := engine.Run(s.c.addresses...); err != nil {
		return err
	}
	return nil
}

func (s *Server) registDefault(engine *gin.Engine) {
	engine.Use(
		PanicRecoverMiddleware(s),
		EnableServerTraceMiddleware(s),
		EnableAttachMiddleware(s),
	)
}
