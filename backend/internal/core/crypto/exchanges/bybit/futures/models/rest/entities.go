package rest

type Position struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		OrderId     string `json:"orderId"`
		OrderLinkId string `json:"orderLinkId"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type Leverage struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type Balance struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		List []struct {
			AccountType            string `json:"accountType"`
			AccountIMRate          string `json:"accountIMRate"`
			AccountMMRate          string `json:"accountMMRate"`
			TotalEquity            string `json:"totalEquity"`
			TotalWalletBalance     string `json:"totalWalletBalance"`
			TotalMarginBalance     string `json:"totalMarginBalance"`
			TotalAvailableBalance  string `json:"totalAvailableBalance"`
			TotalPerpUPL           string `json:"totalPerpUPL"`
			TotalInitialMargin     string `json:"totalInitialMargin"`
			TotalMaintenanceMargin string `json:"totalMaintenanceMargin"`
			AccountLTV             string `json:"accountLTV"`
			Coin                   []struct {
				Coin                string `json:"coin"`
				Equity              string `json:"equity"`
				UsdValue            string `json:"usdValue"`
				WalletBalance       string `json:"walletBalance"`
				BorrowAmount        string `json:"borrowAmount"`
				AvailableToBorrow   string `json:"availableToBorrow"`
				AvailableToWithdraw string `json:"availableToWithdraw"`
				AccruedInterest     string `json:"accruedInterest"`
				TotalOrderIM        string `json:"totalOrderIM"`
				TotalPositionIM     string `json:"totalPositionIM"`
				TotalPositionMM     string `json:"totalPositionMM"`
				UnrealisedPnl       string `json:"unrealisedPnl"`
				CumRealisedPnl      string `json:"cumRealisedPnl"`
			} `json:"coin"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}
