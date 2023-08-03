package hub

import (
	"fmt"
	"github.com/UrbiJr/aqua-io/internal/client"
	"github.com/UrbiJr/aqua-io/internal/core/crypto/common"
	"github.com/UrbiJr/aqua-io/internal/core/crypto/constants"
	bsc "github.com/UrbiJr/aqua-io/internal/core/crypto/exchanges/binance/futures/client/rest"
	bit "github.com/UrbiJr/aqua-io/internal/core/crypto/exchanges/bybit/futures/client/rest"
	okx "github.com/UrbiJr/aqua-io/internal/core/crypto/exchanges/okx/futures/client/rest"
	"github.com/UrbiJr/aqua-io/internal/utils/safemap"
	"github.com/bogdanfinn/tls-client"
	"strings"
)

/*
type Trade struct {
	Trade *common.Trade
}
*/

func NewTrade(b *common.ExchangeTrade) *common.Trade {
	return &common.Trade{
		Event:         b.Event,
		ExchangeTrade: *b,
		Info: common.Info{
			TraderID:   b.Trade.TraderUID,
			TraderName: b.Trade.TraderName,
			Side:       strings.ToLower(b.Trade.Side),
			Coin: common.Coin{
				Symbol:       b.Trade.Symbol,
				BinancePrice: b.Trade.EntryPrice,
			},
		},
		Misc: common.Misc{
			Client: client.NewProxyLess(),
		},
	}
}

var (
	OKX_ORDERBOOK Orderbook[okx.OKX]
	BSC_ORDERBOOK Orderbook[bsc.Binance]
	BIT_ORDERBOOK Orderbook[bit.ByBit]
)

type Orderbook[T any] struct {
	Exchange    string
	Trades      *safemap.SafeMap[string, *T] //symbol:rest
	TradeAmount *safemap.SafeMap[string, int]
	Client      tls_client.HttpClient
}

func NewOrderbook[T any]() Orderbook[T] {
	return Orderbook[T]{
		Trades:      safemap.New[string, *T](),
		TradeAmount: safemap.New[string, int](),
		Client:      client.NewProxyLess(),
	}
}

// GetTrade returns a trade (used for updating/closing a position).
func (ob *Orderbook[T]) GetTrade(symbol string) *T {
	trade, ok := ob.Trades.Get(symbol)
	if !ok {
		return trade
	}
	return trade
}

func (ob *Orderbook[T]) AddTrade(symbol string, t *T) {
	ob.Trades.Set(symbol, t)
}

func (ob *Orderbook[T]) RemoveTrade(symbol string) {
	ob.Trades.Delete(symbol)
}

func (ob *Orderbook[T]) getTradeAmount(traderID string) int {
	i, ok := ob.TradeAmount.Get(traderID)
	if !ok {
		return -1
	}
	return i
}

// Increase increases trade count.
func (ob *Orderbook[T]) Increase(traderID string) {
	i := ob.getTradeAmount(traderID)
	if i == -1 {
		//log.Error(ob.Exchange, "cannot Increase Trade Count | val=-1")
		return
	}
	i++
	ob.TradeAmount.Set(traderID, i)
}

// Decrease decreases trade count.
func (ob *Orderbook[T]) Decrease(traderID string) {
	i := ob.getTradeAmount(traderID)
	if i == -1 {
		//log.Error(ob.Exchange, "cannot Decrease Trade Count | val=-1")
		return
	}
	i--
	ob.TradeAmount.Set(traderID, i)
}

// ValidateTrade checks if a trade exists or not.
func (ob *Orderbook[T]) ValidateTrade(binance *common.ExchangeTrade) (bool, error) {
	switch ob.Exchange {

	case constants.OKX:
		trade, ok := OKX_ORDERBOOK.Trades.Get(binance.Trade.Symbol)
		if !ok {
			return false, fmt.Errorf("ERROR [s:%s] | Unable to Find Trade | Trade != [newTrade]", binance.Trade.Symbol)
		}
		if trade.Trade.Info.TradeSettings.StopControl {
			return false, nil
		}
		return true, nil

	case constants.Binance:
		trade, ok := BSC_ORDERBOOK.Trades.Get(binance.Trade.Symbol)
		if !ok {
			return false, fmt.Errorf("ERROR [s:%s] | Unable to Find Trade | Trade != [newTrade]", binance.Trade.Symbol)
		}
		if trade.Trade.Info.TradeSettings.StopControl {
			return false, nil
		}
		return true, nil

	case constants.ByBit:
		trade, ok := BIT_ORDERBOOK.Trades.Get(binance.Trade.Symbol)
		if !ok {
			return false, fmt.Errorf("ERROR [s:%s] | Unable to Find Trade | Trade != [newTrade]", binance.Trade.Symbol)
		}
		if trade.Trade.Info.TradeSettings.StopControl {
			return false, nil
		}
		return true, nil

	default:
		return false, fmt.Errorf("ValidateTrade: ERROR Unsupported Exchange")
	}

}

// ValidateCounter checks if, from the user's settings, we reached max trades opened simultaneously.
func (ob *Orderbook[T]) ValidateCounter(traderID string) (bool, error) {
	counter, ok := ob.TradeAmount.Get(traderID)
	if !ok {
		return false, fmt.Errorf("ERROR Unable to get TradeAmount [%s]", traderID)
	}
	for _, trader := range config.CopyTradingCfg.Traders {
		for _, trd := range trader.Traders {
			if trd == traderID {
				if counter > trader.MaxOpenPositions {
					return false, nil
				} else {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// StartTradeWebsocketBroadcast initiates websocket connection & the Trade Manager (Orderbook).
func StartTradeWebsocketBroadcast(exchange constants.Exchange) {
	switch exchange {
	case constants.OKX:
		OKX_ORDERBOOK = NewOrderbook[okx.OKX]()
	case constants.Binance:
		BSC_ORDERBOOK = NewOrderbook[bsc.Binance]()
	case constants.ByBit:
		BIT_ORDERBOOK = NewOrderbook[bit.ByBit]()
	default:
		//log.Error("Socket_Broadcast", "Unable to Establish Connection – Unknown Exchange")
		return
	}

	handleTrades(exchange)

	// may not be needed
	var done chan struct{}
	<-done
}
