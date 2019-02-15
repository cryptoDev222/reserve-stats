package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/users/common"
)

const (
	addressesTableName = "addresses"
	usersTableName     = "users"
)

//UserDB is storage of user data
type UserDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//DeleteAllTables delete all table from schema using for test only
func (udb *UserDB) DeleteAllTables() error {
	_, err := udb.db.Exec(fmt.Sprintf(`DROP TABLE "%s", "%s"`, addressesTableName, usersTableName))
	return err
}

//NewDB open a new database connection
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*UserDB, error) {
	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%s"
(
  id           SERIAL PRIMARY KEY,
  email        text      NOT NULL UNIQUE,
  last_updated TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "%s"
(
  id        SERIAL PRIMARY KEY,
  address   text      NOT NULL UNIQUE,
  timestamp TIMESTAMP NOT NULL,
  user_id   SERIAL    NOT NULL REFERENCES users (id)
);
`
	var logger = sugar.With("func", "users/storage.NewDB")

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	logger.Debug("initializing database schema")
	if _, err = tx.Exec(fmt.Sprintf(schemaFmt, usersTableName, addressesTableName)); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &UserDB{
		sugar: sugar,
		db:    db,
	}, nil
}

//Close close db connection and return error if any
func (udb *UserDB) Close() error {
	return udb.db.Close()
}

//CreateOrUpdate store user info to persist in database
func (udb *UserDB) CreateOrUpdate(userData common.UserData) error {
	var (
		logger = udb.sugar.With(
			"func", "users/storage.CreateOrUpdate",
			"email", userData.Email,
		)
		stmt       string
		addresses  []string
		timestamps []int64
	)

	for _, ui := range userData.UserInfo {
		addresses = append(addresses, ui.Address)
	}
	for _, ui := range userData.UserInfo {
		timestamps = append(timestamps, ui.Timestamp)
	}

	tx, err := udb.db.Beginx()
	if err != nil {
		return err
	}

	defer pgsql.CommitOrRollback(tx, logger, &err)

	stmt = fmt.Sprintf(`WITH u AS (
  INSERT INTO "%s" (email, last_updated)
    VALUES ($1, NOW())
    ON CONFLICT ON CONSTRAINT users_email_key
      DO UPDATE SET last_updated = NOW() RETURNING id
),
     a AS (
       SELECT unnest($2::text[])             AS address,
              unnest($3::double precision[]) AS timestamp
     )
INSERT
INTO "%s"(address, timestamp, user_id)
SELECT a.address, to_timestamp(a.timestamp / 1000), u.id
FROM u NATURAL JOIN a
ON CONFLICT ON CONSTRAINT addresses_address_key DO UPDATE SET timestamp = EXCLUDED.timestamp, user_id = EXCLUDED.user_id
`,
		usersTableName,
		addressesTableName)
	logger.Debugw("upsert email and Ethereum addresses",
		"stmt", stmt)
	_, err = tx.Exec(stmt,
		userData.Email,
		pq.StringArray(addresses),
		pq.Int64Array(timestamps),
	)
	if err != nil {
		return err
	}

	logger.Debugw("delete removed Ethereum addresses",
		"stmt", stmt)
	stmt = fmt.Sprintf(`DELETE
FROM "%s"
WHERE user_id IN (SELECT id AS user_id FROM "%s" WHERE email = $1)
 AND address NOT IN (SELECT unnest($2::text[]) as address)
`,
		addressesTableName,
		usersTableName)
	_, err = tx.Exec(stmt,
		userData.Email,
		pq.StringArray(addresses))
	if err != nil {
		return err
	}
	return err
}

//GetAllAddresses return all user address info from addresses table
func (udb *UserDB) GetAllAddresses() ([]string, error) {
	var result []string
	if err := udb.db.Select(&result, fmt.Sprintf(`SELECT address FROM "%s"`, addressesTableName)); err != nil {
		return result, err
	}
	return result, nil
}
