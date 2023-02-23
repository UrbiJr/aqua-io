package repository

import (
	"errors"

	"github.com/UrbiJr/nyx/internal/user"
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

	InsertProxy(p user.Proxy) (*user.Proxy, error)
	InsertProxyGroup(pg user.ProxyGroup) (*user.ProxyGroup, error)
	AllProxies() ([]user.Proxy, error)
	AllProxyGroups() ([]user.ProxyGroup, error)
	UpdateProxyGroup(id int64, updated user.ProxyGroup) error
	DeleteProxyGroup(id int64) error
	DeleteProxy(id int64) error
	DeleteProxiesByGroupID(id int64) error
}
