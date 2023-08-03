package client

import (
	"github.com/UrbiJr/aqua-io/internal/captcha"
	tls_client "github.com/bogdanfinn/tls-client"
)

// TLSClient is used to make requests with custom TLS profiles.
type TLSClient struct {
	captcha.Solver
	tls_client.HttpClient
}

type ClientOptions struct {
	Timeout          int
	AllowRedirects   bool
	CharlesProxy     bool
	Proxy            string
	TlsClientProfile tls_client.ClientProfile
}

func NewTLSClient(captchaOptions *captcha.SolverOptions, clientOptions *ClientOptions) (*TLSClient, error) {

	solver, err := captcha.NewCaptchaSolver(*captchaOptions)
	if err != nil {
		return nil, err
	}

	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithTimeoutSeconds(clientOptions.Timeout),
		tls_client.WithClientProfile(clientOptions.TlsClientProfile),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		//tls_client.WithInsecureSkipVerify(),
	}

	if clientOptions.CharlesProxy {
		options = append(options, tls_client.WithCharlesProxy("127.0.0.1", "8888"))
	} else if clientOptions.Proxy != "" {
		proxyUrl, err := ValidateProxyFormatToUrl(clientOptions.Proxy)
		if err == nil {
			options = append(options, tls_client.WithProxyUrl(proxyUrl.String()))
		}
	}

	if !clientOptions.AllowRedirects {
		options = append(options, tls_client.WithNotFollowRedirects())
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &TLSClient{
		Solver:     solver,
		HttpClient: client,
	}, nil
}
