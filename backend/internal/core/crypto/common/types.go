package common

import (
	"github.com/bogdanfinn/tls-client"
)

// ExchangeTrade is the trade data sent by the websocket.
type ExchangeTrade struct {
	Event OrderType        `json:"event"`
	Trade TradeInformation `json:"trade"`
}

type TradeInformation struct {
	Symbol     string  `json:"symbol"`
	EntryPrice float64 `json:"entryPrice"`
	MarkPrice  float64 `json:"markPrice"`
	Amount     float64 `json:"amount"`
	Leverage   int     `json:"leverage"`
	Pnl        float64 `json:"pnl"`
	Roe        float64 `json:"roe"`
	Side       string  `json:"side"`
	TraderName string  `json:"traderName"`
	TraderUID  string  `json:"traderUID"`
}

type OrderType string

const (
	OpenPos     OrderType = "newTrade"
	ClosePos    OrderType = "closedTrade"
	UpdatePos   OrderType = "updateTrade"
	IncreasePos OrderType = "increaseTrade"
	DecreasePos OrderType = "decreaseTrade"
)

// Trade is the trade data we store on our side.
type Trade struct {
	Event         OrderType
	ExchangeTrade ExchangeTrade
	Info          Info
	Misc          Misc
}

type Info struct {
	TradeSettings TradeSettings
	Coin          Coin
	Wallet        Wallet

	TradeID    string
	TraderID   string
	TraderName string
	TradeMode  string
	Side       string
	Quantity   float64
	Leverage   float64
	PnL        int

	Metrics Metrics
}

// TradeSettings & below comes from the user's settings in config.yml file.
type TradeSettings struct {
	MaxPriceDifferenceBetweenExchange float64
	MaxCoinAllocation                 float64
	MaxOpenPositions                  int
	InitialOpenPercentage             float64
	OpenDelayBetweenPositions         OpenDelayBetweenPositions
	MaxAddMultiplier                  MaxAddMultiplier
	BlockPositionAdds                 bool
	AddPreventionPercent              float64
	StopControl                       bool
}

type MaxAddMultiplier struct {
	Value     float64
	ConfigVal float64
}

type Coin struct {
	Symbol        string
	ExchangePrice float64
	BinancePrice  float64
}

// Misc contains exchange's API keys (passphrase is needed for some). Index is the config's
// number when you iterate.
type Misc struct {
	PublicAPI  string
	SecretAPI  string
	Passphrase string
	Index      int
	Client     tls_client.HttpClient
}

type OpenDelayBetweenPositions struct {
	Value          int
	SleepUntil     int64
	LastTimeAction int64
}

type Wallet struct {
	InitialBalance float64
	Value          float64
}

// Metrics is unused for now
type Metrics struct {
	AverageEntryPrice float64
	PositionChange    map[int64]float64
}
