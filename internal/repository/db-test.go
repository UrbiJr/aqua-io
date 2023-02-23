package repository

import (
	"github.com/UrbiJr/nyx/internal/user"
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
		Title:     "john_red",
		Email:     "john@nyx-robotics.eu",
		FirstName: "John",
		LastName:  "Red",
	}

	all = append(all, p)

	p = user.Profile{
		Title:     "james_green",
		Email:     "james@nyx-robotics.eu",
		FirstName: "James",
		LastName:  "Green",
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
