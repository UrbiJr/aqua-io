package repository

import (
	"testing"

	"github.com/UrbiJr/nyx/internal/user"
)

func TestSQLiteRepository_Migrate(t *testing.T) {
	err := testRepo.Migrate()
	if err != nil {
		t.Error("migrate failed:", err)
	}
}

func TestSQLiteRepository_InsertProfile(t *testing.T) {
	p := user.Profile{
		Title:          "jack_ford",
		OpenDelay:      10.5,
		BlacklistCoins: []string{"coin1", "coin2", "coin3"},
		TestMode:       true,
	}

	result, err := testRepo.InsertProfile(p)
	if err != nil {
		t.Error("insert failed:", err)
	}

	if result.ID <= 0 {
		t.Error("invalid id sent back:", result.ID)
	}
}

func TestSQLiteRepository_AllProfiles(t *testing.T) {
	p, err := testRepo.AllProfiles()
	if err != nil {
		t.Error("get all failed:", err)
	}

	// we inserted 1 row in TestSQLiteRepository_InsertProfile
	if len(p) != 1 {
		t.Error("wrong number of rows returned; expected 1, got:", len(p))
	}
}

func TestSQLiteRepository_UpdateProfile(t *testing.T) {
	p, err := testRepo.AllProfiles()
	if err != nil {
		t.Error("get all failed:", err)
	}

	p[0].OpenDelay = 12
	id := p[0].ID
	err = testRepo.UpdateProfile(id, p[0])
	if err != nil {
		t.Error("update failed:", err)
	}

	p, err = testRepo.AllProfiles()
	if err != nil {
		t.Error("get all failed:", err)
	}

	found := false
	for _, p := range p {
		if p.ID == id {
			found = true
			if p.OpenDelay != 12 {
				t.Errorf("updated failed, expected open_delay 12, got: %f", p.OpenDelay)
			}
		}
	}

	if !found {
		t.Error("get all failed: updated profile not found")
	}

}

func TestSQLiteRepository_DeleteProfile(t *testing.T) {
	err := testRepo.DeleteProfile(1)
	if err != nil {
		t.Error("failed to delete profile", err)
		if err != errDeleteFailed {
			t.Error("wrong error returned")
		}
	}

	err = testRepo.DeleteProfile(299)
	if err == nil {
		t.Error("no error when trying to delete non-existent record")
	}
}
