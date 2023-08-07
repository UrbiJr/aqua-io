package database

import (
	"errors"
	user2 "github.com/UrbiJr/aqua-io/backend/internal/user"
)

var (
	errUpdateFailed = errors.New("update failed")
	errDeleteFailed = errors.New("delete failed")
)

type Repository interface {
	Migrate() error
	InsertProfile(p user2.Profile) (*user2.Profile, error)
	InsertUser(u user2.User) (*user2.User, error)
	AllProfiles() ([]user2.Profile, error)
	GetUser(ID int64) (*user2.User, error)
	GetAllUsers() (*user2.User, error)
	UpdateProfile(id int64, updated user2.Profile) error
	UpdateUser(id int64, updated user2.User) error
	DeleteProfile(id int64) error
	DeleteUser(id int64) error
	DeleteAllProfiles() error
	DeleteAllUsers() error

	InsertOpenedPosition(p user2.OpenedPosition) (*user2.OpenedPosition, error)
	AllOpenedPositions() ([]user2.OpenedPosition, error)
	UpdateOpenedPosition(orderId string, updated user2.OpenedPosition) error
	DeleteOpenedPosition(orderId string) error
}
