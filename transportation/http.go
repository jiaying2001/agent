package transportation

import (
	"github.com/go-resty/resty/v2"
	"github.com/jiaying2001/agent/log"
	"github.com/jiaying2001/agent/store"
)

var (
	HOST = "http://" + store.C.Server.Hostname + `:` + store.C.Server.Port
)

func validate(url string) {
	if store.Pass.AuthToken == "" {
		log.Logger.Error("Auth token is not set when try requesting " + url)
	}
}

func Post(url string, body []byte) []byte {
	if url != "/api/user/login" {
		validate(url)
	}
	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetAuthToken(store.Pass.AuthToken).
		SetBody(body).
		Post(HOST + url)
	return resp.Body()
}

func Get(url string, params map[string]string) []byte {
	validate(url)
	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetAuthToken(store.Pass.AuthToken).
		SetQueryParams(params).
		Get(HOST + url)
	return resp.Body()
}
