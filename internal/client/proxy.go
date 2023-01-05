package client

import (
	"errors"
	"net/url"
	"regexp"
)

// ValidateProxyFormatToUrl takes a proxy string as input and returns it as *url.URL if its format is valid, error otherwise
func ValidateProxyFormatToUrl(proxy string) (*url.URL, error) {
	// Crea una regex per il formato proxy "ip:port:user:pass" e "ip:port".
	r, err := regexp.Compile(`^([^:]+):(\d+)(:([^:]+):([^@]+))?$`)
	if err != nil {
		return nil, err
	}

	// Verifica se il formato della stringa proxy corrisponde alla regex.
	if !r.MatchString(proxy) {
		return nil, errors.New("invalid proxy format")
	}

	// Crea un nuovo URL di tipo "http" utilizzando la stringa proxy.
	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		return nil, err
	}

	return proxyURL, nil
}

// ValidateProxyFormat takes a proxy string as input and returns true if its format is valid, false otherwise
func ValidateProxyFormat(proxy string) bool {
	// Crea una regex per il formato proxy "ip:port:user:pass" e "ip:port".
	r, err := regexp.Compile(`^([^:]+):(\d+)(:([^:]+):([^@]+))?$`)
	if err != nil {
		return false
	}

	// Verifica se il formato della stringa proxy corrisponde alla regex.
	if !r.MatchString(proxy) {
		return false
	}

	return true
}
