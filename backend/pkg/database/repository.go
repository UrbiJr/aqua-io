package database

import (
	"errors"
	"github.com/UrbiJr/aqua-io/backend/internal/user"
)

var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate() error
	InsertProfile(p user.Profile) (*user.Profile, error)
	InsertStrategy(s user.Strategy) (*user.Strategy, error)
	InsertCopiedTrader(t user.Trader) (*user.Trader, error)
	InsertUser(u user.User) (*user.User, error)
	AllProfiles() ([]user.Profile, error)
	GetProfileByTitle(title string) (*user.Profile, error)
	AllStrategies() ([]user.Strategy, error)
	AllCopiedTraders() ([]user.Trader, error)
	AllUsers() ([]user.User, error)
	GetUser(ID int64) (*user.User, error)
	GetStrategy(ID int64) (*user.Strategy, error)
	UpdateProfile(ID int64, updated user.Profile) error
	UpdateStrategy(ID int64, updated user.Strategy) error
	UpdateCopiedTrader(ID int64, updated user.Trader) error
	UpdateUser(ID int64, updated user.User) error
	DeleteProfile(ID int64) error
	DeleteStrategy(ID int64) error
	DeleteCopiedTrader(ID int64) error
	DeleteUser(ID int64) error
	DeleteAllProfiles() error
	DeleteAllStrategies() error
	DeleteAllCopiedTraders() error
	DeleteAllUsers() error
}
