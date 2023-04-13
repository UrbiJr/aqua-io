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
		create table if not exists profile_groups(
			id integer primary key autoincrement,
			name text not null);
	`
	_, err := repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists profiles(
			id integer primary key autoincrement,
			group_id integer not null,
			title text not null,
			bybit_api_key text null,
			bybit_api_secret text null,
			max_bybit_binance_price_difference_percent real null,
			leverage real null,
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
			test_mode boolean null,
			foreign key (group_id) references profile_groups (id));
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists traders(
			encryptedUid primary key,
			futureUid text null,
			nickName text null,
			userPhotoUrl text null,
			rank integer null,
			pnl real null,
			roi real null,
			positionShared boolean null,
			twitterUrl text null,
			updateTime integer null,
			followerCount integer null,
			isTwTrader boolean null,
			openId text null);
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists positions(
			id integer primary key autoincrement,
			trader_id text not null,
			symbol text null,
			entryPrice real null,
			markPrice real null,
			pnl text real null,
			roe text real null,
			amount real null,
			updateTimeStamp integer null,
			yellow boolean null,
			tradeBefore boolean null,
			leverage integer null,
			foreign key (trader_id) references traders (id));
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	return err
}

func (repo *SQLiteRepository) InsertProfile(p user.Profile) (*user.Profile, error) {
	stmt := "insert into profiles (group_id, title, bybit_api_key, bybit_api_secret, max_bybit_binance_price_difference_percent, leverage, initial_open_percent, max_add_multiplier, open_delay, one_coin_max_percent, blacklist_coins, add_prevention_percent, block_adds_above_entry, max_open_positions, auto_tp, auto_sl, test_mode) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
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

	res, err := repo.Conn.Exec(stmt, p.GroupID, p.Title, p.BybitApiKey, p.BybitApiSecret, p.MaxBybitBinancePriceDifferentPercent, p.Leverage, p.InitialOpenPercent, p.MaxAddMultiplier, p.OpenDelay, p.OneCoinMaxPercent, blackListCoins, p.AddPreventionPercent, blockAddsAboveEntry, p.MaxOpenPositions, p.AutoTP, p.AutoSL, testMode)
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

func (repo *SQLiteRepository) InsertProfileGroup(pg user.ProfileGroup) (*user.ProfileGroup, error) {
	stmt := "insert into profile_groups (name) values (?)"

	res, err := repo.Conn.Exec(stmt, pg.Name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	pg.ID = id
	return &pg, nil
}

func (repo *SQLiteRepository) AllProfiles() ([]user.Profile, error) {
	query := "select id, group_id, title, bybit_api_key, bybit_api_secret, max_bybit_binance_price_difference_percent, leverage, initial_open_percent, max_add_multiplier, open_delay, one_coin_max_percent, blacklist_coins, add_prevention_percent, block_adds_above_entry, max_open_positions, auto_tp, auto_sl, test_mode from profiles order by title"

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
			&p.GroupID,
			&p.Title,
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

func (repo *SQLiteRepository) AllProfileGroups() ([]user.ProfileGroup, error) {
	query := "select id, name from profile_groups order by name"

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.ProfileGroup
	for rows.Next() {
		var p user.ProfileGroup
		err := rows.Scan(
			&p.ID,
			&p.Name,
		)
		if err != nil {
			return nil, err
		}
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

	stmt := "update profiles set group_id = ?, title = ?, bybit_api_key = ?,  bybit_api_secret = ?,  max_bybit_binance_price_difference_percent = ?, leverage = ?, initial_open_percent = ?,  max_add_multiplier = ?,  open_delay = ?,  one_coin_max_percent = ?,  blacklist_coins = ?,  add_prevention_percent = ?,  block_adds_above_entry = ?,  max_open_positions = ?,  auto_tp = ?,  auto_sl = ?,  test_mode = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.GroupID, updated.Title, updated.BybitApiKey, updated.BybitApiSecret, updated.MaxBybitBinancePriceDifferentPercent, updated.Leverage, updated.InitialOpenPercent, updated.MaxAddMultiplier, updated.OpenDelay, updated.OneCoinMaxPercent, blackListCoins, updated.AddPreventionPercent, blockAddsAboveEntry, updated.MaxOpenPositions, updated.AutoTP, updated.AutoSL, testMode, id)
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

func (repo *SQLiteRepository) UpdateProfileGroup(id int64, updated user.ProfileGroup) error {
	if id <= 0 {
		return errors.New("invalid updated id")
	}

	stmt := "update profile_groups set name = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.Name, id)
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

func (repo *SQLiteRepository) DeleteProfilesByGroupID(id int64) error {
	res, err := repo.Conn.Exec("delete from profiles where group_id = ?", id)
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

func (repo *SQLiteRepository) DeleteProfileGroup(id int64) error {
	_, err := repo.Conn.Exec("delete from profiles where group_id = ?", id)
	if err != nil {
		return err
	}

	res, err := repo.Conn.Exec("delete from profile_groups where id = ?", id)
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

func (repo *SQLiteRepository) InsertPosition(p user.Position) (*user.Position, error) {
	stmt := "insert into positions (trader_id,symbol,entryPrice,markPrice,pnl,roe,amount,updateTimeStamp,yellow,tradeBefore,leverage) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, p.TraderID, p.Symbol, p.EntryPrice, p.MarkPrice, p.Pnl, p.Roe, p.Amount, p.UpdateTimestamp, p.Yellow, p.TradeBefore, p.Leverage)
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

func (repo *SQLiteRepository) InsertTrader(t user.Trader) (*user.Trader, error) {
	stmt := "insert into traders (encryptedUid,futureUid,nickName,userPhotoUrl,rank,pnl,roi,positionShared,twitterUrl,updateTime,followerCount,isTwTrader,openId) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, t.EncryptedUid, t.FutureUid, t.NickName, t.UserPhotoUrl, t.Rank, t.Pnl, t.Roi, t.PositionShared, t.TwitterUrl, t.UpdateTime, t.FollowerCount, t.IsTwTrader, t.OpenId)
	if err != nil {
		return nil, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (repo *SQLiteRepository) AllPositions() ([]user.Position, error) {
	query := "select * from positions order by symbol"

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Position
	for rows.Next() {
		var p user.Position

		err := rows.Scan(
			&p.ID,
			&p.TraderID,
			&p.Symbol,
			&p.EntryPrice,
			&p.MarkPrice,
			&p.Pnl,
			&p.Roe,
			&p.Amount,
			&p.UpdateTimestamp,
			&p.Yellow,
			&p.TradeBefore,
			&p.Leverage,
		)
		if err != nil {
			return nil, err
		}

		all = append(all, p)
	}

	return all, nil
}

func (repo *SQLiteRepository) AllTraders() ([]user.Trader, error) {
	query := "select * from traders order by nickname"

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Trader
	for rows.Next() {
		var p user.Trader
		err := rows.Scan(
			&p.EncryptedUid,
			&p.FutureUid,
			&p.NickName,
			&p.UserPhotoUrl,
			&p.Rank,
			&p.Pnl,
			&p.Roi,
			&p.PositionShared,
			&p.TwitterUrl,
			&p.UpdateTime,
			&p.FollowerCount,
			&p.IsTwTrader,
			&p.OpenId,
		)
		if err != nil {
			return nil, err
		}
		all = append(all, p)
	}

	return all, nil
}

func (repo *SQLiteRepository) UpdatePosition(id int64, updated user.Position) error {
	if id <= 0 {
		return errors.New("invalid updated id")
	}

	stmt := "update positions set trader_id = ?, symbol = ?, entryPrice = ?,  markPrice = ?,  pnl = ?, roe = ?, amount = ?,  updateTimeStamp = ?,  yellow = ?,  tradeBefore = ?,  leverage = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.TraderID, updated.Symbol, updated.EntryPrice, updated.MarkPrice, updated.Pnl, updated.Roe, updated.Amount, updated.UpdateTimestamp, updated.Yellow, updated.TradeBefore, updated.Leverage, id)
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

func (repo *SQLiteRepository) UpdateTrader(encryptedUid string, updated user.Trader) error {
	if encryptedUid == "" {
		return errors.New("invalid encryptedUid")
	}

	stmt := "update traders set futureUid = ?, nickName = ?, userPhotoUrl = ?,  rank = ?,  pnl = ?, roi = ?, positionShared = ?,  twitterUrl = ?, updateTime = ?,  followerCount = ?,  isTwTrader = ?,  openId = ? where encryptedUid = ?"
	res, err := repo.Conn.Exec(stmt, updated.FutureUid, updated.NickName, updated.UserPhotoUrl, updated.Rank, updated.Pnl, updated.Roi, updated.PositionShared, updated.TwitterUrl, updated.UpdateTime, updated.FollowerCount, updated.IsTwTrader, updated.OpenId, encryptedUid)
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

func (repo *SQLiteRepository) DeletePosition(id int64) error {
	res, err := repo.Conn.Exec("delete from positions where id = ?", id)
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

func (repo *SQLiteRepository) DeletePositionsByTraderID(encryptedUid string) error {
	res, err := repo.Conn.Exec("delete from positions where trader_id = ?", encryptedUid)
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

func (repo *SQLiteRepository) DeleteTrader(encryptedUid string) error {
	_, err := repo.Conn.Exec("delete from positions where trader_id = ?", encryptedUid)
	if err != nil {
		return err
	}

	res, err := repo.Conn.Exec("delete from traders where encryptedUid = ?", encryptedUid)
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
