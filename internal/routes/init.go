package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Init() error {
	port, err := getPort()
	if err != nil {
		return err
	}
	addr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/", homeRoute)

	certFile := os.Getenv("WEBSERVER_TLS_CERT_FILE")
	keyFile := os.Getenv("WEBSERVER_TLS_KEY_FILE")

	if certFile != "" && keyFile != "" {
		log.Printf("Listening (with TLS) on port %d\n", port)
		return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	}

	log.Printf("Listening on port %d\n", port)
	return http.ListenAndServe(addr, nil)
}

func getPort() (int, error) {
	rawPort := os.Getenv("WEBSERVER_PORT")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return port, fmt.Errorf("\"%s\" is not a valid port number", rawPort)
	}
	return port, nil
}
