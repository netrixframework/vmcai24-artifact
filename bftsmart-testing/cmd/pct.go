package cmd

import (
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/netrixframework/bftsmart-testing/client"
	"github.com/netrixframework/bftsmart-testing/util"
	"github.com/netrixframework/netrix/config"
	"github.com/netrixframework/netrix/strategies"
	"github.com/netrixframework/netrix/strategies/pct"
	"github.com/spf13/cobra"
)

func PCTStrategy() *cobra.Command {
	return &cobra.Command{
		Use:  "pct [test]",
		Args: cobra.ExactArgs(1),
		Long: "Run PCT without any guidance and check if the property is satisfied for a specific test",
		RunE: func(cmd *cobra.Command, args []string) error {
			termCh := make(chan os.Signal, 1)
			signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)

			_, property := GetTest(args[0])
			if property == nil {
				return errors.New("invalid test")
			}

			var strategy strategies.Strategy = pct.NewPCTStrategy(&pct.PCTStrategyConfig{
				RandSrc:        rand.NewSource(time.Now().UnixMilli()),
				MaxEvents:      1000,
				Depth:          10,
				RecordFilePath: "results",
			})

			bftSmartClient := client.NewBFTSmartClient(&client.BFTSmartClientConfig{
				CodePath: "/netrixframework/bft-smart",
			})

			strategy = strategies.NewStrategyWithProperty(strategy, property)

			driver := strategies.NewStrategyDriver(
				&config.Config{
					APIServerAddr: "127.0.0.1:7074",
					NumReplicas:   4,
					LogConfig: config.LogConfig{
						Format: "json",
						Level:  "info",
						Path:   "results/checker.log",
					},
				},
				&util.BFTSmartParser{},
				strategy,
				&strategies.StrategyConfig{
					Iterations:       iterations,
					IterationTimeout: 40 * time.Second,
					SetupFunc: func(ctx *strategies.Context) {
						go bftSmartClient.Set("name", "jd")
					},
				},
			)

			go func() {
				<-termCh
				driver.Stop()
			}()

			if err := driver.Start(); err != nil {
				panic(err)
			}
			return nil
		},
	}
}
