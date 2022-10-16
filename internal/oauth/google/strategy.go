package oauthgoogle

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Strategy struct{}

func (s *Strategy) GetAuthenticationURL(redirectURL *url.URL) (*url.URL, error) {
	endpoint, err := fetchAuthorizationEndpoint()
	if err != nil {
		return nil, err
	}

	authURL, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return authURL, err
	}

	q := authURL.Query()
	q.Set("response_type", "code")
	q.Set("scope", "openid")
	q.Set("redirect_uri", redirectURL.String())
	q.Set("client_id", fetchClientID())
	authURL.RawQuery = q.Encode()

	return authURL, nil
}

func (s *Strategy) HandleAuthenticationCallback(redirectURL *url.URL, w http.ResponseWriter, r *http.Request) error {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	code := r.URL.Query().Get("code")
	endpoint, err := fetchTokenEndpoint()
	if err != nil {
		return err
	}

	q, _ := url.ParseQuery("")
	q.Set("code", code)
	q.Set("client_id", fetchClientID())
	q.Set("client_secret", fetchClientSecret())
	q.Set("redirect_uri", redirectURL.String())
	q.Set("grant_type", "authorization_code")

	res, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(q.Encode()))
	if err != nil {
		return err
	}
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	io.WriteString(w, string(content))
	return nil
}
