package httpupgrade

type cliConfig struct {
	Host            string `json:"host"`
	Path            string `json:"path"`
	PaddingMin      uint   `json:"padding_min"`
	PaddingMax      uint   `json:"padding_max"`
	UpgradeProtocol string `json:"upgrade_protocol"`
}

type svrConfig struct {
	Path            string `json:"path"`
	PaddingMin      uint   `json:"padding_min"`
	PaddingMax      uint   `json:"padding_max"`
	FailCode        int    `json:"fail_code"`
	FailReason      string `json:"fail_reason"`
	UpgradeProtocol string `json:"upgrade_protocol"`
}
