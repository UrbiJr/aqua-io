package tls

import (
	"github.com/bogdanfinn/tls-client"
)

func NewProxyLess() tls_client.HttpClient {
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(),
		tls_client.WithClientProfile(tls_client.Chrome_110),
		tls_client.WithTimeoutSeconds(10),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify())
	if err != nil {
		return nil
	}
	return client
}
