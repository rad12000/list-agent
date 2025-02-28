package github

import (
	"fmt"
	"golang.org/x/net/http2"
	"net/http"
	"net/url"
)

type apiTransport struct {
	accessToken string
}

func (t apiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	const baseUrlStr = "https://api.github.com"
	const authHeader = "Authorization"

	githubApiUrl, _ := url.Parse(baseUrlStr)

	if req.URL.Host == "" || req.URL.Host == githubApiUrl.Host {
		req.URL.Scheme = githubApiUrl.Scheme
		req.URL.Host = githubApiUrl.Host

		if t.accessToken != "" {
			req.Header.Set(authHeader, fmt.Sprintf("Bearer %s", t.accessToken))
		}
	}

	defaultTransport := new(http2.Transport)
	return defaultTransport.RoundTrip(req)
}

func newHTTPClient(accessToken string) *http.Client {
	return &http.Client{
		Transport: apiTransport{accessToken: accessToken},
	}
}
