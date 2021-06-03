package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

type Prometheus struct {
	Url       string
	UserName  string
	Password  config.Secret
	CountDays int64
	Timeout   time.Duration
}

func (prom Prometheus) InitPrometheus() v1.API {
	cfg := api.Config{
		Address: prom.Url,
	}

	if prom.UserName != "" && prom.Password != "" {
		cfg.RoundTripper = config.NewBasicAuthRoundTripper(prom.UserName, prom.Password, "", api.DefaultRoundTripper)
	}

	client, err := api.NewClient(cfg)

	if err != nil {
		logrus.WithError(err).Fatalf("Could create prometheus-client!")
	}

	v1api := v1.NewAPI(client)
	return v1api
}

func (prom Prometheus) QueryRange(v1api v1.API, query string, r v1.Range) (model.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), prom.Timeout)
	defer cancel()
	result, warnings, err := v1api.QueryRange(ctx, query, r)

	if len(warnings) > 0 {
		logrus.Warnf("Prometheus-Warnings %v", warnings)
	}

	return result, err
}
