package constants

const (
	IncreaseSymbol = "↗"
	DecreaseSymbol = "↘"
)

// Exchange are the Exchanges we currently support
type Exchange string

const (
	ByBit   = "ByBit"
	OKX     = "OKX"
	Binance = "Binance"
	Phemex  = "Phemex"
	MEXC    = "MEXC"
	Bitget  = "Bitget"
	Bitmart = "Bitmart"
)

// 👇 Futures REST API.
const (
	BinancebaseRestAPI = "https://fapi.binance.com"
	MEXCbaseRestAPI    = "https://contract.mexc.com"
	ByBitbaseRestAPI   = "https://api.bybit.com"
	OKXbaseRestAPI     = "https://www.okx.com"
	PhemexbaseRestAPI  = "https://api.phemex.com"
	BitgetbaseRestAPI  = "https://api.bitget.com"
	BitmartbaseRestAPI = "https://api-cloud.bitmart.com"
)

// 👇 Futures WebSocket API.
const (
	BinancebaseSocketAPI = "wss://stream.binance.com:9443"
	MEXCbaseSocketAPI    = ""
	ByBitbaseSocketAPI   = ""
	OKXbaseSocketAPI     = ""
	BitgetbaseSocketAPI  = "wss://ws.bitget.com/mix/v1/stream"
	BitmartbaseSocketAPI = "wss://openapi-ws.bitmart.com/user?protocol=1.1"
)

// 👇 Exchange's Sides.
var (
	BINANCE_BUY        = "LONG"
	BINANCE_SELL       = "SHORT"
	MEXC_BUY           = ""
	MEX_SELL           = ""
	BYBIT_BUY          = "Buy"
	BYBIT_SELL         = "Sell"
	OKX_BUY            = "buy"
	OKX_SELL           = "sell"
	PHEMEX_BUY         = "Buy"
	PHEMEX_SELL        = "Sell"
	BITGET_BUY         = "open_long"
	BITGET_SELL        = "open_short"
	BITGET_BUY_CLOSE   = "close_long"
	BITGET_SELL_CLOSE  = "close_short"
	BITMART_BUY        = 1
	BITMART_SELL       = 2
	BITMART_BUY_CLOSE  = 3
	BITMART_SELL_CLOSE = 4
)
