package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"brianhang.me/facegraph/internal/api"
	"brianhang.me/facegraph/internal/appurl"
	"brianhang.me/facegraph/internal/oauth"
	oauthgoogle "brianhang.me/facegraph/internal/oauth/google"
	"brianhang.me/facegraph/internal/user"
)

var isSecure bool

func setupRoutes() error {
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	googleOAuth := &oauthgoogle.Strategy{}
	err := oauth.SetupRoutesForStrategy(
		googleOAuth,
		appurl.ForPath("/auth/google"),
		func(w http.ResponseWriter, r *http.Request, googleID string) error {
			u := user.FindOrCreateFromGoogleID(googleID)
			if err := user.SetCookie(w, &u); err != nil {
				return err
			}
			return nil
		},
		api.InternalErrorResponse,
	)
	if err != nil {
		return fmt.Errorf("failed to set up Google OAuth routes: %v", err)
	}

	http.HandleFunc(
		"/api/user/",
		user.RouteWithUser(func(w http.ResponseWriter, r *http.Request, u *user.User) {
			if u != nil {
				io.WriteString(w, fmt.Sprintf("Your user ID is %d", u.ID))
			} else {
				io.WriteString(w, "you are not logged in")
			}
		}),
	)

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
