package naivesvr

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logger"
)

var defaultLogger = logger.Logger()

type Server struct {
	c *Config
}

func NewServer(opts ...Option) (*Server, error) {
	c := &Config{
		attach: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.l == nil {
		c.l = defaultLogger
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
		SupportAttachMiddleware(s),
		EnableServerTraceMiddleware(s),
		SupportServerGetterMiddleware(s),
	)
}
