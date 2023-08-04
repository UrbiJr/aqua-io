package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"net/url"
	"strings"
	"time"
)

const (
	openPositionURL          = "/api/v5/trade/order"
	closePositionURL         = "/api/v5/trade/close-position"
	getBalanceURL            = "/api/v5/account/balance"
	getPnLURL                = "/api/v5/account/positions-history"
	getCoinPriceURL          = "/api/v5/market/ticker"
	updatePositionMargin     = "/api/v5/account/position/margin-balance"
	setLeverageURL           = "/api/v5/account/set-leverage"
	cancelOrderURL           = "/api/v5/trade/cancel-order"
	getInstrumentLeverageURL = "/api/v5/public/instruments"
)

func (o *OKX) Do(method, pathURI, params string, isPrivate bool) (*http.Response, error) {

	var req http.Request
	uri, err := url.Parse(constants.OKXbaseRestAPI + pathURI)
	if err != nil {
		return nil, err
	}

	req = http.Request{
		Method: method,
		URL:    uri,
	}

	if method == http.MethodPost {
		req.Body = io.NopCloser(strings.NewReader(params))
		req.Header = o.sign(params, method, pathURI)
	} else if method == http.MethodGet {
		uri.RawQuery = params
		if isPrivate {
			req.Header = o.sign(params, method, pathURI)
		} else {
			req.Header = http.Header{"User-Agent": {"Arcana/GoLang"}}
		}
	}
	return o.Trade.Misc.Client.Do(&req)
}

func (o *OKX) assert() {
	var accountName string
	for index, trader := range config.CopyTradingCfg.Traders {
		for _, trd := range trader.Traders {
			if trd == o.Trade.Info.TraderID {
				accountName = trader.Account
				o.Trade.Misc.Index = index
				o.Trade.Info.TradeMode = config.CopyTradingCfg.Traders[o.Trade.Misc.Index].TradeMode
				o.Trade.Info.TradeSettings.MaxOpenPositions = config.CopyTradingCfg.Traders[o.Trade.Misc.Index].MaxOpenPositions
				o.Trade.Info.TradeSettings.MaxCoinAllocation = config.CopyTradingCfg.Traders[o.Trade.Misc.Index].MaxCoinAllocation
				o.Trade.Info.Leverage = config.CopyTradingCfg.Traders[o.Trade.Misc.Index].Leverage
				o.Trade.Info.TradeSettings.BlockPositionAdds = config.CopyTradingCfg.Traders[o.Trade.Misc.Index].BlockPositionAdds
				o.Trade.Info.TradeSettings.StopControl = config.CopyTradingCfg.Traders[o.Trade.Misc.Index].StopControl
				break
			}
		}
	}

	for _, okx := range config.CopyTradingCfg.OKX {
		if okx.AccountName == accountName {
			o.Trade.Misc.PublicAPI = okx.PublicAPI
			o.Trade.Misc.SecretAPI = okx.SecretAPI
			o.Trade.Misc.Passphrase = okx.Passphrase
			break
		}
	}

	if o.Trade.Info.Side == "long" {
		o.Trade.Info.Side = constants.OKX_BUY
	} else {
		o.Trade.Info.Side = constants.OKX_SELL
	}
}

func (o *OKX) calcPositionChange() {}

func (o *OKX) formatCoin() {
	if w := o.Trade.Info.Coin.Symbol[len(o.Trade.Info.Coin.Symbol)-4:]; w == "USDT" {
		splitSymbol := strings.Split(o.Trade.Info.Coin.Symbol, "USDT")
		o.Trade.Info.Coin.Symbol = strings.Join(splitSymbol, "-USDT")
	} else {
		log.Error(fmt.Sprintf("ERROR Unsupported Coin? [%s] – Request SUPPORT", o.Trade.Info.Coin.Symbol))
	}
}

func (o *OKX) sign(params, method, pathURI string) http.Header {
	format := "2006-01-02T15:04:05.999Z07:00"
	tNow := time.Now().UTC().Format(format)
	msg := []byte(fmt.Sprintf("%s%s%s%s", tNow, method, pathURI, params))
	mac := hmac.New(sha256.New, []byte(o.Trade.Misc.SecretAPI))
	mac.Write(msg)
	return http.Header{
		"Content-Type":         {"application/json"},
		"OK-ACCESS-KEY":        {o.Trade.Misc.PublicAPI},
		"OK-ACCESS-PASSPHRASE": {o.Trade.Misc.Passphrase},
		"OK-ACCESS-SIGN":       {base64.StdEncoding.EncodeToString(mac.Sum(nil))},
		"OK-ACCESS-TIMESTAMP":  {fmt.Sprint(tNow)},
	}
}
