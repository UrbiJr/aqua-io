package user

// Profile contains information specific to a configuration for a trader
type Profile struct {
	ID                                   int64    `json:"id"`
	Title                                string   `json:"title"`
	TraderID                             string   `json:"trader_id"` // binance encrypted uid
	BybitApiKey                          string   `json:"bybit_api_key"`
	BybitApiSecret                       string   `json:"bybit_api_secret"`
	MaxBybitBinancePriceDifferentPercent float64  `json:"max_bybit_binance_price_difference_percent"`
	Leverage                             int64    `json:"leverage"`
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
	Profiles []Profile
}

func (pfm *ProfileManager) GetProfileByTitle(title string) *Profile {
	for _, p := range pfm.Profiles {
		if p.Title == title {
			return &p
		}
	}

	return nil
}

func (pfm *ProfileManager) UpdateProfile(ID int64, updated Profile) *Profile {
	for i, p := range pfm.Profiles {
		if p.ID == ID {
			pfm.Profiles[i].Title = updated.Title
			pfm.Profiles[i].TraderID = updated.TraderID
			pfm.Profiles[i].BybitApiKey = updated.BybitApiKey
			pfm.Profiles[i].BybitApiSecret = updated.BybitApiSecret
			pfm.Profiles[i].MaxBybitBinancePriceDifferentPercent = updated.MaxBybitBinancePriceDifferentPercent
			pfm.Profiles[i].Leverage = updated.Leverage
			pfm.Profiles[i].InitialOpenPercent = updated.InitialOpenPercent
			pfm.Profiles[i].MaxAddMultiplier = updated.MaxAddMultiplier
			pfm.Profiles[i].OpenDelay = updated.OpenDelay
			pfm.Profiles[i].OneCoinMaxPercent = updated.OneCoinMaxPercent
			pfm.Profiles[i].BlacklistCoins = updated.BlacklistCoins
			pfm.Profiles[i].MaxOpenPositions = updated.MaxOpenPositions
			pfm.Profiles[i].AutoTP = updated.AutoTP
			pfm.Profiles[i].AutoSL = updated.AutoSL
			pfm.Profiles[i].TestMode = updated.TestMode
			return &pfm.Profiles[i]
		}
	}

	return nil
}

func (pfm *ProfileManager) DeleteProfile(ID int64) {
	for i, p := range pfm.Profiles {
		if p.ID == ID {
			pfm.Profiles = append(pfm.Profiles[:i], pfm.Profiles[i+1:]...)
			break
		}
	}
}

func (pfm *ProfileManager) GetProfileByID(ID int64) *Profile {
	for _, p := range pfm.Profiles {
		if p.ID == ID {
			return &p
		}
	}

	return nil
}

func (pfm *ProfileManager) GetProfileByTraderID(traderID string) *Profile {
	for _, p := range pfm.Profiles {
		if p.TraderID == traderID {
			return &p
		}
	}

	return nil
}

func (pfm *ProfileManager) GetProfilesWithTrader() []Profile {
	var profiles []Profile
	for _, p := range pfm.Profiles {
		if p.TraderID != "" {
			profiles = append(profiles, p)
		}
	}

	return profiles
}

func (pfm *ProfileManager) GetAllTitles() []string {
	var titles []string
	for _, p := range pfm.Profiles {
		titles = append(titles, p.Title)
	}

	return titles
}
