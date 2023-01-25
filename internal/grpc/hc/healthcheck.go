package hc

import (
	"context"

	"github.com/bfi-finance/bfi-go-pkg/openapi"
	apidoc "github.com/bfi-finance/bfi-protobuf/gen/docs/bfi/bravoservice/healthcheck"
	"github.com/bfi-finance/bfi-protobuf/gen/go/bfi/bravoservice/healthcheck"
)

// Healthcheck implement common healthcheck gRPC service.
type Healthcheck struct {
	healthcheck.UnimplementedHealthcheckServiceServer
}

func New() *Healthcheck { return &Healthcheck{} }

// OpenAPIFile return embedded file for openapi UI.
func (s *Healthcheck) OpenAPIFile() openapi.EmbeddedFile {
	return openapi.EmbeddedFile{FS: apidoc.APIDocFileFS(), FSFilepath: "openapi.yaml", APIName: "Healthcheck"}
}

func (s *Healthcheck) GetHealthcheck(
	_ context.Context,
	_ *healthcheck.HealthcheckRequest,
) (*healthcheck.HealthcheckResponse, error) {
	// if we need to check infrastructure dependencies for readiness, we should check it here concurrently.
	return &healthcheck.HealthcheckResponse{Status: "OK"}, nil
}
