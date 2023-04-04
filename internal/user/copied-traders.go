package user

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
}

type Position struct {
	ID              int64   `json:"id"`
	TraderID        int64   `json:"trader_id"`
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
	ProfileGroupID  int64   `json:"-"`
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

type CopiedTradersManager struct {
	Traders []Trader
}

func (ctm *CopiedTradersManager) GetTraderByUid(encryptedUid string) *Trader {
	for idx, t := range ctm.Traders {
		if t.EncryptedUid == encryptedUid {
			return &ctm.Traders[idx]
		}
	}

	return nil
}

func (ctm *CopiedTradersManager) RemoveTraderByUid(encryptedUid string) {
	for idx, t := range ctm.Traders {
		if t.EncryptedUid == encryptedUid {
			ctm.Traders = append(ctm.Traders[:idx], ctm.Traders[idx+1:]...)
		}
	}
}
