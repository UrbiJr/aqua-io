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
	InsertUser(u user.User) (*user.User, error)
	AllProfiles() ([]user.Profile, error)
	GetUser(ID int64) (*user.User, error)
	GetAllUsers() (*user.User, error)
	UpdateProfile(id int64, updated user.Profile) error
	UpdateUser(id int64, updated user.User) error
	DeleteProfile(id int64) error
	DeleteUser(id int64) error
	DeleteAllProfiles() error
	DeleteAllUsers() error

	InsertOpenedPosition(p user.OpenedPosition) (*user.OpenedPosition, error)
	AllOpenedPositions() ([]user.OpenedPosition, error)
	UpdateOpenedPosition(orderId string, updated user.OpenedPosition) error
	DeleteOpenedPosition(orderId string) error
}
