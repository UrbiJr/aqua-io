package repository

import (
	"errors"

	"github.com/UrbiJr/copy-io/internal/user"
)

var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate() error
	InsertProfile(p user.Profile) (*user.Profile, error)
	InsertProfileGroup(pg user.ProfileGroup) (*user.ProfileGroup, error)
	AllProfiles() ([]user.Profile, error)
	AllProfileGroups() ([]user.ProfileGroup, error)
	UpdateProfile(id int64, updated user.Profile) error
	UpdateProfileGroup(id int64, updated user.ProfileGroup) error
	DeleteProfileGroup(id int64) error
	DeleteProfile(id int64) error
	DeleteProfilesByGroupID(id int64) error
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
