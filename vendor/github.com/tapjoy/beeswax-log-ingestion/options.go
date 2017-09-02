package main

import (
	newrelic "github.com/newrelic/go-agent"
	"github.com/tapjoy/beeswax-log-ingestion/configs"
	"github.com/tapjoy/beeswax-log-ingestion/ingestion"
)

// TODO: fix configs.Data dependency to allow for testing
func getServerOptions() ([]ingestion.Option, error) {
	opts := []ingestion.Option{
		ingestion.ReadTimeout(1),
		ingestion.WriteTimeout(1),
		ingestion.IdleTimeout(1),
		ingestion.ClientTimeout(1),
		ingestion.RateLimit(10),
	}

	cfg := configs.Data
	if cfg.NewRelicEnabled {
		app, err := createNewRelicApp(cfg.NewRelicAppName, cfg.NewRelicID)
		if err != nil {
			return []ingestion.Option{}, err
		}
		opts = append(opts, ingestion.NewRelic(app))
	}

	return opts, nil
}

func createNewRelicApp(appName, license string) (*newrelic.Application, error) {
	cfg := newrelic.NewConfig(appName, license)
	cfg.Enabled = true

	app, err := newrelic.NewApplication(cfg)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
