package oauthgoogle

import "os"

var clientID string
var clientSecret string

func fetchClientID() string {
	if clientID != "" {
		return clientID
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	return clientID
}

func fetchClientSecret() string {
	if clientSecret != "" {
		return clientSecret
	}

	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	return clientSecret
}
