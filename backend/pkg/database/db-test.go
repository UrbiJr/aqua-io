package database

import (
	"github.com/UrbiJr/aqua-io/backend/internal/user"
)

type TestRepository struct{}

func NewTestRepository() *TestRepository {
	return &TestRepository{}
}

func (repo *TestRepository) Migrate() error {
	return nil
}

func (repo *TestRepository) InsertProfile(p user.Profile) (*user.Profile, error) {
	return &p, nil
}

func (repo *TestRepository) AllProfiles() ([]user.Profile, error) {
	var all []user.Profile

	p := user.Profile{
		Title:           "john_red",
		StopIfFallUnder: 1000,
		AccountName:     "johnny1",
		TestMode:        true,
	}

	all = append(all, p)

	p = user.Profile{
		Title:           "james_green",
		StopIfFallUnder: 2000,
		AccountName:     "james-1",
		TestMode:        true,
	}

	all = append(all, p)

	return all, nil
}

func (repo *TestRepository) UpdateProfile(id int64, updated user.Profile) error {
	return nil
}

func (repo *TestRepository) DeleteProfile(id int64) error {
	return nil
}
