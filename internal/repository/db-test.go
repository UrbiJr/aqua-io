package repository

import (
	"github.com/UrbiJr/copy-io/internal/user"
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
		Title:          "john_red",
		OpenDelay:      10.5,
		BlacklistCoins: []string{"coin1", "coin2", "coin3"},
		TestMode:       true,
	}

	all = append(all, p)

	p = user.Profile{
		Title:          "james_green",
		OpenDelay:      5.7,
		BlacklistCoins: []string{"usdt"},
		TestMode:       false,
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
