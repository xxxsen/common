package es

import "time"

type EsConfig struct {
	User     string   `json:"user"`
	Password string   `json:"password"`
	Timeout  int      `json:"timeout"`
	Host     []string `json:"host"`
}

type config struct {
	user, password string
	urls           []string
	timeout        time.Duration
}

type Option func(c *config)

func WithAuth(u, p string) Option {
	return func(c *config) {
		c.user, c.password = u, p
	}
}

func WithHost(host ...string) Option {
	return func(c *config) {
		c.urls = host
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.timeout = timeout
	}
}
