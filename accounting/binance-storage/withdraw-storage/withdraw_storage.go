package withdrawstorage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

//BinanceStorage is storage for binance fetcher including trade history and withdraw history
type BinanceStorage struct {
	sugar     *zap.SugaredLogger
	db        *sqlx.DB
	tableName string
}

//NewDB return a new instance of binance storage
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, tableName string) (*BinanceStorage, error) {
	var (
		logger = sugar.With("func", "accounting/binance-storage/binancestorage.NewDB")
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%[1]s"
	(
	  id   text NOT NULL,
	  data JSONB,
	  CONSTRAINT %[1]s_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS %[1]s_time_idx ON %[1]s ((data ->> 'applyTime'));
	`

	query := fmt.Sprintf(schemaFmt, tableName)
	logger.Debugw("create table query", "query", query)

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	logger.Info("binance table init successfully")

	return &BinanceStorage{
		sugar:     sugar,
		db:        db,
		tableName: tableName,
	}, nil
}

//Close database connection
func (bd *BinanceStorage) Close() error {
	if bd.db != nil {
		return bd.db.Close()
	}
	return nil
}

//DeleteTable remove trades table
func (bd *BinanceStorage) DeleteTable() error {
	query := fmt.Sprintf("DROP TABLE %s", bd.tableName)
	if _, err := bd.db.Exec(query); err != nil {
		return err
	}
	return nil
}

//UpdateWithdrawHistory save withdraw history to db
func (bd *BinanceStorage) UpdateWithdrawHistory(withdrawHistories []binance.WithdrawHistory) (err error) {
	var (
		logger       = bd.sugar.With("func", "accounting/binance_storage.UpdateWithdrawHistory")
		withdrawJSON []byte
	)
	const updateQuery = `INSERT INTO %[1]s (id, data)
	VALUES(
		$1,
		$2
	) ON CONFLICT ON CONSTRAINT %[1]s_pk DO NOTHING;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	query := fmt.Sprintf(updateQuery, bd.tableName)
	logger.Debugw("query update withdraw history", "query", query)

	for _, withdraw := range withdrawHistories {
		withdrawJSON, err = json.Marshal(withdraw)
		if err != nil {
			return
		}
		if _, err = tx.Exec(query, withdraw.ID, withdrawJSON); err != nil {
			return
		}
	}

	return
}

//GetWithdrawHistory return list of withdraw fromTime to toTime
func (bd *BinanceStorage) GetWithdrawHistory(fromTime, toTime time.Time) ([]binance.WithdrawHistory, error) {
	var (
		logger   = bd.sugar.With("func", "account/binance_storage.GetTradeHistory")
		result   []binance.WithdrawHistory
		dbResult [][]byte
		tmp      binance.WithdrawHistory
	)
	const selectStmt = `SELECT data FROM %s WHERE data->>'applyTime'>=$1 AND data->>'applyTime'<=$2`
	query := fmt.Sprintf(selectStmt, bd.tableName)

	logger.Debugw("querying trade history...", "query", query)

	from := timeutil.TimeToTimestampMs(fromTime)
	to := timeutil.TimeToTimestampMs(toTime)
	if err := bd.db.Select(&dbResult, query, from, to); err != nil {
		return result, err
	}

	for _, data := range dbResult {
		if err := json.Unmarshal(data, &tmp); err != nil {
			return result, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

//GetLastStoredTimestamp return last timestamp stored in database
func (bd *BinanceStorage) GetLastStoredTimestamp() (time.Time, error) {
	var (
		logger   = bd.sugar.With("func", "account/binance_storage.GetLastStoredTimestamp")
		result   = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		dbResult uint64
	)
	const selectStmt = `SELECT data->>'applyTime' FROM %s ORDER BY data->>'applyTime' DESC LIMIT 1`
	query := fmt.Sprintf(selectStmt, bd.tableName)

	logger.Debugw("querying last stored timestamp...", "query", query)

	if err := bd.db.Get(&dbResult, query); err != nil {
		return result, err
	}

	result = timeutil.TimestampMsToTime(dbResult)

	return result, nil
}
