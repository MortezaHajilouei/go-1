package middleware

import (
	"context"
	"go.uber.org/fx"
	"micro/adapter/redis"
	"micro/config"
	"reflect"

	"google.golang.org/grpc"
)

type Params struct {
	fx.In
}

func NewMiddleware(params Params) Middleware {
	return &middle{}
}

// Middleware interface
type Middleware interface {
	JWT(ctx context.Context) (context.Context, error)
	SJWT(ctx context.Context, req ...interface{}) (context.Context, error)
	assignMiddleware(ctx context.Context, req interface{}, middlewares []string) error
	MiddlewareUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// middle struct
type middle struct {
	redis redis.AuthStore
}

func (m *middle) MiddlewareUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// loop for all routes we have in config file
	for _, r := range config.C().Service.Router {
		// if method name != proto rpc name, then go to next method
		if r.Method != info.FullMethod {
			continue
		}
		if err := m.assignMiddleware(ctx, req, r.Middlewares); err != nil {
			return nil, err
		}
	}

	h, err := handler(ctx, req)

	return h, err
}

func (m *middle) assignMiddleware(ctx context.Context, req interface{}, middlewares []string) error {

	// loop middlewares for every route
	for _, m := range middlewares {

		// get middleware methods by name
		method := reflect.ValueOf(&middle{}).MethodByName(m)
		if !method.IsValid() {
			continue
		}

		// check every middleware for method
		responses := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
		if err := responses[1].Interface(); err != nil {
			return err.(error)
		}
	}

	return nil
}
