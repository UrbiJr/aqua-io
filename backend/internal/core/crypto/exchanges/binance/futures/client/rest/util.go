package rest

import (
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/constants"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/exchanges/binance/futures/client/util"
	http "github.com/bogdanfinn/fhttp"
	"net/url"
	"strconv"
	"time"
)

const (
	positionURL     = "/fapi/v1/order"
	getBalanceURL   = "/fapi/v2/balance"
	getPnLURL       = "/fapi/v1/userTrades"
	getCoinPriceURL = "/fapi/v1/ticker/price"
	modifyMarginURL = "/fapi/v1/marginType"
	setLeverageURL  = "/fapi/v1/leverage"
)

func (b *Binance) Do(method string, params url.Values, pathURI string) (*http.Response, error) {
	var req http.Request

	uri, err := url.Parse(constants.BinancebaseRestAPI + pathURI)
	if err != nil {
		return nil, err
	}

	params.Add("signature", util.Sign(b.Misc.SecretAPI, params.Encode()))

	uri.RawQuery = params.Encode()
	req = http.Request{
		Method: method,
		URL:    uri,
	}

	return b.Misc.Client.Do(&req)
}

func (b *Binance) assert() {
	var accountName string
	for index, traders := range config.CopyTradingCfg.Traders {
		for _, trader := range traders.Traders {
			if trader == b.Info.TraderID {
				accountName = traders.Account
				b.Misc.Index = index
				b.Info.TradeMode = config.CopyTradingCfg.Traders[index].TradeMode
				b.Info.TradeSettings.MaxOpenPositions = config.CopyTradingCfg.Traders[index].MaxOpenPositions
				b.Info.TradeSettings.MaxCoinAllocation = config.CopyTradingCfg.Traders[index].MaxCoinAllocation
				b.Info.Leverage = config.CopyTradingCfg.Traders[index].Leverage
				b.Info.TradeSettings.BlockPositionAdds = config.CopyTradingCfg.Traders[index].BlockPositionAdds
				b.Info.TradeSettings.StopControl = config.CopyTradingCfg.Traders[index].StopControl
				b.Info.TradeSettings.MaxPriceDifferenceBetweenExchange = config.CopyTradingCfg.Traders[index].MaxPriceDifferenceBetweenExchange
				break
			}
		}
	}
	for _, binance := range config.CopyTradingCfg.Binance {
		if binance.AccountName == accountName {
			b.Misc.PublicAPI = binance.PublicAPI
			b.Misc.SecretAPI = binance.SecretAPI
			break
		}
	}
}

func now() string {
	n := time.Now().UnixMilli()
	return strconv.FormatInt(n, 64)
}
