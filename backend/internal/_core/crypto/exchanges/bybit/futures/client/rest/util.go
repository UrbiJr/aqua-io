package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/constants"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	positionURL     = "/v5/order/create"
	betBalanceURL   = "/v5/account/wallet-balance"
	getPnLURL       = "/v5/position/closed-pnl"
	getCoinPriceURL = "/v5/market/tickers"
	setLeverageURL  = "/v5/position/set-leverage"
	moveToBreakEven = "/v5/position/trading-stop"
)

func (b *ByBit) Do(method, pathURI string, payload string) (*http.Response, error) {
	var req http.Request

	uri, err := url.Parse(constants.ByBitbaseRestAPI + pathURI)
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost {
		req.Body = io.NopCloser(strings.NewReader(payload))
	} else if method == http.MethodGet {
		uri.RawQuery = payload
	} else {
		return nil, fmt.Errorf("ERR Unsupported HTTP Request Method [POST, GET]")
	}

	req.Method = method
	req.URL = uri
	req.Header = b.V5Sign(payload)

	return b.Misc.Client.Do(&req)
}

func (b *ByBit) assert() {
	var accountName string
	for index, trader := range config.CopyTradingCfg.Traders {
		for _, trd := range trader.Traders {
			if trd == b.Info.TraderID {
				accountName = trader.Account
				b.Misc.Index = index
				b.Info.TradeMode = config.CopyTradingCfg.Traders[b.Misc.Index].TradeMode
				b.Info.TradeSettings.MaxOpenPositions = config.CopyTradingCfg.Traders[b.Misc.Index].MaxOpenPositions
				b.Info.TradeSettings.MaxCoinAllocation = config.CopyTradingCfg.Traders[b.Misc.Index].MaxCoinAllocation
				b.Info.Leverage = config.CopyTradingCfg.Traders[b.Misc.Index].Leverage
				b.Info.TradeSettings.BlockPositionAdds = config.CopyTradingCfg.Traders[b.Misc.Index].BlockPositionAdds
				b.Info.TradeSettings.StopControl = config.CopyTradingCfg.Traders[b.Misc.Index].StopControl
				b.Info.TradeSettings.MaxPriceDifferenceBetweenExchange = config.CopyTradingCfg.Traders[b.Misc.Index].MaxPriceDifferenceBetweenExchange
				break
			}
		}
	}
	for _, bybit := range config.CopyTradingCfg.ByBit {
		if bybit.AccountName == accountName {
			b.Misc.PublicAPI = bybit.PublicAPI
			b.Misc.SecretAPI = bybit.SecretAPI
			break
		}
	}
}

func (b *ByBit) V5Sign(params string) http.Header {
	tNow := time.Now().UnixMilli()
	msg := fmt.Sprintf("%d%s%s%s", tNow, b.Misc.PublicAPI, "5000", params)
	mac := hmac.New(sha256.New, []byte(b.Misc.SecretAPI))
	mac.Write([]byte(msg))
	sign := hex.EncodeToString(mac.Sum(nil))

	return http.Header{
		"Content-Type":       {"application/json"},
		"X-BAPI-API-KEY":     {b.Misc.PublicAPI},
		"X-BAPI-TIMESTAMP":   {strconv.FormatInt(tNow, 10)},
		"X-BAPI-SIGN":        {sign},
		"X-BAPI-RECV-WINDOW": {"5000"},
	}
}
