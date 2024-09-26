package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrCreateCounter(t *testing.T) {
	name := "test_counter"
	counter := GetOrCreateCounter(name)
	assert.NotNil(t, counter)
	assert.Equal(t, "finance_symbol_test_counter", counter.Name())
}

func TestGetOrCreateGauge(t *testing.T) {
	name := "test_gauge"
	gauge := GetOrCreateGauge(name)
	assert.NotNil(t, gauge)
	assert.Equal(t, "finance_symbol_test_gauge", gauge.Name())
}
