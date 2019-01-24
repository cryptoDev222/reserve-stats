package storage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	logschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	walletschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/walletfee"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

const (
	//timePrecision is the precision configured for influxDB
	timePrecision           = "s"
	tradeLogMeasurementName = "trades"
	burnFeesMeasurementName = "burn_fees"
	walletMeasurementName   = "wallet_fees"
)

// InfluxStorage represent a client to store trade data to influx DB
type InfluxStorage struct {
	sugar                *zap.SugaredLogger
	dbName               string
	influxClient         client.Client
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface
	kycChecker           kycChecker

	// traded stored traded addresses to use in a single SaveTradeLogs
	traded map[ethereum.Address]struct{}
}

// NewInfluxStorage init an instance of InfluxStorage
func NewInfluxStorage(sugar *zap.SugaredLogger, dbName string, influxClient client.Client,
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface, kycChecker kycChecker) (*InfluxStorage, error) {
	storage := &InfluxStorage{
		sugar:                sugar,
		dbName:               dbName,
		influxClient:         influxClient,
		tokenAmountFormatter: tokenAmountFormatter,
		kycChecker:           kycChecker,
		traded:               make(map[ethereum.Address]struct{}),
	}
	if err := storage.createDB(); err != nil {
		return nil, err
	}
	return storage, nil
}

// SaveTradeLogs persist trade logs to DB
func (is *InfluxStorage) SaveTradeLogs(logs []common.TradeLog) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  is.dbName,
		Precision: timePrecision,
	})
	if err != nil {
		return err
	}
	for _, log := range logs {
		points, err := is.tradeLogToPoint(log)
		if err != nil {
			return err
		}

		for _, pt := range points {
			bp.AddPoint(pt)
		}
	}

	if err := is.influxClient.Write(bp); err != nil {
		return err
	}

	if len(logs) > 0 {
		is.sugar.Debugw("saved trade logs into influxdb",
			"first_block", logs[0].BlockNumber,
			"last_block", logs[len(logs)-1].BlockNumber,
			"trade_logs", len(logs))
	} else {
		is.sugar.Debugw("no trade log to store")
	}

	// reset traded map to avoid ever growing size
	is.traded = make(map[ethereum.Address]struct{})
	return nil
}

// LastBlock returns last stored trade log block number from database.
func (is InfluxStorage) LastBlock() (int64, error) {
	q := fmt.Sprintf(`SELECT "block_number","eth_amount" from "trades" ORDER BY time DESC limit 1`)

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return 0, err
	}

	if len(res) != 1 || len(res[0].Series) != 1 || len(res[0].Series[0].Values[0]) != 3 {
		is.sugar.Info("no result returned for last block query")
		return 0, nil
	}

	return influxdb.GetInt64FromTagValue(res[0].Series[0].Values[0][1])
}

func prepareTradeLogQuery() string {
	var (
		tradeLogQueryFields = []logschema.FieldName{
			logschema.Time,
			logschema.BlockNumber,
			logschema.EthAmount,
			logschema.UserAddr,
			logschema.SrcAddr,
			logschema.DstAddr,
			logschema.SrcAmount,
			logschema.DstAmount,
			logschema.IP,
			logschema.Country,
			logschema.IntegrationApp,
			logschema.DestBurnFee,
			logschema.SourceBurnFee,
			logschema.LogIndex,
			logschema.TxHash,
		}
		tradeLogQuery string
	)
	for _, field := range tradeLogQueryFields {
		tradeLogQuery += field.String() + ", "
	}
	fiatAmount := fmt.Sprintf("(%s * %s) AS %s", logschema.EthAmount.String(), logschema.EthUSDRate.String(), logschema.FiatAmount.String())
	tradeLogQuery += fiatAmount
	return tradeLogQuery
}

func prepareWalletFeeQuery() string {
	var (
		walletFeeFields = []walletschema.FieldName{
			walletschema.Time,
			walletschema.ReserveAddr,
			walletschema.WalletAddr,
			walletschema.Amount,
			walletschema.LogIndex,
		}
		walletQuery string
	)
	for i, field := range walletFeeFields {
		if i != 0 {
			walletQuery += ", "
		}
		walletQuery += field.String()
	}
	return walletQuery
}

// LoadTradeLogs return trade logs from DB
func (is *InfluxStorage) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	var (
		result = make([]common.TradeLog, 0)
		q      = fmt.Sprintf(
			`
		SELECT %[1]s FROM %[5]s WHERE time >= '%[3]s' AND time <= '%[4]s' GROUP BY %[7]s;
		SELECT %[2]s FROM %[6]s WHERE time >= '%[3]s' AND time <= '%[4]s';
		`,
			prepareWalletFeeQuery(),
			prepareTradeLogQuery(),
			from.Format(time.RFC3339),
			to.Format(time.RFC3339),
			walletMeasurementName,
			tradeLogMeasurementName,
			walletschema.TxHash.String()+", "+walletschema.TradeLogIndex.String(),
		)

		logger = is.sugar.With(
			"func", "tradelogs/storage/InfluxStorage.LoadTradLogs",
			"from", from,
			"to", to,
		)
	)
	logger.Debug("prepared query statement", "query", q)

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return nil, err
	}

	// Get WalletFees
	// map [tx_hash][trade_log_index][]common.WalletFee
	walletFeesByTxHash := make(map[ethereum.Hash]map[uint][]common.WalletFee)

	if len(res[0].Series) == 0 {
		is.sugar.Debug("empty wallet fee in query result")
	} else {
		for _, row := range res[0].Series {
			txHash, tradeLogIndex, walletFees, err := is.rowToWalletFees(row)
			if err != nil {
				return nil, err
			}
			_, exist := walletFeesByTxHash[txHash]
			if !exist {
				walletFeesByTxHash[txHash] = make(map[uint][]common.WalletFee)
			}
			walletFeesByTxHash[txHash][uint(tradeLogIndex)] = walletFees
		}
	}

	// Get TradeLogs
	if len(res[1].Series) == 0 {
		is.sugar.Debug("empty trades in query result")
		return result, nil
	}
	idxs, err := logschema.NewFieldsRegistrar(res[1].Series[0].Columns)
	if err != nil {
		return nil, err
	}
	for _, row := range res[1].Series[0].Values {

		tradeLog, err := is.rowToTradeLog(row, walletFeesByTxHash, idxs)
		if err != nil {
			return nil, err
		}
		result = append(result, tradeLog)
	}

	return result, nil
}

// createDB creates the database will be used for storing trade logs measurements.
func (is *InfluxStorage) createDB() error {
	_, err := is.queryDB(is.influxClient, fmt.Sprintf("CREATE DATABASE %s", is.dbName))
	return err
}

// queryDB convenience function to query the database
func (is *InfluxStorage) queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: is.dbName,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (is *InfluxStorage) setBurnFeeTagsAndFields(log common.TradeLog, tags map[string]string, fields map[string]interface{}) error {
	var logger = is.sugar.With(
		"func", "tradelogs/storage/setBurnFeeTagsAndFields",
		"log", log,
	)
	if blockchain.IsBurnable(log.SrcAddress) {
		if len(log.BurnFees) < 1 {
			logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "at least 1 burn fees (src)")
			return nil
		}
		tags[logschema.SrcReserveAddr.String()] = log.BurnFees[0].ReserveAddress.String()
		burnAmount, err := is.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[0].Amount)
		if err != nil {
			return err
		}
		fields[logschema.SourceBurnFee.String()] = burnAmount

		if blockchain.IsBurnable(log.DestAddress) {
			if len(log.BurnFees) < 2 {
				logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "2 burn fees (src-dst)")
				return nil
			}
			tags[logschema.DstReserveAddr.String()] = log.BurnFees[1].ReserveAddress.String()
			burnAmount, err = is.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[1].Amount)
			if err != nil {
				return err
			}
			fields[logschema.DestBurnFee.String()] = burnAmount
			return nil
		}

		return nil
	}

	if blockchain.IsBurnable(log.DestAddress) {
		if len(log.BurnFees) < 1 {
			logger.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "at least 1 burn fees (dst)")
			return nil
		}
		tags[logschema.DstReserveAddr.String()] = log.BurnFees[0].ReserveAddress.String()
		burnAmount, err := is.tokenAmountFormatter.FromWei(blockchain.KNCAddr, log.BurnFees[0].Amount)
		if err != nil {
			return err
		}
		fields[logschema.DestBurnFee.String()] = burnAmount
	}
	return nil
}

func (is *InfluxStorage) tradeLogToPoint(log common.TradeLog) ([]*client.Point, error) {
	var points []*client.Point
	var walletAddr ethereum.Address
	if len(log.WalletFees) > 0 {
		walletAddr = log.WalletFees[0].WalletAddress
	}

	tags := map[string]string{

		logschema.UserAddr.String(): log.UserAddress.String(),

		logschema.SrcAddr.String():        log.SrcAddress.String(),
		logschema.DstAddr.String():        log.DestAddress.String(),
		logschema.IntegrationApp.String(): log.IntegrationApp,
		logschema.WalletAddress.String():  walletAddr.String(),
		logschema.LogIndex.String():       strconv.FormatUint(uint64(log.Index), 10),

		logschema.Country.String(): log.Country,

		logschema.LogIndex.String():       strconv.FormatUint(uint64(log.Index), 10),
		logschema.SrcReserveAddr.String(): log.SrcReserveAddress.String(),
		logschema.DstReserveAddr.String(): log.DstReserveAddress.String(),
	}

	ethReceivalAmount, err := is.tokenAmountFormatter.FromWei(blockchain.ETHAddr, log.EtherReceivalAmount)
	if err != nil {
		return nil, err
	}

	srcAmount, err := is.tokenAmountFormatter.FromWei(log.SrcAddress, log.SrcAmount)
	if err != nil {
		return nil, err
	}

	dstAmount, err := is.tokenAmountFormatter.FromWei(log.DestAddress, log.DestAmount)
	if err != nil {
		return nil, err
	}

	var ethAmount float64

	if log.SrcAddress == blockchain.ETHAddr {
		ethAmount = srcAmount
	} else if log.DestAddress == blockchain.ETHAddr {
		ethAmount = dstAmount
	} else {
		ethAmount = ethReceivalAmount
	}

	fields := map[string]interface{}{

		logschema.SrcAmount.String():  srcAmount,
		logschema.DstAmount.String():  dstAmount,
		logschema.EthUSDRate.String(): log.ETHUSDRate,

		logschema.EthAmount.String():      ethAmount,
		logschema.BlockNumber.String():    strconv.FormatUint(log.BlockNumber, 10),
		logschema.TxHash.String():         log.TransactionHash.String(),
		logschema.IP.String():             log.IP,
		logschema.EthUSDProvider.String(): log.ETHUSDProvider,
	}

	if err = is.setBurnFeeTagsAndFields(log, tags, fields); err != nil {
		return nil, err
	}

	tradePoint, err := client.NewPoint(tradeLogMeasurementName, tags, fields, log.Timestamp)
	if err != nil {
		return nil, err
	}

	points = append(points, tradePoint)
	// build walletFeePoint
	for _, walletFee := range log.WalletFees {
		tags := map[string]string{
			walletschema.ReserveAddr.String():   walletFee.ReserveAddress.String(),
			walletschema.WalletAddr.String():    walletFee.WalletAddress.String(),
			walletschema.Country.String():       log.Country,
			walletschema.TradeLogIndex.String(): strconv.FormatUint(uint64(log.Index), 10),
			walletschema.TxHash.String():        log.TransactionHash.String(),
			walletschema.LogIndex.String():      strconv.FormatUint(uint64(walletFee.Index), 10),
		}

		amount, err := is.tokenAmountFormatter.FromWei(blockchain.KNCAddr, walletFee.Amount)
		if err != nil {
			return nil, err
		}

		fields := map[string]interface{}{
			walletschema.Amount.String(): amount,
		}

		walletFeePoint, err := client.NewPoint(walletMeasurementName, tags, fields, log.Timestamp)
		if err != nil {
			return nil, err
		}

		points = append(points, walletFeePoint)
	}

	firstTradePoint, err := is.assembleFirstTradePoint(log)
	if err != nil {
		return nil, err
	}
	if firstTradePoint != nil {
		points = append(points, firstTradePoint)
	}

	kycedPoint, err := is.assembleKYCPoint(log)
	if err != nil {
		return nil, err
	}

	if kycedPoint != nil {
		points = append(points, kycedPoint)
	}

	return points, nil
}

func (is *InfluxStorage) assembleFirstTradePoint(logItem common.TradeLog) (*client.Point, error) {
	var logger = is.sugar.With(
		"func", "tradelogs/storage/InfluxStorage.assembleFirstTradePoint",
		"timestamp", logItem.Timestamp.String(),
		"user_addr", logItem.UserAddress.Hex(),
		"country", logItem.Country,
	)

	if _, ok := is.traded[logItem.UserAddress]; ok {
		logger.Debug("user has already traded, ignoring")
		return nil, nil
	}

	traded, err := is.userTraded(logItem.UserAddress)
	if err != nil {
		return nil, err
	}

	if traded {
		return nil, nil
	}

	logger.Debugw("user first trade")
	tags := map[string]string{
		"user_addr": logItem.UserAddress.Hex(),
		"country":   logItem.Country,
	}

	for _, walletFee := range logItem.WalletFees {
		tags["wallet_addr"] = walletFee.WalletAddress.Hex()
	}

	fields := map[string]interface{}{
		"traded": true,
	}

	point, err := client.NewPoint("first_trades", tags, fields, logItem.Timestamp)
	if err != nil {
		return nil, err
	}

	is.traded[logItem.UserAddress] = struct{}{}
	return point, nil
}

func (is *InfluxStorage) userTraded(addr ethereum.Address) (bool, error) {
	q := fmt.Sprintf("SELECT traded FROM first_trades WHERE user_addr='%s'", addr.String())
	response, err := is.queryDB(is.influxClient, q)
	if err != nil {
		return false, err
	}
	// if there is no record, this mean the address has not traded yet
	if (len(response) == 0) || (len(response[0].Series) == 0) || (len(response[0].Series[0].Values) == 0) {
		return false, nil
	}
	return true, nil
}

func (is *InfluxStorage) assembleKYCPoint(logItem common.TradeLog) (*client.Point, error) {
	var logger = is.sugar.With(
		"func", "tradelogs/storage/InfluxStorage.assembleKYCPoint",
		"timestamp", logItem.Timestamp.String(),
		"user_addr", logItem.UserAddress.Hex(),
		"country", logItem.Country,
	)

	kyced, err := is.kycChecker.IsKYCed(logItem.UserAddress, logItem.Timestamp)
	if err != nil {
		return nil, err
	}

	if !kyced {
		logger.Debugw("user has not been kyced yet")
		return nil, nil
	}

	logger.Debugw("user has been kyced")
	tags := map[string]string{
		"user_addr": logItem.UserAddress.Hex(),
		"country":   logItem.Country,
	}

	for _, walletFee := range logItem.WalletFees {
		tags["wallet_addr"] = walletFee.WalletAddress.Hex()
	}

	fields := map[string]interface{}{
		"kyced": true,
	}

	point, err := client.NewPoint("kyced", tags, fields, logItem.Timestamp)
	return point, err
}
