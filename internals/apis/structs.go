package apis

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

type Client struct {
	Mail_Token   string
	User_UUID    string
	User_Id      string
	Login_Id     string
	User_Country string
	Promo        string
	CapKey       string
	HTTPClient   *tls_client.HttpClient
}
