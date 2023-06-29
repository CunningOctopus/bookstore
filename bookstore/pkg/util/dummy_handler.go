package util

import "net/http"

func Dummy(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}
