package observability

import (
	"context"
	"fmt"
	"os"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	log "github.com/sirupsen/logrus"
)

var (
	// harvester methods are null safe
	harvester *telemetry.Harvester
)

func Init() {
	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if licenseKey == "" {
		Log("environment variable NEW_RELIC_LICENSE_KEY not set, skipping instrumentation")
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
	Log("NewRelic telemetry initialized")
}

func Shutdown() {
	harvester.HarvestNow(context.Background())
}

func Log(message string) {
	log.Infof(message)
	harvester.RecordLog(
		telemetry.Log{
			Message: message,
		},
	)
}

func Logf(format string, args ...interface{}) {
	Log(fmt.Sprintf(format, args...))
}
