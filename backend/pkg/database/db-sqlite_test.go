package database

import (
	"testing"

	"github.com/UrbiJr/aqua-io/backend/internal/user"
)

func TestSQLiteRepository_Migrate(t *testing.T) {
	err := testRepo.Migrate()
	if err != nil {
		t.Error("migrate failed:", err)
	}
}

func TestSQLiteRepository_InsertProfile(t *testing.T) {
	bAddress, err := testRepo.InsertAddress(
		user.Address{
			FirstName:    "Federico",
			LastName:     "Urbinelli",
			Phone:        "1112221212",
			AddressLine1: "Piazza Roma 1",
			AddressLine2: "",
			ZipCode:      "48015",
			Province:     "RA",
			CountryCode:  "IT",
		})
	if err != nil {
		t.Error("insert billing address failed:", err)
	}

	sAddress, err := testRepo.InsertAddress(
		user.Address{
			FirstName:    "Federico",
			LastName:     "Urbinelli",
			Phone:        "1112221212",
			AddressLine1: "Via Molveno 6",
			AddressLine2: "Netrising s.r.l.",
			ZipCode:      "48015",
			Province:     "RA",
			CountryCode:  "IT",
		})
	if err != nil {
		t.Error("insert shipping address failed:", err)
	}

	p := user.Profile{
		Title:             "federico 1",
		BillingAddressID:  bAddress.ID,
		ShippingAddressID: sAddress.ID,
		CardNumber:        "",
		CardMonth:         "",
		CardYear:          "",
		CardCvv:           "",
		TestMode:          true,
	}

	result, err := testRepo.InsertProfile(p)
	if err != nil {
		t.Error("insert profile failed:", err)
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

	p[0].Title = "Personal 1"
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
			if p.Title != "Personal 1" {
				t.Errorf("updated failed, expected title 'Personal 1', got: %s", p.Title)
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
