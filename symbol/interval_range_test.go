package symbol

import (
	"testing"
)

func TestIntervalRangeString(t *testing.T) {
	tests := []struct {
		input    IntervalRange
		expected string
	}{
		{OneDay, "1d"},
		{FiveDay, "5d"},
		{OneMonth, "1mo"},
		{ThreeMonth, "3mo"},
		{SixMonth, "6mo"},
		{OneYear, "1y"},
		{TwoYear, "2y"},
		{FiveYear, "5y"},
		{TenYear, "10y"},
		{YearToDate, "ytd"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Errorf("IntervalRange.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
