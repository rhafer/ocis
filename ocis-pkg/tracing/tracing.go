package tracing

import (
	"fmt"
	"net/url"
	"strings"

	rtrace "github.com/cs3org/reva/v2/pkg/trace"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Propagator ensures the importer module uses the same trace propagation strategy.
var Propagator = propagation.NewCompositeTextMapPropagator(
	propagation.Baggage{},
	propagation.TraceContext{},
)

// GetTraceProvider returns a configured open-telemetry trace provider.
func GetTraceProvider(agentEndpoint, collectorEndpoint, serviceName, traceType string) (*sdktrace.TracerProvider, error) {
	switch t := traceType; t {
	case "", "jaeger":
		var (
			exp *jaeger.Exporter
			err error
		)

		if agentEndpoint != "" {
			var agentHost string
			var agentPort string

			agentHost, agentPort, err = parseAgentConfig(agentEndpoint)
			if err != nil {
				return nil, err
			}

			exp, err = jaeger.New(
				jaeger.WithAgentEndpoint(
					jaeger.WithAgentHost(agentHost),
					jaeger.WithAgentPort(agentPort),
				),
			)
		} else if collectorEndpoint != "" {
			exp, err = jaeger.New(
				jaeger.WithCollectorEndpoint(
					jaeger.WithEndpoint(collectorEndpoint),
				),
			)
		}
		if err != nil {
			return nil, err
		}

		rtrace.InitDefaultTracerProvider(collectorEndpoint, agentEndpoint)
		return sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName)),
			),
		), nil

	case "agent":
		fallthrough
	case "zipkin":
		fallthrough
	default:
		return nil, fmt.Errorf("unknown trace type %s", traceType)
	}
}

func parseAgentConfig(ae string) (string, string, error) {
	u, err := url.Parse(ae)
	// as per url.go:
	// [...] Trying to parse a hostname and path
	// without a scheme is invalid but may not necessarily return an
	// error, due to parsing ambiguities.
	if err == nil && u.Hostname() != "" && u.Port() != "" {
		return u.Hostname(), u.Port(), nil
	}

	p := strings.Split(ae, ":")
	if len(p) != 2 {
		return "", "", fmt.Errorf(fmt.Sprintf("invalid agent endpoint `%s`. expected format: `hostname:port`", ae))
	}

	switch {
	case p[0] == "" && p[1] == "": // case ae = ":"
		return "", "", fmt.Errorf(fmt.Sprintf("invalid agent endpoint `%s`. expected format: `hostname:port`", ae))
	case p[0] == "":
		return "", "", fmt.Errorf(fmt.Sprintf("invalid agent endpoint `%s`. expected format: `hostname:port`", ae))
	}
	return p[0], p[1], nil
}

// Configure for Reva serves only as informational / instructive log messages. Tracing config will be delegated directly
// to Reva services.
func Configure(enabled bool, tracingType string, logger log.Logger) {
	if enabled {
		switch tracingType {
		case "agent":
			logger.Error().
				Str("type", tracingType).
				Msg("Reva only supports the jaeger tracing backend")

		case "jaeger":
			logger.Info().
				Str("type", tracingType).
				Msg("configuring storage to use the jaeger tracing backend")

		case "zipkin":
			logger.Error().
				Str("type", tracingType).
				Msg("Reva only supports the jaeger tracing backend")

		default:
			logger.Warn().
				Str("type", tracingType).
				Msg("Unknown tracing backend")
		}

	} else {
		logger.Debug().
			Msg("Tracing is not enabled")
	}
}
