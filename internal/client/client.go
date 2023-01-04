package client

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/cactus-aio/go-cactus/internal/captcha"
)

// Client is used by packages that use/implement inernal/bot.Bot.
type Client struct {
	captcha.Solver
	TLSClient *tls_client.HttpClient
}

type ClientOptions struct {
	timeout        int
	allowRedirects bool
	proxy          string
	tlsClient      tls_client.ClientProfile
}

func NewClient(captchaOptions *captcha.SolverOptions, clientOptions *ClientOptions) (*Client, error) {

	solver, err := captcha.NewCaptchaSolver(*captchaOptions)
	if err != nil {
		return nil, err
	}

	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(clientOptions.timeout),
		tls_client.WithClientProfile(clientOptions.tlsClient),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		//tls_client.WithInsecureSkipVerify(),
	}

	if clientOptions.proxy != "" {
		proxyUrl, err := ValidateProxyFormat(clientOptions.proxy)
		if err == nil {
			options = append(options, tls_client.WithProxyUrl(proxyUrl.String()))
		}
	}

	if !clientOptions.allowRedirects {
		options = append(options, tls_client.WithNotFollowRedirects())
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Solver:    solver,
		TLSClient: &client,
	}, nil
}
