package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	fetcher "github.com/KyberNetwork/reserve-stats/accounting/binance-fetcher"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

const (
	fromIDFlag        = "from-id"
	retryDelayFlag    = "retry-delay"
	attemptFlag       = "attempt"
	batchSizeFlag     = "batch-size"
	defaultRetryDelay = 2 // minute
	defaultAttempt    = 4
	defaultBatchSize  = 100
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting binance trades fetcher"
	app.Usage = "Fetch and store trades history from binance"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.IntFlag{
			Name:   retryDelayFlag,
			Usage:  "delay time when do a retry",
			EnvVar: "RETRY_DELAY",
			Value:  defaultRetryDelay,
		},
		cli.IntFlag{
			Name:   attemptFlag,
			Usage:  "number of time doing retry",
			EnvVar: "ATTEMPT",
			Value:  defaultAttempt,
		},
		cli.IntFlag{
			Name:   batchSizeFlag,
			Usage:  "batch to request to binance",
			EnvVar: "BATCH_SIZE",
			Value:  defaultBatchSize,
		},
		cli.Uint64Flag{
			Name:   fromIDFlag,
			Usage:  "id to get trade history from",
			EnvVar: "FROM_ID",
		},
	)

	app.Flags = append(app.Flags, binance.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}

	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("initiate fetcher")

	binanceClient, err := binance.NewClientFromContext(c, sugar)
	if err != nil {
		return err
	}

	fromID := c.Uint64(fromIDFlag)

	retryDelay := c.Int(retryDelayFlag)
	attempt := c.Int(attemptFlag)
	batchSize := c.Int(batchSizeFlag)
	binanceFetcher := fetcher.NewFetcher(sugar, binanceClient, retryDelay, attempt, batchSize)

	tradeHistories, err := binanceFetcher.GetTradeHistory(fromID)
	if err != nil {
		return err
	}

	sugar.Debugw("trade histories", "result", tradeHistories)

	return nil
}