package user

// Trader contains information specific to a trader (copied or not)
type Trader struct {
	ID                                int64    `json:"id"`
	ProfileID                         int64    `json:"profile_id"` // used to retrieve exchange settings
	EncryptedUid                      string   `json:"encrypted_uid"`
	FutureUid                         any      `json:"future_uid"`
	NickName                          string   `json:"nickName"`
	UserPhotoUrl                      string   `json:"userPhotoUrl"`
	Rank                              int64    `json:"rank"`
	Pnl                               float64  `json:"pnl"`
	Roi                               float64  `json:"roi"`
	PositionShared                    bool     `json:"positionShared"`
	TwitterUrl                        any      `json:"twitterUrl"`
	UpdateTime                        int64    `json:"updateTime"`
	FollowerCount                     int64    `json:"followerCount"`
	TwShared                          string   `json:"-"`
	IsTwTrader                        bool     `json:"isTwTrader"`
	OpenId                            any      `json:"openId"`
	PortfolioId                       any      `json:"portfolioId"`
	TradeMode                         string   `json:"trade_mode"`
	Leverage                          float64  `json:"leverage"`
	MaxOpenPositions                  int      `json:"max_open_positions"`
	MaxCoinPercentagePosition         int      `json:"max_coin_percentage_position"`
	MaxPriceDifferenceBetweenExchange float64  `json:"price_difference_between_exchanges"`
	OpenDelayBetweenPositions         int      `json:"open_delay_between_positions"`
	BlockPositionAdds                 bool     `json:"block_position_adds"`
	AutoTakeProfit                    Strategy `json:"-"`
	AutoStopLoss                      Strategy `json:"-"`
	MaxCoinAllocation                 float64  `json:"max_coin_allocation"` // 25%
	MaxAddMultiplier                  float64  `json:"max_add_multiplier"`  // 10%
	AddPreventionPercent              float64  `json:"add_prevention_percent"`
	BlackListedCoins                  []string `json:"blacklisted_coins"`
	StopControl                       bool     `json:"stop_control"`
}

type Strategy struct {
	ID           int64   `json:"id"`
	PositionSize float64 `json:"position_size"`
	Percentage   float64 `json:"percentage"`
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
	ProfileID        int64   `json:"-"`
	Symbol           string  `json:"symbol"`
	OrderID          string  `json:"orderId"`
	OrderLinkID      string  `json:"orderLinkId"`
	OrderStatus      string  `json:"orderStatus"`
	OrderType        string  `json:"orderType"`
	CreatedTime      int64   `json:"createdTime"`
	FilledQty        float64 `json:"cumExecQty"`
	Qty              float64 `json:"qty"`
	AvgFilledPrice   float64 `json:"avgPrice"`
	Price            float64 `json:"price"`
	TakeProfit       float64 `json:"takeProfit"`
	StopLoss         float64 `json:"stopLoss"`
	TriggerPrice     float64 `json:"triggerPrice"`
	TriggerDirection float64 `json:"triggerDirection"`
	Side             string  `json:"side"`
	IsLeverage       int64   `json:"isLeverage"`
}

type PositionInfo struct {
	OrderID        string  `json:"order_id"`
	PositionIdx    float64 `json:"positionIdx"`
	Side           string  `json:"side"`
	Symbol         string  `json:"symbol"`
	Size           float64 `json:"size"`
	AvgPrice       float64 `json:"avgPrice"` // Average entry price
	MarkPrice      float64 `json:"markPrice"`
	LiqPrice       float64 `json:"liqPrice"` // Position liquidation price
	UnrealisedPnl  float64 `json:"unrealisedPnl"`
	CumRealisedPnl float64 `json:"cumRealisedPnl"`
	Leverage       int64   `json:"leverage"`
	TakeProfit     any     `json:"takeProfit"`
	StopLoss       any     `json:"stopLoss"`
	PositionValue  float64 `json:"positionValue"`
	CreatedTime    int64   `json:"createdTime"`
	UpdatedTime    int64   `json:"updatedTime"`
	PositionStatus string  `json:"positionStatus"`
}

type CopiedTradersManager struct {
	CopiedTraders []Trader
}

func (ctm *CopiedTradersManager) GetCopiedTraderByID(ID int64) *Trader {
	for _, t := range ctm.CopiedTraders {
		if t.ID == ID {
			return &t
		}
	}
	return nil
}