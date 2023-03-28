package user

type ProfileGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Profile contains information specific to a single account of a particular site i.e. BSTN
type Profile struct {
	ID                                   int64    `json:"id"`
	GroupID                              int64    `json:"group_id"`
	Title                                string   `json:"title"`
	BybitApiKey                          string   `json:"bybit_api_key"`
	BybitApiSecret                       string   `json:"bybit_api_secret"`
	MaxBybitBinancePriceDifferentPercent float64  `json:"max_bybit_binance_price_difference_percent"`
	Leverage                             float64  `json:"leverage"`
	InitialOpenPercent                   float64  `json:"initial_open_percent"`
	MaxAddMultiplier                     float64  `json:"max_add_multiplier"`
	OpenDelay                            float64  `json:"open_delay"`
	OneCoinMaxPercent                    float64  `json:"one_coin_max_percent"`
	BlacklistCoins                       []string `json:"blacklist_coins"` // stored in DB as comma separated: coin1,coin2,coin3
	AddPreventionPercent                 float64  `json:"add_prevention_percent"`
	BlockAddsAboveEntry                  bool     `json:"block_adds_above_entry"`
	MaxOpenPositions                     int64    `json:"max_open_positions"`
	AutoTP                               float64  `json:"auto_tp"`
	AutoSL                               float64  `json:"auto_sl"`
	TestMode                             bool     `json:"test_mode"`
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

func (pfm *ProfileManager) GetGroupByName(name string) *ProfileGroup {
	for _, p := range pfm.Groups {
		if p.Name == name {
			return &p
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

func (pfm *ProfileManager) GetProfileByID(ID int64, groupID int64) *Profile {
	filtered := pfm.FilterByGroupID(groupID)
	for _, p := range filtered {
		if p.ID == ID {
			return &p
		}
	}

	return nil
}
