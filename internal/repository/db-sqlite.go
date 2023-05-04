package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/UrbiJr/aqua-io/internal/user"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

func (repo *SQLiteRepository) Migrate() error {
	query := `
		create table if not exists profiles(
			id integer primary key autoincrement,
			title text not null,
			trader_id text null,
			bybit_api_key text null,
			bybit_api_secret text null,
			max_bybit_binance_price_difference_percent real null,
			leverage integer null,
			initial_open_percent real null,
			max_add_multiplier real null,
			open_delay real null,
			one_coin_max_percent real null,
			blacklist_coins text null,
			add_prevention_percent real null,
			block_adds_above_entry boolean null,
			max_open_positions integer null,
			auto_tp real null,
			auto_sl real null,
			test_mode boolean null);
	`
	_, err := repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	return err
}

func (repo *SQLiteRepository) InsertProfile(p user.Profile) (*user.Profile, error) {
	stmt := "insert into profiles (title, trader_id, bybit_api_key, bybit_api_secret, max_bybit_binance_price_difference_percent, leverage, initial_open_percent, max_add_multiplier, open_delay, one_coin_max_percent, blacklist_coins, add_prevention_percent, block_adds_above_entry, max_open_positions, auto_tp, auto_sl, test_mode) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var blackListCoins string
	var blockAddsAboveEntry, testMode int

	if p.BlockAddsAboveEntry {
		blockAddsAboveEntry = 1
	} else {
		blockAddsAboveEntry = 0
	}

	if p.TestMode {
		testMode = 1
	} else {
		testMode = 0
	}

	blackListCoins = strings.Join(p.BlacklistCoins, ",")

	res, err := repo.Conn.Exec(stmt, p.Title, p.BybitApiKey, p.BybitApiSecret, p.MaxBybitBinancePriceDifferentPercent, p.Leverage, p.InitialOpenPercent, p.MaxAddMultiplier, p.OpenDelay, p.OneCoinMaxPercent, blackListCoins, p.AddPreventionPercent, blockAddsAboveEntry, p.MaxOpenPositions, p.AutoTP, p.AutoSL, testMode)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.ID = id
	return &p, nil
}

func (repo *SQLiteRepository) AllProfiles() ([]user.Profile, error) {
	query := "select id, title, trader_id, bybit_api_key, bybit_api_secret, max_bybit_binance_price_difference_percent, leverage, initial_open_percent, max_add_multiplier, open_delay, one_coin_max_percent, blacklist_coins, add_prevention_percent, block_adds_above_entry, max_open_positions, auto_tp, auto_sl, test_mode from profiles order by title"

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Profile
	for rows.Next() {
		var p user.Profile
		var blackListCoins string
		var blockAddsAboveEntry, testMode int

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.TraderID,
			&p.BybitApiKey,
			&p.BybitApiSecret,
			&p.MaxBybitBinancePriceDifferentPercent,
			&p.Leverage,
			&p.InitialOpenPercent,
			&p.MaxAddMultiplier,
			&p.OpenDelay,
			&p.OneCoinMaxPercent,
			&blackListCoins,
			&p.AddPreventionPercent,
			&blockAddsAboveEntry,
			&p.MaxOpenPositions,
			&p.AutoTP,
			&p.AutoSL,
			&testMode,
		)
		if err != nil {
			return nil, err
		}

		if testMode == 0 {
			p.TestMode = false
		} else {
			p.TestMode = true
		}
		if blockAddsAboveEntry == 0 {
			p.BlockAddsAboveEntry = false
		} else {
			p.BlockAddsAboveEntry = true
		}
		p.BlacklistCoins = strings.Split(blackListCoins, ",")

		all = append(all, p)
	}

	return all, nil
}

func (repo *SQLiteRepository) UpdateProfile(id int64, updated user.Profile) error {
	if id <= 0 {
		return errors.New("invalid updated id")
	}

	var blackListCoins string
	var blockAddsAboveEntry, testMode int

	if updated.BlockAddsAboveEntry {
		blockAddsAboveEntry = 1
	} else {
		blockAddsAboveEntry = 0
	}
	if updated.TestMode {
		testMode = 1
	} else {
		testMode = 0
	}
	blackListCoins = strings.Join(updated.BlacklistCoins, ",")

	stmt := "update profiles set title = ?, trader_id = ?, bybit_api_key = ?,  bybit_api_secret = ?,  max_bybit_binance_price_difference_percent = ?, leverage = ?, initial_open_percent = ?,  max_add_multiplier = ?,  open_delay = ?,  one_coin_max_percent = ?,  blacklist_coins = ?,  add_prevention_percent = ?,  block_adds_above_entry = ?,  max_open_positions = ?,  auto_tp = ?,  auto_sl = ?,  test_mode = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.Title, updated.TraderID, updated.BybitApiKey, updated.BybitApiSecret, updated.MaxBybitBinancePriceDifferentPercent, updated.Leverage, updated.InitialOpenPercent, updated.MaxAddMultiplier, updated.OpenDelay, updated.OneCoinMaxPercent, blackListCoins, updated.AddPreventionPercent, blockAddsAboveEntry, updated.MaxOpenPositions, updated.AutoTP, updated.AutoSL, testMode, id)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows <= 0 {
		return errUpdateFailed
	}

	return nil
}

func (repo *SQLiteRepository) DeleteProfile(id int64) error {
	res, err := repo.Conn.Exec("delete from profiles where id = ?", id)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows <= 0 {
		return errDeleteFailed
	}

	return nil
}

func (repo *SQLiteRepository) DeleteAllProfiles() error {
	res, err := repo.Conn.Exec("delete from profiles where 1 = 1")
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows <= 0 {
		return errDeleteFailed
	}

	return nil
}
