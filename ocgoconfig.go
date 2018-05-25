package ocgoconfig

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

var (
	env       func(string) string = os.Getenv
	exporters                     = map[string]func(in *url.URL) (trace.Exporter, error){}

	ErrInvalidConfiguration = errors.New("invalid configuration")
	ErrUnknownExporter      = errors.New("unknown exporter")
)

func init() {
	exporters["jaeger"] = func(in *url.URL) (trace.Exporter, error) {
		q := in.Query()
		serviceName := q.Get("name")

		protocol := "https"
		if q.Get("insecure") == "true" {
			protocol = "http"
		}

		pass, _ := in.User.Password()
		return jaeger.NewExporter(jaeger.Options{
			Endpoint:    fmt.Sprintf("%s://%s", protocol, in.Host),
			Username:    in.User.Username(),
			Password:    pass,
			ServiceName: serviceName,
		})
	}
	exporters["stackdriver"] = func(in *url.URL) (trace.Exporter, error) {
		q := in.Query()
		return stackdriver.NewExporter(stackdriver.Options{
			ProjectID:    in.Hostname(),
			MetricPrefix: q.Get("metrics_prefix"),
		})
	}
}

func Setup(config string) error {
	if config == "" {
		config = env("OCGO_EXPORTER")
		if config == "" {
			return nil
		}
	}

	parts, err := url.Parse(config)
	if err != nil {
		return ErrInvalidConfiguration
	}

	factory, found := exporters[strings.ToLower(parts.Scheme)]
	if !found {
		return ErrUnknownExporter
	}

	exporter, err := factory(parts)
	if err != nil {
		return err
	}

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(getTraceConfig(parts.Query()))

	return nil
}

func Available() []string {
	i := 0
	a := make([]string, len(exporters))
	for n, _ := range exporters {
		a[i] = n
		i++
	}
	return a
}

func getTraceConfig(config url.Values) trace.Config {
	c := trace.Config{}

	// Sampler
	switch config.Get("sampler") {
	case "never":
		// This seems odd yes it may mean don't provide it in the
		// first place, however a HTTP call could change this value
		// if the application wanted that to happen.
		c.DefaultSampler = trace.NeverSample()
	case "5050":
		c.DefaultSampler = trace.ProbabilitySampler(0.5)
	default:
		c.DefaultSampler = trace.AlwaysSample()
	}
	// End Sampler

	return c
}
