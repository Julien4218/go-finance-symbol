package observability

import (
	"context"
	"os"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	log "github.com/sirupsen/logrus"
)

var (
	harvester *telemetry.Harvester
)

func Init() {
	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if licenseKey == "" {
		log.Warn("environment variable NEW_RELIC_LICENSE_KEY not set, skipping instrumentation")
		return
	}

	key := telemetry.ConfigAPIKey(licenseKey)
	var err error
	var metricUrl func(*telemetry.Config)
	metricApi := os.Getenv("NEW_RELIC_METRIC_API")
	if len(metricApi) > 0 {
		metricUrl = telemetry.ConfigMetricsURLOverride(metricApi)
		harvester, err = telemetry.NewHarvester(key, metricUrl)
	} else {
		harvester, err = telemetry.NewHarvester(key)
	}
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("NewRelic telemetry initialized")
}

func Shutdown() {
	harvester.HarvestNow(context.Background())
}
