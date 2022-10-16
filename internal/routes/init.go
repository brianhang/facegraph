package routes

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"brianhang.me/facegraph/internal/oauth"
	oauthgoogle "brianhang.me/facegraph/internal/oauth/google"
)

var isSecure bool
var baseURL *url.URL

func getBaseURL() *url.URL {
	if baseURL != nil {
		return baseURL
	}

	rawURL := os.Getenv("WEBSERVER_BASE_URL")
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("The WEBSERVER_BASE_URL environment variable is not a valid URL: %v", err)
	}
	return baseURL
}

func getURLForPath(path string) *url.URL {
	return getBaseURL().JoinPath(path)
}

func setupRoutes() error {
	http.HandleFunc("/", handleHome)

	googleOAuth := &oauthgoogle.Strategy{}
	err := oauth.SetupRoutesForStrategy(
		googleOAuth,
		getURLForPath("/auth/google"),
		internalErrorResponse,
	)
	if err != nil {
		return fmt.Errorf("failed to set up routes: %v", err)
	}

	return nil
}

func Init() error {
	isSecure = false

	port, err := getPort()
	if err != nil {
		return err
	}
	addr := fmt.Sprintf(":%d", port)

	err = setupRoutes()
	if err != nil {
		return err
	}

	certFile := os.Getenv("WEBSERVER_TLS_CERT_FILE")
	keyFile := os.Getenv("WEBSERVER_TLS_KEY_FILE")

	if certFile != "" && keyFile != "" {
		isSecure = true
		log.Printf("Listening (with TLS) on port %d\n", port)

		return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	}

	log.Printf("Listening on port %d\n", port)
	return http.ListenAndServe(addr, nil)
}

func IsSecure() bool {
	return isSecure
}

func getPort() (int, error) {
	rawPort := os.Getenv("WEBSERVER_PORT")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return port, fmt.Errorf("\"%s\" is not a valid port number", rawPort)
	}
	return port, nil
}
