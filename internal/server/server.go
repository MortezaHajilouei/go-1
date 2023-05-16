package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"micro/adapter/trace"
	"micro/api/pb/base"
	"micro/config"
	"micro/internal/middleware"
	"micro/internal/modules/sampleBase/delivery/grpc"
	"micro/internal/modules/sampleBase/delivery/rest"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Micro usecase
type server struct {
	gwServer   *http.Server
	grpcServer *grpc.Server
	middleware middleware.Middleware
}

type ControllerContainer struct {
	BaseController controller.BaseController
	Middleware     middleware.Middleware
	fx.In
}

func New(cc ControllerContainer) (IServer, error) {
	srv := server{
		middleware: cc.Middleware,
	}
	srv.setupGrpcServer(cc)
	srv.SetupGatewayServer()
	return &srv, nil
}

func (s *server) setupGrpcServer(cc ControllerContainer) {
	s.grpcServer = grpc.NewServer(s.ServerOptions()...)
	base.RegisterSampleAPIServer(s.grpcServer, &cc.BaseController)
	reflection.Register(s.grpcServer)
}

func (s *server) SetupGatewayServer() error {
	conn, err := grpc.DialContext(context.Background(), config.C().Service.GRPC.Host+config.C().Service.GRPC.Port, s.DialOptions()...)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	s.gwServer, err = s.GatewayServer(context.Background(), conn)
	return err
}

func (s *server) ListenAndServe() error {
	lis, err := net.Listen(config.C().Service.GRPC.Protocol, config.C().Service.GRPC.Port)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			zap.L().Error("error serving the grpc server")
		}
	}()

	go func() {
		if err := s.gwServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			zap.L().Error("error serving gateway server")
		}
	}()
	log.Printf("******** router connected ports[GRPC:%s,HTTP:%s] ********** \n",
		config.C().Service.GRPC.Port,
		config.C().Service.HTTP.Port)
	return nil
}

func (s *server) Shutdown() error {
	s.grpcServer.GracefulStop()
	return s.gwServer.Shutdown(context.TODO())
}

// defaultGRPCOptions
// add options for grpc connection
// In order to enable tracing of both upstream and downstream requests of the gRPC usecase, the gRPC adapter must also be initialized with adapter-side opentracing interceptor
// The parent spans created by the gRPC middleware are injected to the go context
func (x *server) ServerOptions() []grpc.ServerOption {
	options := []grpc.ServerOption{}

	// UnaryInterceptor and OpenTracingServerInterceptor for tracer
	options = append(options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		otgrpc.OpenTracingServerInterceptor(trace.T().GetTracer(), otgrpc.LogPayloads()),
		grpc_auth.UnaryServerInterceptor(x.middleware.JWT), // middleware for all routes example
		grpc_prometheus.UnaryServerInterceptor,
		x.middleware.MiddlewareUnaryInterceptor, // middleware for specific methods
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(panicRecover)),
	)))

	options = append(options, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_auth.StreamServerInterceptor(x.middleware.JWT), // middleware for all routes example
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(panicRecover)),
		grpc_prometheus.StreamServerInterceptor,
	)))

	return options
}

// dialOptions, options for dial connections
func (*server) DialOptions() []grpc.DialOption {
	options := []grpc.DialOption{}
	options = append(options, grpc.WithInsecure())
	return options
}

// HandleFuncs method for handler your basci methods
func (*server) HandleFuncs(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", rest.M.Metrics)
	mux.Handle("/health", http.HandlerFunc(rest.M.Health))
}

// with this method we can recover panics in usecase
func panicRecover(ctx context.Context, p interface{}) error {
	zap.L().Error(fmt.Sprintf("panic error: %v %v", ctx, p))
	return fmt.Errorf("PANICED: %+v", p)
}

func (s *server) GatewayServer(ctx context.Context, conn *grpc.ClientConn) (*http.Server, error) {
	// new server from http package
	gwmux := runtime.NewServeMux(runtime.WithErrorHandler(func(c context.Context, sm *runtime.ServeMux, m runtime.Marshaler, rw http.ResponseWriter, r *http.Request, e error) {
		status := status.Convert(e)
		if uint32(status.Code())-(config.C().Service.ID*1000) == http.StatusUnauthorized {
			rw.WriteHeader(http.StatusUnauthorized)
		}
		runtime.DefaultHTTPErrorHandler(c, sm, m, rw, r, e)
	}))

	// register handler
	if err := base.RegisterSampleAPIHandler(ctx, gwmux, conn); err != nil {
		return nil, err
	}

	// handle methods
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	s.HandleFuncs(mux)

	gwServer := &http.Server{
		Addr:    config.C().Service.HTTP.Port,
		Handler: corsHandler(mux),
	}

	return gwServer, nil
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.C().CORS)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, TE, User-Agent, Cache-Control, Sec-Fetch-Dest, Sec-Fetch-Mode, Sec-Fetch-Site, Referer, Content-Type, Pragma, Connection, Content-Length, Accept-Language, Accept-Encoding, Authorization, ResponseType")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
