package cgi

type config struct {
	addresses  []string
	registerFn HandlerRegisterFunc
	attach     map[string]interface{}
}

type Option func(c *config)

func WithHandlerRegister(fn HandlerRegisterFunc) Option {
	return func(c *config) {
		c.registerFn = fn
	}
}

func WithAddress(address string) Option {
	return func(c *config) {
		c.addresses = append(c.addresses, address)
	}
}

func WithAttach(key string, val interface{}) Option {
	return func(c *config) {
		c.attach[key] = val
	}
}
