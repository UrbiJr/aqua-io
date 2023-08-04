package rest

import (
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"net/url"
	"strconv"
	"strings"
)

func (b *Binance) getPNL() (int, error) {

	params := url.Values{
		"symbol":    {b.Info.Coin.Symbol},
		"limit":     {"1"},
		"timestamp": {now()},
	}

	resp, err := b.Do(http.MethodGet, params, getPnLURL)
	if err != nil {
		return -1, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != 200 {
		return -1, errors.HandleERR("constants.Binance", fmt.Errorf("fr je sais pas quelle erreur c"))
	}

	defer resp.Body.Close()

	var pnl rest.PNL
	err = json.Unmarshal(body, &pnl)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(pnl[0].RealizedPnl)
}

func (b *Binance) getCoinPrice() (float64, error) {
	payload := url.Values{
		"symbol":    {b.Info.Coin.Symbol},
		"timestamp": {now()},
	}

	resp, err := b.Do(http.MethodGet, payload, getCoinPriceURL)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != 200 {
		return -1, fmt.Errorf("")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var ticker rest.Ticker
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		return -1, err
	}

	for _, tck := range ticker {
		if tck.Symbol == b.Info.Coin.Symbol {
			var price float64
			price, err = strconv.ParseFloat(tck.Price, 64)
			if err != nil {
				return -1, err
			}
			return price, nil
		}
	}

	return -1, fmt.Errorf("ERROR Fetching Binance Price ")
}

func (b *Binance) changeMargin() error {
	payload := url.Values{
		"symbol":     {b.Info.Coin.Symbol},
		"marginType": {strings.ToUpper(b.Info.TradeMode)},
		"timestamp":  {now()},
	}

	resp, err := b.Do(http.MethodGet, payload, modifyMarginURL)
	if err != nil {
		return err
	}

	log.Info("change margin resp", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Info("body of change margin", string(body))

	return nil
}

func (b *Binance) getBalance() error {
	payload := url.Values{
		"timestamp": {now()},
	}

	resp, err := b.Do(http.MethodGet, payload, getBalanceURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var balance rest.Balance
	err = json.Unmarshal(body, &balance)
	if err != nil {
		return err
	}

	for _, asset := range balance {
		if asset.Asset == "USDT" {
			b.Trade.Info.Wallet.Value, err = strconv.ParseFloat(asset.Balance, 64)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("Unable to Fetch Balance")
}

func (b *Binance) setLeverage() error {
	payload := url.Values{
		"symbol":    {b.Info.Coin.Symbol},
		"leverage":  {strconv.FormatFloat(b.Info.Leverage, 'f', 0, 64)},
		"timestamp": {now()},
	}

	resp, err := b.Do(http.MethodPost, payload, setLeverageURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Info("resp body lev set", string(body))

	if resp.StatusCode != 200 {
		return errors.HandleERR("constants.Binance", fmt.Errorf("unknown [%s]", resp.Status))
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return errors.HandleERR("constants.Binance", err)
	}

	if response["leverage"] != int(b.Info.Leverage) {
		return errors.HandleERR("constants.Binance", fmt.Errorf("setting leverage: wanted [%d], got [%s]", int(b.Info.Leverage), response["leverage"]))
	} else {
		log.Success("Successfully Set Leverage")
		return nil
	}
}
