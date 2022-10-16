package oauth

import (
	"net/http"
	"net/url"
)

type Strategy interface {
	GetAuthenticationURL(redirectURL *url.URL) (*url.URL, error)
	HandleAuthenticationCallback(redirectURL *url.URL, w http.ResponseWriter, r *http.Request) error
}

type errorHandler func(w http.ResponseWriter, r *http.Request, err error)

func SetupRoutesForStrategy(strategy Strategy, baseURL *url.URL, handleError errorHandler) error {
	authURL, err := strategy.GetAuthenticationURL(baseURL)
	if err != nil {
		return err
	}

	authURLStr := authURL.String()
	http.HandleFunc("/"+baseURL.Path, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, authURLStr, http.StatusTemporaryRedirect)
	})

	callbackURL := authURL.JoinPath("/callback")
	http.HandleFunc("/"+callbackURL.Path, func(w http.ResponseWriter, r *http.Request) {
		err := strategy.HandleAuthenticationCallback(baseURL, w, r)
		if err != nil {
			handleError(w, r, err)
		}
	})

	return nil
}
