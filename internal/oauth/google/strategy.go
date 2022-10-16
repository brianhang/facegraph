package oauthgoogle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"brianhang.me/facegraph/internal/oauth"
	"github.com/golang-jwt/jwt/v4"
)

type Strategy struct{}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token"`
}

type googleClaims struct {
	jwt.StandardClaims
}

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

func (s *Strategy) HandleAuthenticationCallback(
	redirectURL *url.URL,
	w http.ResponseWriter,
	r *http.Request,
	handleAuthenticated oauth.AuthenticatedHandler,
) error {
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

	tokenRes := tokenResponse{}
	if err = json.Unmarshal(content, &tokenRes); err != nil {
		return fmt.Errorf("failed to parse token response (%s): %v", content, err)
	}

	claims, err := fetchClaimsFromJWT(tokenRes.IDToken)
	if err != nil {
		return fmt.Errorf("failed to parse JWT (%s): %v", tokenRes.IDToken, err)
	}

	err = handleAuthenticated(w, r, claims.Subject)
	if err != nil {
		return err
	}

	return nil
}

func fetchClaimsFromJWT(token string) (googleClaims, error) {
	claims := googleClaims{}
	jwksURL, err := fetchJWKSURI()
	if err != nil {
		return claims, fmt.Errorf("failed to fetch JWKS URL: %v", err)
	}

	jwks, err := oauth.FetchJWKS(jwksURL)
	if err != nil {
		return claims, fmt.Errorf("failed to fetch JWKS: %v", err)
	}

	parsed, err := jwt.ParseWithClaims(
		token,
		&claims,
		jwks.Keyfunc,
	)
	if err != nil || !parsed.Valid {
		return claims, fmt.Errorf("failed to parse JWT (%s): %v", token, err)
	}

	return claims, nil
}
