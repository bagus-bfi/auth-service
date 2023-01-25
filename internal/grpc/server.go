package grpc

import (
	"database/sql"
	"fmt"

	"bravo-go-template/internal/config"
	"bravo-go-template/internal/grpc/hc"

	"github.com/bfi-finance/bfi-go-pkg/eventbus"
	"github.com/bfi-finance/bfi-go-pkg/grpc/grpcinterceptor/datadogloggerinterceptor"
	"github.com/bfi-finance/bfi-go-pkg/grpc/grpcinterceptor/loggerinterceptor"
	"github.com/bfi-finance/bfi-go-pkg/grpc/grpcinterceptor/panicinterceptor"
	"github.com/bfi-finance/bfi-go-pkg/grpc/grpcinterceptor/requestloggerinterceptor"
	"github.com/bfi-finance/bfi-go-pkg/grpc/grpcserver"
	"github.com/bfi-finance/bfi-go-pkg/openapi"
	"github.com/bfi-finance/bfi-protobuf/gen/go/bfi/bravoservice/healthcheck"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func New(
	cfg *config.Config,
	zlogger zerolog.Logger,
	evBus eventbus.EventBus,
	sqlDB *sql.DB,
) (*grpcserver.Server, error) {
	// -----------------------------------------------------------------------------------------------------------------
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),                 // open telemetry tracing
		loggerinterceptor.UnaryServer(zlogger),            // request scoped context logger set
		datadogloggerinterceptor.UnaryServer(cfg.AppInfo), // to connect logs and trace
		requestloggerinterceptor.UnaryServer(),            // grpc request logger
		panicinterceptor.UnaryServer(),                    // catch and log panic
		// more interceptors here
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),                 // open telemetry tracing
		loggerinterceptor.StreamServer(zlogger),            // request scoped context logger set
		datadogloggerinterceptor.StreamServer(cfg.AppInfo), // to connect logs and trace
		requestloggerinterceptor.StreamServer(),            // grpc request logger
		panicinterceptor.StreamServer(),                    // catch and log panic
		// more interceptors here
	}
	// -----------------------------------------------------------------------------------------------------------------
	// Setup gRPC server and its provided services.
	grpcServer, err := grpcserver.NewGRPC(
		cfg.GRPCServer,
		grpcserver.WithGRPCServerOption(grpc.ChainUnaryInterceptor(unaryInterceptors...)),
		grpcserver.WithGRPCServerOption(grpc.ChainStreamInterceptor(streamInterceptors...)),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc: failed to construct gRPC server: %w", err)
	}
	// -----------------------------------------------------------------------------------------------------------------
	restProxies := make([]grpcserver.RestProxyEndpointRegisterFunc, 0)
	openapiUIFiles := make([]openapi.EmbeddedFile, 0)
	// -----------------------------------------------------------------------------------------------------------------
	// ===== TODO: CONSTRUCT AND REGISTER GRPC SERVICE HERE =====
	_ = sqlDB
	_ = evBus
	// -----------------------------------------------------------------------------------------------------------------
	// put healthcheck service at the end
	hcGRPCHandler := hc.New()
	healthcheck.RegisterHealthcheckServiceServer(grpcServer, hcGRPCHandler)                      // REGISTER GRPC
	restProxies = append(restProxies, healthcheck.RegisterHealthcheckServiceHandlerFromEndpoint) // REGISTER REST PROXY
	openapiUIFiles = append(openapiUIFiles, hcGRPCHandler.OpenAPIFile())                         // REGISTER API DOC
	// -----------------------------------------------------------------------------------------------------------------
	server := grpcserver.New(
		grpcServer,
		zlogger,
		cfg.AppInfo,
		cfg.GRPCServer,
		cfg.GRPCRESTProxyServer,
		cfg.OpenAPIUI,
		// ===== OPTIONAL VARIADIC PARAMS =====
		grpcserver.WithServerRESTProxyOpenAPIUIFiles(openapiUIFiles...),
		grpcserver.WithServerRESTProxyEndpoints(restProxies...),
		// to set custom rest proxy chi mux: grpcserver.WithServerRESTProxyChiMux()
		// to add custom middleware to default rest proxy mux: grpcserver.WithServerRESTProxyMiddlewares()
	)
	return server, nil
}
