package versionmgr

type config struct {
	tabName string
}

type Option func(c *config)

func WithTabName(t string) Option {
	return func(c *config) {
		c.tabName = t
	}
}
