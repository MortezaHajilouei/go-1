package server

import (
	"net/http"

	"google.golang.org/grpc"
)

type IServer interface {
	ListenAndServe() error
	ServerOptions() []grpc.ServerOption
	Shutdown() error
	DialOptions() []grpc.DialOption
	HandleFuncs(*http.ServeMux)
}
