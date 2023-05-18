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

	query = `
		create table if not exists opened_positions(
			order_id text primary key,
			profile_id integer not null,
			symbol text not null,
			FOREIGN KEY(profile_id) REFERENCES profiles(id));
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists users(
			id integer primary key autoincrement,
			profile_picture_path text null,
			license_key text not null,
			persistent_login boolean null,
			theme text null);
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	return err
}

func (repo *SQLiteRepository) InsertProfile(p user.Profile) (*user.Profile, error) {
	stmt := "insert into profiles (title, trader_id, bybit_api_key, bybit_api_secret, max_bybit_binance_price_difference_percent, leverage, initial_open_percent, max_add_multiplier, open_delay, one_coin_max_percent, blacklist_coins, add_prevention_percent, block_adds_above_entry, max_open_positions, auto_tp, auto_sl, test_mode) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
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

	res, err := repo.Conn.Exec(stmt, p.Title, p.TraderID, p.BybitApiKey, p.BybitApiSecret, p.MaxBybitBinancePriceDifferentPercent, p.Leverage, p.InitialOpenPercent, p.MaxAddMultiplier, p.OpenDelay, p.OneCoinMaxPercent, blackListCoins, p.AddPreventionPercent, blockAddsAboveEntry, p.MaxOpenPositions, p.AutoTP, p.AutoSL, testMode)
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

func (repo *SQLiteRepository) InsertOpenedPosition(p user.OpenedPosition) (*user.OpenedPosition, error) {
	stmt := "insert into opened_positions (order_id, profile_id, symbol) values (?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, p.OrderID, p.ProfileID, p.Symbol)
	if err != nil {
		return nil, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *SQLiteRepository) InsertUser(u user.User) (*user.User, error) {
	stmt := "insert into users (profile_picture_path, license_key, persistent_login, theme) values (?, ?, ?, ?)"
	var persistent int

	if u.PersistentLogin {
		persistent = 1
	} else {
		persistent = 0
	}

	res, err := repo.Conn.Exec(stmt, u.ProfilePicturePath, u.LicenseKey, persistent, u.Theme)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return &u, nil
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

func (repo *SQLiteRepository) AllOpenedPositions() ([]user.OpenedPosition, error) {
	query := "select order_id, profile_id, symbol from opened_positions order by profile_id"

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.OpenedPosition
	for rows.Next() {
		var p user.OpenedPosition

		err := rows.Scan(
			&p.OrderID,
			&p.ProfileID,
			&p.Symbol,
		)
		if err != nil {
			return nil, err
		}

		all = append(all, p)
	}

	return all, nil
}

func (repo *SQLiteRepository) GetUser(ID int64) (*user.User, error) {
	query := "select id, profile_picture_path, license_key, persistent_login, theme from users where id = ?"

	rows, err := repo.Conn.Query(query, ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var u *user.User
	for rows.Next() {
		var persistent int

		err := rows.Scan(
			&u.ID,
			&u.ProfilePicturePath,
			&u.LicenseKey,
			&persistent,
			&u.Theme,
		)
		if err != nil {
			return nil, err
		}

		if persistent == 0 {
			u.PersistentLogin = false
		} else {
			u.PersistentLogin = true
		}
	}

	return u, nil
}

func (repo *SQLiteRepository) GetAllUsers() (*user.User, error) {
	query := "select id, profile_picture_path, license_key, persistent_login, theme from users"

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var u user.User
	for rows.Next() {
		var persistent int

		err := rows.Scan(
			&u.ID,
			&u.ProfilePicturePath,
			&u.LicenseKey,
			&persistent,
			&u.Theme,
		)
		if err != nil {
			return nil, err
		}

		if persistent == 0 {
			u.PersistentLogin = false
		} else {
			u.PersistentLogin = true
		}
	}

	return &u, nil
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

func (repo *SQLiteRepository) UpdateOpenedPosition(orderId string, updated user.OpenedPosition) error {
	if orderId == "" {
		return errors.New("invalid updated id")
	}

	stmt := "update opened_positions set profile_id = ?, symbol = ? where order_id = ?"
	res, err := repo.Conn.Exec(stmt, updated.ProfileID, updated.Symbol, orderId)
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

func (repo *SQLiteRepository) UpdateUser(id int64, updated user.User) error {
	if id <= 0 {
		return errors.New("invalid updated id")
	}

	var persistent int

	if updated.PersistentLogin {
		persistent = 1
	} else {
		persistent = 0
	}

	stmt := "update users set profile_picture_path = ?, license_key = ?,  persistent_login = ?, theme = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.ProfilePicturePath, updated.LicenseKey, persistent, updated.Theme, id)
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

func (repo *SQLiteRepository) DeleteOpenedPosition(orderId string) error {
	res, err := repo.Conn.Exec("delete from opened_positions where order_id = ?", orderId)
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

func (repo *SQLiteRepository) DeleteUser(id int64) error {
	res, err := repo.Conn.Exec("delete from users where id = ?", id)
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

func (repo *SQLiteRepository) DeleteAllUsers() error {
	res, err := repo.Conn.Exec("delete from users where 1 = 1")
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
