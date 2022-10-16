package appurl

import (
	"log"
	"net/url"
	"os"
)

var baseURL *url.URL

func Base() *url.URL {
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

func ForPath(path string) *url.URL {
	return Base().JoinPath(path)
}
