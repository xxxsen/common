package httprewrite

type config struct {
	ForceUseProxy bool              `json:"force_use_proxy"`
	RewritePath   string            `json:"rewrite_path"`
	RewriteQuery  map[string]string `json:"rewrite_query"`
	RewriteHeader map[string]string `json:"rewrite_header"`
}
