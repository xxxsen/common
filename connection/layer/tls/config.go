package tls

type cliConfig struct {
	SNI                string `json:"sni"`
	SkipInsecureVerify bool   `json:"skip_insecure_verify"`
	FingerPrint        string `json:"fingerprint"`
	MinTLSVersion      string `json:"min_tls_version"`
	MaxTLSVersion      string `json:"max_tls_version"`
}

type svrConfig struct {
	CertFile      string `json:"cert_file"`
	KeyFile       string `json:"key_file"`
	MinTLSVersion string `json:"min_tls_version"`
	MaxTLSVersion string `json:"max_tls_version"`
}
