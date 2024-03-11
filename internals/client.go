package client

import (
	"alien/internals/apis"
	"alien/internals/logger"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func NewClient(proxi, capkey string) (*apis.Client, error) {
	logger.DebugLogger.Println("Creating New Client")
	cuki_jar := tls_client.NewCookieJar()
	opts := []tls_client.HttpClientOption{
		tls_client.WithProxyUrl("http://" + proxi),
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithCookieJar(cuki_jar),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithRandomTLSExtensionOrder(),
	}
	http_client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), opts...)
	if err != nil {
		logger.ErrorLogger.Fatal("http_ClientErr: ", err)
		return nil, err
	}
	return &apis.Client{
		CapKey:     capkey,
		HTTPClient: &http_client,
	}, nil
}
