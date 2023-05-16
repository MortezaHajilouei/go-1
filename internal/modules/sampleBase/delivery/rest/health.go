package rest

import (
	"net/http"
)

// check usecase health
func (a *Micro) Health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}
