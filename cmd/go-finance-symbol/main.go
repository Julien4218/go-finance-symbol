package main

import (
	"flag"
	"log"

	"github.com/Julien4218/go-finance-symbol/observability"
	s "github.com/Julien4218/go-finance-symbol/symbol"
	"github.com/spf13/cobra"
)

func main() {
	if err := Command.Execute(); err != nil {
		if err != flag.ErrHelp {
			log.Fatal(err)
		}
	}
}

func init() {
}

func globalInit(cmd *cobra.Command, args []string) {
	observability.Init()
}

var Command = &cobra.Command{
	Use:              "go-finance-symbol",
	Short:            "Test App",
	PersistentPreRun: globalInit,
	Long:             `Execute`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			observability.Log("Execute with a symbol as a command argument, for example AAPL for apple")
			return
		}
		for _, symbol := range args {
			intervals := []s.IntervalRange{
				s.FiveDay,
				s.OneMonth,
				s.ThreeMonth,
				s.SixMonth,
				s.OneYear,
				s.YearToDate}
			s.Execute(symbol, intervals)
		}
		observability.Shutdown()
	},
}
