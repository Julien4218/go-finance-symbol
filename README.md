# go-finance-symbol

Get price information for 1 or many finance symbols and optionally send this as dimensional metrics to New Relic.
Example

```bash
 make && ./bin/linux/go-finance-symbol AAPL
```

## Observability

Optionally, Metrics for each result and symbols can be sent to New Relic. To do so use the following environment variable `NEW_RELIC_LICENSE_KEY=<my-license-key>` with the application.
If using an alternate environment than US, specify the metric endpoint by using an environment variable like this `NEW_RELIC_METRIC_API=https://metric-api.eu.newrelic.com/metric/v1`.
