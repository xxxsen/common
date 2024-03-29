package s3

type S3Config struct {
	Endpoint  string `json:"endpoint"`
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	UseSSL    bool   `json:"use_ssl"`
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
}

type config struct {
	secretId  string
	secretKey string
	ssl       bool
	endpoint  string
	bucket    string
	region    string
}

type Option func(c *config)

func WithSecret(id, key string) Option {
	return func(c *config) {
		c.secretId = id
		c.secretKey = key
	}
}

func WithSSL(v bool) Option {
	return func(c *config) {
		c.ssl = v
	}
}

func WithEndpoint(ep string) Option {
	return func(c *config) {
		c.endpoint = ep
	}
}

func WithBucket(bk string) Option {
	return func(c *config) {
		c.bucket = bk
	}
}

func WithRegion(rg string) Option {
	return func(c *config) {
		c.region = rg
	}
}
