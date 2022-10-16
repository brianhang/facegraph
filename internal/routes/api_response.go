package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
)

type errorResponse struct {
	Message string `json:"error_message"`
}

func internalErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	dump, _ := httputil.DumpRequest(r, true)
	log.Printf("Request has failed: %v\n\nRequest Dump:\n%s", err, dump)

	errRes := errorResponse{Message: "The server has an error :("}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errRes)
}
