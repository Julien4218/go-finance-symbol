package symbol

type IntervalRange int

const (
	OneDay IntervalRange = iota
	FiveDay
	OneMonth
	ThreeMonth
	SixMonth
	OneYear
	TwoYear
	FiveYear
	TenYear
	YearToDate
)

func (d IntervalRange) String() string {
	return [...]string{"1d", "5d", "1mo", "3mo", "6mo", "1y", "2y", "5y", "10y", "ytd"}[d]
}
