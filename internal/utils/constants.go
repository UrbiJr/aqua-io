package utils

type PositionDirection string
type ByBitProductType string
type OrderSide string
type OrderType string

const (
	SHORT_POSITION       = "Short"
	LONG_POSITION        = "Long"
	BYBIT_PRODUCT_LINEAR = "linear"
	BYBIT_PRODUCT_SPOT   = "spot"
	ORDER_SELL           = "Sell"
	ORDER_BUY            = "Buy"
	ORDER_MARKET         = "Market"
	ORDER_LIMIT          = "Limit"
)
