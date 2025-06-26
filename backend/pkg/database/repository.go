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
	InsertAddress(p user.Address) (*user.Address, error)
	InsertProxyList(s user.ProxyList) (*user.ProxyList, error)
	InsertTask(t user.Task) (*user.Task, error)
	InsertUser(u user.User) (*user.User, error)
	AllProfiles() ([]user.Profile, error)
	GetProfileByTitle(title string) (*user.Profile, error)
	GetProfileByID(ID int64) (*user.Profile, error)
	GetAddressByID(ID int64) (*user.Address, error)
	GetTaskByID(ID int64) (*user.Task, error)
	AllAddresses() ([]user.Address, error)
	AllProxyLists() ([]user.ProxyList, error)
	AllTasks() ([]user.Task, error)
	AllUsers() ([]user.User, error)
	GetUser(ID int64) (*user.User, error)
	GetProxyList(ID int64) (*user.ProxyList, error)
	UpdateProfile(ID int64, updated user.Profile) error
	UpdateAddress(ID int64, updated user.Address) error
	UpdateProxyList(ID int64, updated user.ProxyList) error
	UpdateTask(ID int64, updated user.Task) error
	UpdateUser(ID int64, updated user.User) error
	DeleteAddress(ID int64) error
	DeleteProfile(ID int64) error
	DeleteProxyList(ID int64) error
	DeleteTask(ID int64) error
	DeleteUser(ID int64) error
	DeleteAllAddresses() error
	DeleteAllProfiles() error
	DeleteAllProxyLists() error
	DeleteAllTasks() error
	DeleteAllUsers() error
}
