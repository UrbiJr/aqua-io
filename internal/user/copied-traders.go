package user

// Trader contains information specific to a trader
type Trader struct {
	EncryptedUid   string  `json:"encryptedUid"`
	FutureUid      any     `json:"futureUid"`
	NickName       string  `json:"nickName"`
	UserPhotoUrl   string  `json:"userPhotoUrl"`
	Rank           int64   `json:"rank"`
	Pnl            float64 `json:"pnl"`
	Roi            float64 `json:"roi"`
	PositionShared bool    `json:"positionShared"`
	TwitterUrl     any     `json:"twitterUrl"`
	UpdateTime     int64   `json:"updateTime"`
	FollowerCount  int64   `json:"followerCount"`
	TwShared       string  `json:"-"`
	IsTwTrader     bool    `json:"isTwTrader"`
	OpenId         any     `json:"openId"`
	PortfolioId    any     `json:"portfolioId"`
}

type Position struct {
	ID              int64   `json:"id"`
	Symbol          string  `json:"symbol"`
	EntryPrice      float64 `json:"entryPrice"`
	MarkPrice       float64 `json:"markPrice"`
	Pnl             float64 `json:"pnl"`
	Roe             float64 `json:"roe"`
	UpdateTime      []int64 `json:"-"`
	Amount          float64 `json:"amount"`
	UpdateTimestamp int64   `json:"updateTimeStamp"`
	Yellow          bool    `json:"yellow"`
	TradeBefore     bool    `json:"tradeBefore"`
	Leverage        int64   `json:"leverage"`
}

type Transaction struct {
	ProfileID       int64   `json:"-"`
	OrderID         string  `json:"orderId"`
	TradeID         string  `json:"tradeId"`
	Symbol          string  `json:"symbol"`
	Currency        string  `json:"currency"`
	Funding         float64 `json:"funding"`
	TradePrice      float64 `json:"tradePrice"`
	Qty             float64 `json:"qty"`
	Size            float64 `json:"size"`
	Side            string  `json:"side"`
	TransactionTime int64   `json:"transactionTime"`
}

type Order struct {
	ProfileID   int64   `json:"-"`
	Symbol      string  `json:"symbol"`
	OrderID     string  `json:"orderId"`
	OrderLinkID string  `json:"orderLinkId"`
	OrderStatus string  `json:"orderStatus"`
	OrderType   string  `json:"orderType"`
	CreatedTime int64   `json:"createdTime"`
	Qty         float64 `json:"qty"`
	Price       float64 `json:"price"`
	Side        string  `json:"side"`
	IsLeverage  int64   `json:"isLeverage"`
}

type PositionInfo struct {
	PositionIdx    float64 `json:"positionIdx"`
	Symbol         string  `json:"symbol"`
	Leverage       int64   `json:"leverage"`
	AvgPrice       float64 `json:"avgPrice"` // Average entry price
	LiqPrice       float64 `json:"liqPrice"` // Position liquidation price
	TakeProfit     any     `json:"takeProfit"`
	StopLoss       any     `json:"stopLoss"`
	PositionValue  float64 `json:"positionValue"`
	UnrealisedPnl  float64 `json:"unrealisedPnl"`
	CumRealisedPnl float64 `json:"cumRealisedPnl"`
	MarkPrice      float64 `json:"markPrice"`
	CreatedTime    int64   `json:"createdTime"`
	UpdatedTime    int64   `json:"updatedTime"`
	Side           string  `json:"side"`
	PositionStatus string  `json:"positionStatus"`
}
