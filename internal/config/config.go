package config

import (
	"runtime"
	"strings"
	"time"

	"github.com/bfi-finance/bfi-go-pkg/appinfo"
	"github.com/bfi-finance/bfi-go-pkg/grpc/grpcserver"
	"github.com/bfi-finance/bfi-go-pkg/logger"
	"github.com/bfi-finance/bfi-go-pkg/openapi"
	"github.com/bfi-finance/bfi-go-pkg/pgsql"
	"github.com/bfi-finance/bfi-go-pkg/tracer"
	"github.com/caarlos0/env/v6"
)

// flatEnv is struct to parse environment variable values.
// Each struct field representing environment value key to be parsed.
// Reference: https://github.com/caarlos0/env
//
//nolint:govet
type flatEnv struct {
	AppName string `env:"APP_NAME"`
	AppEnv  string `env:"APP_ENV"` // local|dev|uat|sit|prod

	LoggerLevel  string `env:"LOGGER_LEVEL"`
	LoggerOutput string `env:"LOGGER_OUTPUT"`

	TracingEnable           bool    `env:"TRACING_ENABLE"`
	TracingExporter         string  `env:"TRACING_EXPORTER"`
	TracingExporterEndpoint string  `env:"TRACING_EXPORTER_ENDPOINT"`
	TracingSampling         string  `env:"TRACING_SAMPLING"`
	TracingSamplingRatio    float64 `env:"TRACING_SAMPLING_RATIO"`

	GRPCPort int `env:"GRPC_PORT"`

	GRPCRestProxyPort              int    `env:"GRPC_REST_PROXY_PORT"`
	GRPCRestProxyGrpcHost          string `env:"GRPC_REST_PROXY_GRPC_HOST"`
	GRPCRestProxyReadHeaderTimeout int    `env:"GRPC_REST_PROXY_READ_HEADER_TIMEOUT"`
	GRPCRestProxyShutdownTimeout   int    `env:"GRPC_REST_PROXY_SHUTDOWN_TIMEOUT"`
	GRPCRestProxyEnable            bool   `env:"GRPC_REST_PROXY_ENABLE"`

	OpenAPIUIEnable         bool   `env:"OPENAPI_UI_ENABLE"`
	OpenAPIUITemplate       string `env:"OPENAPI_UI_TEMPLATE"`
	OpenAPIUIPath           string `env:"OPENAPI_UI_PATH"`
	OpenAPIUIServerBasePath string `env:"OPENAPI_UI_SERVER_BASE_PATH"`

	PostgresHost        string `env:"POSTGRES_HOST"`
	PostgresPort        int    `env:"POSTGRES_PORT"`
	PostgresDatabase    string `env:"POSTGRES_DATABASE"`
	PostgresUsername    string `env:"POSTGRES_USERNAME"`
	PostgresPassword    string `env:"POSTGRES_PASSWORD,unset"`
	PostgresLogging     bool   `env:"POSTGRES_LOGGING"`
	PostgresConnMaxOpen int    `env:"POSTGRES_CONN_MAX_OPEN"`
	PostgresConnMaxIdle int    `env:"POSTGRES_CONN_MAX_IDLE"`
	PostgresSSLMode     string `env:"POSTGRES_SSL_MODE"`
}

// Config is the main application config that is set from env var parse result.
//
//nolint:govet
type Config struct {
	AppInfo             appinfo.Info
	Tracer              tracer.Config
	Logger              logger.Config
	OpenAPIUI           openapi.UIConfig
	GRPCServer          grpcserver.GRPCConfig
	GRPCRESTProxyServer grpcserver.RESTProxyConfig
	//
	PostgreSQL pgsql.Config
}

// LoadFromEnv load environment variables into a private flat struct then parse to actual internal config.
// Environment variables should be set right before running the app and we don't have to specify .env file manually.
// See docker-compose.yml in base folder for reference.
func LoadFromEnv() (*Config, error) {
	var envCfg flatEnv
	err := env.Parse(&envCfg)
	if err != nil {
		return nil, err
	}
	// ========== BASE CONFIG ==========
	cfg := Config{
		AppInfo: appinfo.Info{
			Name:          envCfg.AppName,
			Env:           envCfg.AppEnv,
			GitURL:        strings.TrimSpace(appinfo.GitURL),        // set back from build set
			GitCommitHash: strings.TrimSpace(appinfo.GitCommitHash), // set back from build set
			GitTag:        strings.TrimSpace(appinfo.GitTag),        // set back from build set
			BuildOS:       strings.TrimSpace(appinfo.BuildOS),       // set back from build set
			BuildTime:     strings.TrimSpace(appinfo.BuildTime),     // set back from build set
			GoVersion:     runtime.Version(),
		},
		Logger: logger.Config{
			Level:  envCfg.LoggerLevel,
			Output: envCfg.LoggerOutput,
		},
		Tracer: tracer.Config{
			Enable:           envCfg.TracingEnable,
			Exporter:         envCfg.TracingExporter,
			ExporterEndpoint: envCfg.TracingExporterEndpoint,
			Sampling:         envCfg.TracingSampling,
			SamplingRatio:    envCfg.TracingSamplingRatio,
		},
		GRPCServer: grpcserver.GRPCConfig{
			Port: envCfg.GRPCPort,
		},
		GRPCRESTProxyServer: grpcserver.RESTProxyConfig{
			GRPCServerHost:    envCfg.GRPCRestProxyGrpcHost,
			Port:              envCfg.GRPCRestProxyPort,
			ReadHeaderTimeout: time.Second * time.Duration(envCfg.GRPCRestProxyReadHeaderTimeout),
			ShutdownTimeout:   time.Second * time.Duration(envCfg.GRPCRestProxyShutdownTimeout),
			Enable:            envCfg.GRPCRestProxyEnable,
		},
		OpenAPIUI: openapi.UIConfig{
			Enable:         envCfg.OpenAPIUIEnable,
			Template:       envCfg.OpenAPIUITemplate,
			Path:           envCfg.OpenAPIUIPath,
			ServerBasePath: envCfg.OpenAPIUIServerBasePath,
		},
	}
	// ========== IF REQUIRE POSTGRESQL
	cfg.PostgreSQL = pgsql.Config{
		Username:    envCfg.PostgresUsername,
		Password:    envCfg.PostgresPassword,
		Database:    envCfg.PostgresDatabase,
		Host:        envCfg.PostgresHost,
		Port:        envCfg.PostgresPort,
		ConnMaxOpen: envCfg.PostgresConnMaxOpen,
		ConnMaxIdle: envCfg.PostgresConnMaxIdle,
		Logging:     envCfg.PostgresLogging,
		LogLevel:    envCfg.LoggerLevel, // use global logger log level
		SSLMode:     envCfg.PostgresSSLMode,
		Tracing:     envCfg.TracingEnable, // use global tracing flag
	}
	return &cfg, nil
}
