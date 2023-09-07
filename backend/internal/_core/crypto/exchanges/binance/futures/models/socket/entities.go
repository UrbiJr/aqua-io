package socket

type PositionsPayload struct {
	ID     string                 `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type Positions struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Result []struct {
		Id           int    `json:"id"`
		Price        string `json:"price"`
		Qty          string `json:"qty"`
		QuoteQty     string `json:"quoteQty"`
		Time         int64  `json:"time"`
		IsBuyerMaker bool   `json:"isBuyerMaker"`
		IsBestMatch  bool   `json:"isBestMatch"`
	} `json:"result"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
		Count         int    `json:"count"`
	} `json:"rateLimits"`
}

type BalancePayload struct {
	ID     string                 `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type Balances struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Result struct {
		MakerCommission  int  `json:"makerCommission"`
		TakerCommission  int  `json:"takerCommission"`
		BuyerCommission  int  `json:"buyerCommission"`
		SellerCommission int  `json:"sellerCommission"`
		CanTrade         bool `json:"canTrade"`
		CanWithdraw      bool `json:"canWithdraw"`
		CanDeposit       bool `json:"canDeposit"`
		CommissionRates  struct {
			Maker  string `json:"maker"`
			Taker  string `json:"taker"`
			Buyer  string `json:"buyer"`
			Seller string `json:"seller"`
		} `json:"commissionRates"`
		Brokered                   bool   `json:"brokered"`
		RequireSelfTradePrevention bool   `json:"requireSelfTradePrevention"`
		UpdateTime                 int64  `json:"updateTime"`
		AccountType                string `json:"accountType"`
		Balances                   []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
		Permissions []string `json:"permissions"`
	} `json:"result"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
		Count         int    `json:"count"`
	} `json:"rateLimits"`
}

// AccountUpdate is the Socket Response for event.ACCOUNT_UPDATE
type AccountUpdate struct {
	E int64   `json:"E"`
	T int64   `json:"T"`
	A Account `json:"a"`
}

type Account struct {
	M string     `json:"m"`
	B []Balance  `json:"B"`
	P []Position `json:"P"`
}

type Balance struct {
	A  string `json:"a"`
	WB string `json:"wb"`
	CW string `json:"cw"`
	BC string `json:"bc"`
}

type Position struct {
	S  string `json:"s"`
	PA string `json:"pa"`
	EP string `json:"ep"`
	CR string `json:"cr"`
	UP string `json:"up"`
	MT string `json:"mt"`
	IW string `json:"iw"`
	PS string `json:"ps"`
}
