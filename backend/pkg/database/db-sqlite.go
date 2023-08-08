package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/UrbiJr/aqua-io/backend/internal/user"
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
			exchange text null,
			account_name text null,
			public_api text null,
			secret_api text null,
			passphrase text null,
			stop_if_fall_under real null,
			test_mode boolean null);
	`
	_, err := repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists strategies(
			id integer primary key autoincrement,
			position_size real null,
			percentage real null);
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists copied_traders(
			id string primary key,
			profile_id integer null,
			encrypted_uid text null,
			trade_mode text null,
			leverage real null,
			max_open_positions integer null,
			max_coin_percentage_position integer null,
			price_difference_between_exchanges real null,
			open_delay_between_positions integer null,
			block_position_adds boolean null,
			auto_take_profit_strategy integer null,
			auto_stop_loss_strategy integer null,
			max_coin_allocation real null,
			max_add_multiplier real null,
			add_prevention_percent real null,
			blacklisted_coins text null,
			stop_control boolean null,
			FOREIGN KEY(profile_id) REFERENCES profiles(id),
			FOREIGN KEY(auto_take_profit_strategy) REFERENCES strategies(id),
			FOREIGN KEY(auto_stop_loss_strategy) REFERENCES strategies(id));
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
			theme text null,
			close_all_trades_when_closing boolean null);
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	return err
}

func (repo *SQLiteRepository) InsertProfile(p user.Profile) (*user.Profile, error) {
	stmt := `
		insert into profiles (
			title,
			exchange,
			account_name,
			public_api,
			secret_api,
			passphrase,
			stop_if_fall_under,
			test_mode
		) values (?, ?, ?, ?, ?, ?, ?, ?)
	`

	var testMode int
	if p.TestMode {
		testMode = 1
	} else {
		testMode = 0
	}

	res, err := repo.Conn.Exec(stmt, p.Title, p.Exchange, p.AccountName, p.PublicAPI, p.SecretAPI, p.Passphrase, p.StopIfFallUnder, testMode)
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

func (repo *SQLiteRepository) InsertStrategy(s user.Strategy) (*user.Strategy, error) {
	stmt := `
		insert into strategies (
			position_size,
			percentage,
		) values (?, ?)
	`

	res, err := repo.Conn.Exec(stmt, s.PositionSize, s.Percentage)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	s.ID = id
	return &s, nil
}

func (repo *SQLiteRepository) InsertCopiedTrader(t user.Trader) (*user.Trader, error) {
	stmt := `
		insert into copied_traders (
		profile_id,
		encrypted_uid,
		trade_mode, 
		leverage, 
		max_open_positions, 
		max_coin_percentage_position, 
		price_difference_between_exchanges, 
		open_delay_between_positions, 
		block_position_adds, 
		auto_take_profit_strategy, 
		auto_stop_loss_strategy, 
		max_coin_allocation, 
		max_add_multiplier, 
		add_prevention_percent, 
		blacklisted_coins, 
		stop_control
	) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	var blackListCoins string
	var blockPositionAdds, stopControl int

	if t.BlockPositionAdds {
		blockPositionAdds = 1
	} else {
		blockPositionAdds = 0
	}

	if t.StopControl {
		stopControl = 1
	} else {
		stopControl = 0
	}

	blackListCoins = strings.Join(t.BlackListedCoins, ",")

	res, err := repo.Conn.Exec(stmt,
		t.ProfileID,
		t.TradeMode,
		t.Leverage,
		t.MaxOpenPositions,
		t.MaxCoinPercentagePosition,
		t.MaxPriceDifferenceBetweenExchange,
		t.OpenDelayBetweenPositions,
		blockPositionAdds,
		t.AutoTakeProfit.ID,
		t.AutoStopLoss.ID,
		t.MaxCoinAllocation,
		t.MaxAddMultiplier,
		t.AddPreventionPercent,
		blackListCoins,
		stopControl,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id
	return &t, nil
}

func (repo *SQLiteRepository) InsertUser(u user.User) (*user.User, error) {
	stmt := "insert into users (profile_picture_path, license_key, persistent_login, theme, close_all_trades_when_closing) values (?, ?, ?, ?, ?)"
	var persistent, closeAllTradesWhenClosing int

	if u.PersistentLogin {
		persistent = 1
	} else {
		persistent = 0
	}

	if u.CloseAllTradesWhenClosing {
		closeAllTradesWhenClosing = 1
	} else {
		closeAllTradesWhenClosing = 0
	}

	res, err := repo.Conn.Exec(stmt, u.ProfilePicturePath, u.LicenseKey, persistent, u.Theme, closeAllTradesWhenClosing)
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
	query := `
		select id,
		title,
		exchange,
		account_name,
		public_api,
		secret_api,
		passphrase,
		stop_if_fall_under,
		test_mode
	 	from profiles order by title
	`

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Profile
	for rows.Next() {
		var p user.Profile
		var testMode int

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Exchange,
			&p.AccountName,
			&p.PublicAPI,
			&p.SecretAPI,
			&p.Passphrase,
			&p.StopIfFallUnder,
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

		all = append(all, p)
	}

	return all, nil
}

func (repo *SQLiteRepository) AllStrategies() ([]user.Strategy, error) {
	query := `
		select id,
		position_size,
		percentage,
	 	from strategies
	`

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Strategy
	for rows.Next() {
		var s user.Strategy

		err := rows.Scan(
			&s.ID,
			&s.PositionSize,
			&s.Percentage,
		)
		if err != nil {
			return nil, err
		}

		all = append(all, s)
	}

	return all, nil
}

func (repo *SQLiteRepository) AllCopiedTraders() ([]user.Trader, error) {
	query := `
		select id,
		profile_id,
		encrypted_uid,
		trade_mode, 
		leverage, 
		max_open_positions, 
		max_coin_percentage_position, 
		price_difference_between_exchanges, 
		open_delay_between_positions, 
		block_position_adds, 
		auto_take_profit_strategy, 
		auto_stop_loss_strategy, 
		max_coin_allocation, 
		max_add_multiplier, 
		add_prevention_percent, 
		blacklisted_coins, 
		stop_control
	 	from copied_traders order by profile_id
	 `

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Trader
	for rows.Next() {
		var t user.Trader
		var blackListCoins string
		var blockPositionAdds, stopControl int

		err := rows.Scan(
			&t.ID,
			&t.ProfileID,
			&t.TradeMode,
			&t.Leverage,
			&t.MaxOpenPositions,
			&t.MaxCoinPercentagePosition,
			&t.MaxPriceDifferenceBetweenExchange,
			&t.OpenDelayBetweenPositions,
			&blockPositionAdds,
			&t.AutoTakeProfit.ID,
			&t.AutoStopLoss.ID,
			&t.MaxCoinAllocation,
			&t.MaxAddMultiplier,
			&t.AddPreventionPercent,
			&blackListCoins,
			&stopControl,
		)
		if err != nil {
			return nil, err
		}

		if stopControl == 0 {
			t.StopControl = false
		} else {
			t.StopControl = true
		}
		if blockPositionAdds == 0 {
			t.BlockPositionAdds = false
		} else {
			t.BlockPositionAdds = true
		}

		t.BlackListedCoins = strings.Split(blackListCoins, ",")

		s, err := repo.GetStrategy(t.AutoTakeProfit.ID)
		if err != nil {
			t.AutoTakeProfit = user.Strategy{}
		} else {
			t.AutoTakeProfit = *s
		}

		s, err = repo.GetStrategy(t.AutoStopLoss.ID)
		if err != nil {
			t.AutoStopLoss = user.Strategy{}
		} else {
			t.AutoStopLoss = *s
		}

		all = append(all, t)
	}

	return all, nil
}

func (repo *SQLiteRepository) AllUsers() ([]user.User, error) {
	query := "select id, profile_picture_path, license_key, persistent_login, theme, close_all_trades_when_closing from users"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.User
	for rows.Next() {
		var u user.User
		var persistent, closeAllTradesWhenClosing int

		err := rows.Scan(
			&u.ID,
			&u.ProfilePicturePath,
			&u.LicenseKey,
			&persistent,
			&u.Theme,
			&closeAllTradesWhenClosing,
		)
		if err != nil {
			return nil, err
		}

		if persistent == 0 {
			u.PersistentLogin = false
		} else {
			u.PersistentLogin = true
		}

		if closeAllTradesWhenClosing == 0 {
			u.CloseAllTradesWhenClosing = false
		} else {
			u.CloseAllTradesWhenClosing = true
		}

		all = append(all, u)
	}

	return all, nil
}

func (repo *SQLiteRepository) GetUser(ID int64) (*user.User, error) {
	query := "select id, profile_picture_path, license_key, persistent_login, theme, close_all_trades_when_closing from users where id = ?"

	rows, err := repo.Conn.Query(query, ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var u *user.User
	for rows.Next() {
		var persistent, closeAllTradesWhenClosing int

		err := rows.Scan(
			&u.ID,
			&u.ProfilePicturePath,
			&u.LicenseKey,
			&persistent,
			&u.Theme,
			&closeAllTradesWhenClosing,
		)
		if err != nil {
			return nil, err
		}

		if persistent == 0 {
			u.PersistentLogin = false
		} else {
			u.PersistentLogin = true
		}

		if closeAllTradesWhenClosing == 0 {
			u.CloseAllTradesWhenClosing = false
		} else {
			u.CloseAllTradesWhenClosing = true
		}
	}

	return u, nil
}

func (repo *SQLiteRepository) GetStrategy(ID int64) (*user.Strategy, error) {
	query := "select id, position_size, percentage from strategies where id = ?"

	rows, err := repo.Conn.Query(query, ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var s *user.Strategy
	for rows.Next() {
		err := rows.Scan(
			&s.ID,
			&s.PositionSize,
			&s.Percentage,
		)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (repo *SQLiteRepository) UpdateProfile(ID int64, updated user.Profile) error {
	if ID <= 0 {
		return errors.New("invalid updated id")
	}

	var testMode int

	if updated.TestMode {
		testMode = 1
	} else {
		testMode = 0
	}

	stmt := `
		update profiles set
			title = ?,
			exchange = ?,
			account_name = ?,
			public_api = ?,
			secret_api = ?,
			passphrase = ?,
			stop_if_fall_under = ?,
			test_mode = ?
		where id = ?
	`
	res, err := repo.Conn.Exec(stmt, updated.Title, updated.Exchange, updated.AccountName, updated.PublicAPI, updated.SecretAPI, updated.Passphrase, updated.StopIfFallUnder, testMode, ID)
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

func (repo *SQLiteRepository) UpdateStrategy(ID int64, updated user.Strategy) error {
	if ID <= 0 {
		return errors.New("invalid updated id")
	}

	stmt := `
		update strategies set
			position_size = ?,
			percentage = ?,
		where id = ?
	`

	res, err := repo.Conn.Exec(stmt, updated.PositionSize, updated.Percentage, ID)
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

func (repo *SQLiteRepository) UpdateCopiedTrader(ID int64, updated user.Trader) error {
	if ID <= 0 {
		return errors.New("invalid updated id")
	}

	stmt := `
		update copied_traders set
		profile_id = ?, 
		encrypted_uid = ?,
		trade_mode = ?, 
		leverage = ?, 
		max_open_positions = ?, 
		max_coin_percentage_position = ?, 
		price_difference_between_exchanges = ?, 
		open_delay_between_positions = ?, 
		block_position_adds = ?, 
		auto_take_profit_strategy = ?, 
		auto_stop_loss_strategy = ?, 
		max_coin_allocation = ?, 
		max_add_multiplier = ?, 
		add_prevention_percent = ?, 
		blacklisted_coins = ?, 
		stop_control = ?
	where id  = ?
	`
	var blackListCoins string
	var blockPositionAdds, stopControl int

	if updated.BlockPositionAdds {
		blockPositionAdds = 1
	} else {
		blockPositionAdds = 0
	}

	if updated.StopControl {
		stopControl = 1
	} else {
		stopControl = 0
	}

	blackListCoins = strings.Join(updated.BlackListedCoins, ",")

	res, err := repo.Conn.Exec(stmt,
		updated.ProfileID,
		updated.EncryptedUid,
		updated.TradeMode,
		updated.Leverage,
		updated.MaxOpenPositions,
		updated.MaxCoinPercentagePosition,
		updated.MaxPriceDifferenceBetweenExchange,
		updated.OpenDelayBetweenPositions,
		blockPositionAdds,
		updated.AutoTakeProfit.ID,
		updated.AutoStopLoss.ID,
		updated.MaxCoinAllocation,
		updated.MaxAddMultiplier,
		updated.AddPreventionPercent,
		blackListCoins,
		stopControl,
		ID,
	)
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

func (repo *SQLiteRepository) UpdateUser(ID int64, updated user.User) error {
	if ID <= 0 {
		return errors.New("invalid updated id")
	}

	var persistent, closeAllTradesWhenClosing int

	if updated.PersistentLogin {
		persistent = 1
	} else {
		persistent = 0
	}

	if updated.CloseAllTradesWhenClosing {
		closeAllTradesWhenClosing = 1
	} else {
		closeAllTradesWhenClosing = 0
	}

	stmt := "update users set profile_picture_path = ?, license_key = ?,  persistent_login = ?, theme = ?, close_all_trades_when_closing = ?, where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.ProfilePicturePath, updated.LicenseKey, persistent, updated.Theme, closeAllTradesWhenClosing, ID)
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

func (repo *SQLiteRepository) DeleteProfile(ID int64) error {
	res, err := repo.Conn.Exec("delete from profiles where id = ?", ID)
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

func (repo *SQLiteRepository) DeleteStrategy(ID int64) error {
	res, err := repo.Conn.Exec("delete from strategies where id = ?", ID)
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

func (repo *SQLiteRepository) DeleteCopiedTrader(ID int64) error {
	res, err := repo.Conn.Exec("delete from copied_traders where id = ?", ID)
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

func (repo *SQLiteRepository) DeleteUser(ID int64) error {
	res, err := repo.Conn.Exec("delete from users where id = ?", ID)
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

func (repo *SQLiteRepository) DeleteAllStrategies() error {
	res, err := repo.Conn.Exec("delete from strategies where 1 = 1")
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

func (repo *SQLiteRepository) DeleteAllCopiedTraders() error {
	res, err := repo.Conn.Exec("delete from copied_traders where 1 = 1")
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
