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
			billing_address_id integer null,
			shipping_address_id integer null,
			card_number text null,
			card_month text null,
			card_year text null,
			card_cvv text null,
			test_mode boolean null,
			FOREIGN KEY(billing_address_id) REFERENCES addresses(id),
			FOREIGN KEY(shipping_address_id) REFERENCES addresses(id));
	`
	_, err := repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists addresses(
			id integer primary key autoincrement,
			first_name text null,
			last_name text null,
			phone text null,
			address_line_1 text null,
			address_line_2 text null,
			zip_code text null,
			province text null,
			country_code text null);
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists proxy_lists(
			id integer primary key autoincrement,
			title text not null,
			proxies text null);
	`
	_, err = repo.Conn.Exec(query)
	if err != nil {
		return err
	}

	query = `
		create table if not exists tasks(
			id string primary key,
			profile_id integer null,
			proxy_list_id integer null,
			module text null,
			payment_mode text null,
			FOREIGN KEY(profile_id) REFERENCES profiles(id),
			FOREIGN KEY(proxy_list_id) REFERENCES proxy_lists(id));
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
	stmt := `
		insert into profiles (
			title,
			billing_address_id,
			shipping_address_id,
			card_number,
			card_month,
			card_year,
			card_cvv,
			test_mode
		) values (?, ?, ?, ?, ?, ?, ?, ?)
	`

	var testMode int
	if p.TestMode {
		testMode = 1
	} else {
		testMode = 0
	}

	res, err := repo.Conn.Exec(stmt, p.Title, p.BillingAddressID, p.ShippingAddressID, p.CardNumber, p.CardMonth, p.CardYear, p.CardCvv, testMode)
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

func (repo *SQLiteRepository) InsertAddress(p user.Address) (*user.Address, error) {
	stmt := `
		insert into addresses (
			first_name,
			last_name,
			phone,
			address_line_1,
			address_line_2,
			zip_code,
			province,
			country_code
		) values (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := repo.Conn.Exec(stmt, p.FirstName, p.LastName, p.Phone, p.AddressLine1, p.AddressLine2, p.ZipCode, p.Province, p.CountryCode)
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

func (repo *SQLiteRepository) InserProxyList(s user.ProxyList) (*user.ProxyList, error) {
	stmt := `
		insert into proxy_lists (
			title,
			proxies
		) values (?, ?)
	`

	var proxies string
	proxies = strings.Join(s.Proxies, ",")

	res, err := repo.Conn.Exec(stmt, s.Title, proxies)
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

func (repo *SQLiteRepository) InserTask(t user.Task) (*user.Task, error) {
	stmt := `
		insert into tasks (
		profile_id,
		proxy_list_id,
		module, 
		payment_mode
	) values (?, ?, ?, ?)
	`
	res, err := repo.Conn.Exec(stmt,
		t.ProfileID,
		t.ProxyListID,
		t.Module,
		t.PaymentMode,
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
	stmt := "insert into users (profile_picture_path, license_key, persistent_login, theme) values (?, ?, ?, ?)"
	var persistent, closeAllTradesWhenClosing int

	if u.PersistentLogin {
		persistent = 1
	} else {
		persistent = 0
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
		billing_address_id,
		shipping_address_id,
		card_number,
		card_month,
		card_year,
		card_cvv,
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
			&p.BillingAddressID,
			&p.ShippingAddressID,
			&p.CardNumber,
			&p.CardMonth,
			&p.CardYear,
			&p.CardCvv,
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

func (repo *SQLiteRepository) AllAddresses() ([]user.Address, error) {
	query := `
		select id,
		first_name,
		last_name,
		phone,
		address_line_1,
		address_line_2,
		zip_code,
		province,
		country_code
	 	from addresses order by title
	`

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Address
	for rows.Next() {
		var p user.Address

		err := rows.Scan(
			&p.ID,
			&p.FirstName,
			&p.LastName,
			&p.Phone,
			&p.AddressLine1,
			&p.AddressLine2,
			&p.ZipCode,
			&p.Province,
			&p.CountryCode,
		)
		if err != nil {
			return nil, err
		}

		all = append(all, p)
	}

	return all, nil
}

func (repo *SQLiteRepository) GetProfileByTitle(title string) (*user.Profile, error) {
	query := "select * from profiles where title = ? limit 1"

	rows, err := repo.Conn.Query(query, title)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var p *user.Profile
	for rows.Next() {
		var testMode int
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.BillingAddressID,
			&p.ShippingAddressID,
			&p.CardNumber,
			&p.CardMonth,
			&p.CardYear,
			&p.CardCvv,
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
	}

	return p, nil
}

func (repo *SQLiteRepository) GetProfileByID(ID int64) (*user.Profile, error) {
	query := "select * from profiles where id = ? limit 1"

	rows, err := repo.Conn.Query(query, ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var p *user.Profile
	for rows.Next() {
		var testMode int
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.BillingAddressID,
			&p.ShippingAddressID,
			&p.CardNumber,
			&p.CardMonth,
			&p.CardYear,
			&p.CardCvv,
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
	}

	return p, nil
}

func (repo *SQLiteRepository) AllProxyLists() ([]user.ProxyList, error) {
	query := `
		select id,
		title,
		proxies
	 	from proxy_lists
	`

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.ProxyList
	for rows.Next() {
		var s user.ProxyList

		err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Proxies,
		)
		if err != nil {
			return nil, err
		}

		all = append(all, s)
	}

	return all, nil
}

func (repo *SQLiteRepository) AllTasks() ([]user.Task, error) {
	query := `
		select id,
		profile_id,
		proxy_list_id,
		module, 
		payment_mode
	 	from tasks order by id
	 `

	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.Task
	for rows.Next() {
		var t user.Task

		err := rows.Scan(
			&t.ID,
			&t.ProfileID,
			&t.ProxyListID,
			&t.Module,
			&t.PaymentMode,
		)
		if err != nil {
			return nil, err
		}

		all = append(all, t)
	}

	return all, nil
}

func (repo *SQLiteRepository) AllUsers() ([]user.User, error) {
	query := "select id, profile_picture_path, license_key, persistent_login, theme from users"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []user.User
	for rows.Next() {
		var u user.User
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

		all = append(all, u)
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

func (repo *SQLiteRepository) GetProxyList(ID int64) (*user.ProxyList, error) {
	query := "select id, title, proxies from proxy_lists where id = ?"

	rows, err := repo.Conn.Query(query, ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var s *user.ProxyList
	for rows.Next() {
		err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Proxies,
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
			billing_address_id = ?,
			shipping_address_id = ?,
			card_number = ?,
			card_month = ?,
			card_year = ?,
			card_cvv = ?,
			test_mode = ?
		where id = ?
	`
	res, err := repo.Conn.Exec(stmt, updated.Title, updated.BillingAddressID, updated.ShippingAddressID, updated.CardNumber, updated.CardMonth, updated.CardYear, updated.CardCvv, testMode, ID)
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

func (repo *SQLiteRepository) UpdateAddress(ID int64, updated user.Address) error {
	if ID <= 0 {
		return errors.New("invalid updated id")
	}

	stmt := `
		update addresses set
			first_name = ?,
			last_name = ?,
			phone = ?,
			address_line_1 = ?,
			address_line_2 = ?,
			zip_code = ?,
			province = ?,
			country_code = ?,
		where id = ?
	`
	res, err := repo.Conn.Exec(stmt, updated.FirstName, updated.LastName, updated.Phone, updated.AddressLine1, updated.AddressLine2, updated.ZipCode, updated.Province, updated.CountryCode, ID)
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

func (repo *SQLiteRepository) UpdateProxyList(ID int64, updated user.ProxyList) error {
	if ID <= 0 {
		return errors.New("invalid updated id")
	}

	stmt := `
		update proxy_lists set
			title = ?,
			proxies = ?,
		where id = ?
	`

	var proxies string

	proxies = strings.Join(updated.Proxies, ",")

	res, err := repo.Conn.Exec(stmt, updated.Title, proxies, ID)
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

func (repo *SQLiteRepository) UpdateTask(ID int64, updated user.Task) error {
	if ID <= 0 {
		return errors.New("invalid updated id")
	}

	stmt := `
		update copied_traders set
		profile_id = ?, 
		proxy_list_id = ?,
		module = ?, 
		payment_mode = ? 
	where id  = ?
	`
	res, err := repo.Conn.Exec(stmt,
		updated.ProfileID,
		updated.ProxyListID,
		updated.Module,
		updated.PaymentMode,
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

	var persistent int

	if updated.PersistentLogin {
		persistent = 1
	} else {
		persistent = 0
	}

	stmt := "update users set profile_picture_path = ?, license_key = ?,  persistent_login = ?, theme = ?, where id = ?"
	res, err := repo.Conn.Exec(stmt, updated.ProfilePicturePath, updated.LicenseKey, persistent, updated.Theme, ID)
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

func (repo *SQLiteRepository) DeleteAddress(ID int64) error {
	res, err := repo.Conn.Exec("delete from addresses where id = ?", ID)
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

func (repo *SQLiteRepository) DeleteProxyList(ID int64) error {
	res, err := repo.Conn.Exec("delete from proxy_lists where id = ?", ID)
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

func (repo *SQLiteRepository) DeleteTask(ID int64) error {
	res, err := repo.Conn.Exec("delete from tasks where id = ?", ID)
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

func (repo *SQLiteRepository) DeleteAllAddresses() error {
	res, err := repo.Conn.Exec("delete from addresses where 1 = 1")
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

func (repo *SQLiteRepository) DeleteAllProxyLists() error {
	res, err := repo.Conn.Exec("delete from proxy_lists where 1 = 1")
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

func (repo *SQLiteRepository) DeleteAllTasks() error {
	res, err := repo.Conn.Exec("delete from tasks where 1 = 1")
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
