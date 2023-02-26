package user

type ProfileGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Profile contains information specific to a single account of a particular site i.e. BSTN
type Profile struct {
	ID           int64  `json:"id"`
	GroupID      int64  `json:"group_id"`
	Title        string `json:"title"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	Postcode     string `json:"postcode"`
	State        string `json:"state"`
	CountryCode  string `json:"country_code"`
	Phone        string `json:"phone"`
	CardNumber   string `json:"card_number"`
	CardMonth    string `json:"card_month"`
	CardYear     string `json:"card_year"`
	CardCvv      string `json:"card_cvv"`
}

type ProfileManager struct {
	Groups   []ProfileGroup
	Profiles []Profile
}

func (pfm *ProfileManager) FilterByGroupName(groupName string) []Profile {
	var filtered []Profile

	// iterate through proxy groups
	for _, g := range pfm.Groups {
		// if group is the one requested
		if g.Name == groupName {
			// keep it as current and iterate through proxies
			for _, p := range pfm.Profiles {
				// if proxy belongs to current group
				if p.GroupID == g.ID {
					// add it to filtered
					filtered = append(filtered, p)
				}
			}
			return filtered
		}
	}

	return filtered
}

func (pfm *ProfileManager) FilterByGroupID(ID int64) []Profile {
	var filtered []Profile

	// keep it as current and iterate through proxies
	for _, p := range pfm.Profiles {
		// if proxy belongs to current group
		if p.GroupID == ID {
			// add it to filtered
			filtered = append(filtered, p)
		}
	}

	return filtered
}

func (pfm *ProfileManager) GetGroupByID(ID int64) *ProfileGroup {
	for idx, g := range pfm.Groups {
		if g.ID == ID {
			return &pfm.Groups[idx]
		}
	}

	return nil
}

func (pfm *ProfileManager) GetProfileByTitle(title string, groupID int64) *Profile {
	filtered := pfm.FilterByGroupID(groupID)
	for _, p := range filtered {
		if p.Title == title {
			return &p
		}
	}

	return nil
}
