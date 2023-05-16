package rest

import (
	"net/http"
)

var (
	M controller = &Micro{}
)

type controller interface {
	Metrics(w http.ResponseWriter, req *http.Request)
	Health(w http.ResponseWriter, req *http.Request)
}

// micro usecase
type Micro struct{}
