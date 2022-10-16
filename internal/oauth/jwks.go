package oauth

import (
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
)

var instances map[string]*keyfunc.JWKS

func FetchJWKS(jwksURL string) (*keyfunc.JWKS, error) {
	if instances == nil {
		instances = make(map[string]*keyfunc.JWKS, 1)
	}

	if instance, ok := instances[jwksURL]; ok {
		return instance, nil
	}

	// Create the keyfunc options. Use an error handler that logs. Refresh the JWKS when a JWT signed by an unknown KID
	// is found or at the specified interval. Rate limit these refreshes. Timeout the initial JWKS refresh request after
	// 10 seconds. This timeout is also used to create the initial context.Context for keyfunc.Get.
	options := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return jwks, err
	}

	instances[jwksURL] = jwks
	return jwks, nil
}
