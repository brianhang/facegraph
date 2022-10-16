package oauthgoogle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DISCOVERY_DOCUMENT_URL = "https://accounts.google.com/.well-known/openid-configuration"

type discoveryDocument struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

var lastDiscoveryDocument *discoveryDocument

func fetchAuthorizationEndpoint() (string, error) {
	document, err := fetchDiscoveryDocument()
	if err != nil {
		return "", err
	}
	return document.AuthorizationEndpoint, nil
}

func fetchTokenEndpoint() (string, error) {
	document, err := fetchDiscoveryDocument()
	if err != nil {
		return "", err
	}
	return document.TokenEndpoint, nil
}

func fetchDiscoveryDocument() (*discoveryDocument, error) {
	if lastDiscoveryDocument != nil {
		return lastDiscoveryDocument, nil
	}

	res, err := http.Get(DISCOVERY_DOCUMENT_URL)
	if err != nil {
		return lastDiscoveryDocument, fmt.Errorf("failed to request Google discovery document: %v", err)
	}
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return lastDiscoveryDocument, fmt.Errorf("failed to read Google discovery document: %v", err)
	}

	lastDiscoveryDocument := &discoveryDocument{}
	json.Unmarshal(rawBody, &lastDiscoveryDocument)
	return lastDiscoveryDocument, nil
}
