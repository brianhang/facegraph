package routes

import (
	"io"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><h1>Hello!</h1></html>")
}
