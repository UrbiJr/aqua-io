package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/UrbiJr/aqua-io/internal/core/crypto/exchanges/okx/futures/models/rest"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"net/url"
	"strconv"
)

func (o *OKX) getBalance() error {

	resp, err := o.Do(http.MethodGet, getBalanceURL, "", true)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Fetching Balance [%s]", resp.Status)
	}

	var balance rest.Balance
	err = json.NewDecoder(resp.Body).Decode(&balance)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	for _, val := range balance.Data {
		for _, detail := range val.Details {
			if detail.Ccy == "USDT" {
				o.Trade.Info.Wallet.Value, err = strconv.ParseFloat(detail.AvailBal, 64)
				return err
			}
		}
	}

	return nil
}

func (o *OKX) setLeverage() error {
	var buf bytes.Buffer

	payload := map[string]string{
		"instId":  o.Trade.Info.Coin.Symbol,
		"lever":   fmt.Sprintf("%d", int(o.Trade.Info.Leverage)),
		"mgnMode": o.Trade.Info.TradeMode,
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := o.Do(http.MethodPost, setLeverageURL, buf.String(), true)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Setting Leverage [%s]", resp.Status)
	}

	var leverage rest.Leverage
	err = json.NewDecoder(resp.Body).Decode(&leverage)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if leverage.Code == "0" && leverage.Msg == "" {
		return nil
	} else {
		return fmt.Errorf("ERROR Set Leverage [%s | %s]", leverage.Msg, leverage.Code)
	}
}

func (o *OKX) GetPnL() error {

	resp, err := o.Do(http.MethodGet, getPnLURL, "", true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Getting PnL [%s]", resp.Status)
	}

	var positionsHistory rest.PositionsHistory
	err = json.Unmarshal(body, &positionsHistory)
	if err != nil {
		return err
	}

	o.Trade.Info.PnL, err = strconv.Atoi(positionsHistory.Data[0].Pnl)
	if err != nil {
		return err
	}

	return nil
}

func (o *OKX) getCoinPrice() error {

	params := url.Values{
		"instId": {o.Trade.Info.Coin.Symbol},
	}

	resp, err := o.Do(http.MethodGet, getCoinPriceURL, params.Encode(), false)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Getting Coin Price [%s | %s]", o.Trade.Info.Coin.ExchangePrice, resp.Status)
	}

	var ticker rest.Ticker
	err = json.NewDecoder(resp.Body).Decode(&ticker)
	if err != nil {
		return err
	}

	if ticker.Msg == "" && ticker.Code == "0" {
		o.Trade.Info.Coin.ExchangePrice, err = strconv.ParseFloat(ticker.Data[0].Last, 64)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("ERROR Getting Price [%s | %s]", ticker.Msg, ticker.Code)
	}
}

func (o *OKX) getInstrumentsLeverage() error {

	o.formatCoin()

	params := url.Values{
		"instType":   {"FUTURES"},
		"instFamily": {o.Trade.Info.Coin.Symbol},
	}

	resp, err := o.Do(http.MethodGet, getInstrumentLeverageURL, params.Encode(), false)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			return fmt.Errorf("[OKX] Unsupported Coin on Futures [%s]", o.Trade.Info.Coin.Symbol)
		}
		return fmt.Errorf("ERROR Getting Instrument Leverage [%s]", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var instruments rest.Instruments
	err = json.Unmarshal(body, &instruments)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if instruments.Msg == "" && instruments.Code == "0" {
		for _, instrument := range instruments.Data {
			if instrument.Alias == "this_week" {
				if instrument.Lever < fmt.Sprint(o.Trade.Info.Leverage) {
					var lev float64
					lev, err = strconv.ParseFloat(instruments.Data[0].Lever, 64)
					if err != nil {
						return err
					}
					o.Trade.Info.Leverage = lev
					break
				}
			}
		}
		return nil
	} else {
		return fmt.Errorf("ERROR Getting Price [%s | %s]", instruments.Msg, instruments.Code)
	}

}

func (o *OKX) cancelOrder() error {

	var buf bytes.Buffer

	var payload = map[string]interface{}{
		"ordId":  o.Trade.Info.TradeID,
		"instId": o.Trade.Info.Coin.Symbol,
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := o.Do(http.MethodPost, cancelOrderURL, buf.String(), true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Cancelling Order [%s | %s]", o.Trade.Info.TradeID, resp.Status)
	}

	var cancel rest.Cancel
	err = json.NewDecoder(resp.Body).Decode(&cancel)
	if err != nil {
		return err
	}

	if cancel.Code == "0" && cancel.Msg == "" {
		log.Success(constants.OKX, "Position Cancelled [Not Filled by Market]")
	} else {
		return fmt.Errorf("ERROR Cancelling Position [%s]", cancel.Msg)
	}

	return nil
}

func parseTickers() {}
