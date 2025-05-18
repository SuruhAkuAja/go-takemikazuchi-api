package configs

import (
	"net/http"
	"time"
)

type HttpClient struct {
	BaseURL    *string
	apiKey     *string
	HTTPClient *http.Client
}

func NewHttpClient(apiKey, baseUrl *string) *HttpClient {
	return &HttpClient{
		BaseURL: baseUrl,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}
