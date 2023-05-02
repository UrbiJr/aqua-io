package repository

import (
	"errors"

	"github.com/UrbiJr/aqua-io/internal/user"
)

var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate() error
	InsertProfile(p user.Profile) (*user.Profile, error)
	AllProfiles() ([]user.Profile, error)
	UpdateProfile(id int64, updated user.Profile) error
	DeleteProfile(id int64) error
	DeleteAllProfiles() error
	InsertPosition(p user.Position) (*user.Position, error)
	InsertTrader(t user.Trader) (*user.Trader, error)
	AllPositions() ([]user.Position, error)
	AllTraders() ([]user.Trader, error)
	UpdatePosition(id int64, updated user.Position) error
	UpdateTrader(encryptedUid string, updated user.Trader) error
	DeletePosition(id int64) error
	DeletePositionsByTraderID(encryptedUid string) error
	DeleteTrader(encryptedUid string) error
}
