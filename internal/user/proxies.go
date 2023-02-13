package user

type ProxyProfile struct {
	Id      string
	Name    string
	Proxies []string
}

// ReadProxies reads proxies from the DB and returns read data as []Profile
func ReadProxies() ([]ProxyProfile, error) {

	var proxies []ProxyProfile

	return proxies, nil

}

// WriteProxies writes proxies to the DB+
func WriteProxies(proxies []ProxyProfile) {

}
