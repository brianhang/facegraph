package oauth

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Strategy interface {
	GetAuthenticationURL(redirectURL *url.URL) (*url.URL, error)
	HandleAuthenticationCallback(redirectURL *url.URL, w http.ResponseWriter, r *http.Request) error
}

type errorHandler func(w http.ResponseWriter, r *http.Request, err error)

func SetupRoutesForStrategy(strategy Strategy, baseURL *url.URL, handleError errorHandler) error {
	path := pathWithSlashes(baseURL.Path)
	callbackURL := baseURL.JoinPath("callback")
	callbackPath := pathWithSlashes(callbackURL.Path)

	authURL, err := strategy.GetAuthenticationURL(callbackURL)
	if err != nil {
		return err
	}
	authURLStr := authURL.String()

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, authURLStr, http.StatusTemporaryRedirect)
	})
	log.Printf("Setting up Google authentication at %s\n", path)

	http.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		err := strategy.HandleAuthenticationCallback(callbackURL, w, r)
		if err != nil {
			handleError(w, r, err)
		}
	})
	log.Printf("Setting up Google authentication callback at %s\n", callbackPath)

	return nil
}

func pathWithSlashes(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}
