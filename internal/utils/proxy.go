package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// ParseProxy parses a colon delimited string that may contain both a proxy host and credentials.
// The port and credentials are optional.
// Example Colon Separated Values: "hostname:port:username:password"
func ParseProxy(csv string) (*url.URL, error) {
	validScheme := regexp.MustCompile(`^https?:\/\/`)
	if validScheme.MatchString(csv) {
		return url.Parse(csv)
	}

	csv += ":::"
	values := strings.Split(csv, ":")

	hostname := values[0]
	port := values[1]
	username := values[2]
	password := values[3]

	host := ""

	if hostname != "" && port != "" {
		host = hostname + ":" + port
	} else {
		host = hostname
	}

	if username != "" || password != "" {
		if username != "" && password != "" {
			return url.Parse("http://" + username + ":" + password + "@" + host)
		} else if username != "" {
			return url.Parse("http://" + username + "@" + host)
		} else {
			return url.Parse("http://" + password + "@" + host)
		}
	}

	return url.Parse("http://" + host)
}
