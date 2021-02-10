package pkg

import "net/http"

func handleResponseError(statusCode int, errorString string, w http.ResponseWriter){
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(errorString))
}
