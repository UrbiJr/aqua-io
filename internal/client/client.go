package client

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

// Client is used by packages that use/implement inernal/bot.Bot.
type Client struct {
	TLSClient *tls_client.HttpClient
}

func NewClient(timeout int, allowRedirects bool, proxy string, tlsClient tls_client.ClientProfile) (*Client, error) {

	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(timeout),
		tls_client.WithClientProfile(tlsClient),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		//tls_client.WithInsecureSkipVerify(),
	}

	if proxy != "" {
		proxyUrl, err := ValidateProxyFormat(proxy)
		if err == nil {
			options = append(options, tls_client.WithProxyUrl(proxyUrl.String()))
		}
	}

	if !allowRedirects {
		options = append(options, tls_client.WithNotFollowRedirects())
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		TLSClient: &client,
	}, nil
}
