package repository

import (
	"database/sql"
	"errors"

	"github.com/UrbiJr/nyx/internal/user"
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
			bybit_api_key text not null,
			bybit_api_secret text not null,
			max_bybit_binance_price_difference_percent float not null,
			initial_open_percent float not null,
			max_add_multiplier float not null,
			open_delay float not null,
			one_coin_max_percent float not null,
			blacklist_coins text not null,
			add_prevention_percent float not null,
			block_adds_above_entry bool not null,
			max_open_positions integer not null,
			auto_tp float not null,
			auto_sl float not null,
			test_mode bool not null,
			foreign key (group_id) references profile_groups (id));
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	return err
}

func (repo *SQLiteRepository) InsertProfile(p user.Profile) (*user.Profile, error) {
	stmt := "insert into profiles (group_id, title, email, first_name, last_name, address_line_1, address_line_2, city, postcode, state, country_code, phone, card_number, card_month, card_year, card_cvv) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, p.GroupID, p.Title, p.Email, p.FirstName, p.LastName, p.AddressLine1, p.AddressLine2, p.City, p.Postcode, p.State, p.CountryCode, p.Phone, p.CardNumber, p.CardMonth, p.CardYear, p.CardCvv)
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
	query := "select id, group_id, title, email, first_name, last_name, address_line_1, address_line_2, city, postcode, state, country_code, phone, card_number, card_month, card_year, card_cvv from profiles order by title"

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Profile
	for rows.Next() {
		var p user.Profile
		err := rows.Scan(
			&p.ID,
			&p.GroupID,
			&p.Title,
			&p.Email,
			&p.FirstName,
			&p.LastName,
			&p.AddressLine1,
			&p.AddressLine2,
			&p.City,
			&p.Postcode,
			&p.State,
			&p.CountryCode,
			&p.Phone,
			&p.CardNumber,
			&p.CardMonth,
			&p.CardYear,
			&p.CardCvv,
		)
		if err != nil {
			return nil, err
		}
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

	stmt := "update profiles set group_id = ?, title = ?, email = ?, first_name = ?, last_name = ?, address_line_1 = ?, address_line_2 = ?, city = ?, postcode = ?, state = ?, country_code = ?, phone = ?, card_number = ?, card_month = ?, card_year = ?, card_cvv = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.GroupID, updated.Title, updated.Email, updated.FirstName, updated.LastName, updated.AddressLine1, updated.AddressLine2, updated.City, updated.Postcode, updated.State, updated.CountryCode, updated.Phone, updated.CardNumber, updated.CardMonth, updated.CardYear, updated.CardCvv, id)
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
