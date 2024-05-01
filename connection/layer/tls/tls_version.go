package tls

import tls "github.com/refraction-networking/utls"

const defaultTlsVersion = tls.VersionTLS13

var versionMapping = map[string]uint16{
	"1.1": tls.VersionTLS11,
	"1.2": tls.VersionTLS12,
	"1.3": tls.VersionTLS13,
}

func GetVersionIdByName(name string) (uint16, bool) {
	if v, ok := versionMapping[name]; ok {
		return v, true
	}
	return 0, false
}

func GetVersionIdByNameOrDefault(name string) uint16 {
	if v, ok := GetVersionIdByName(name); ok {
		return v
	}
	return defaultTlsVersion
}
