package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/UrbiJr/aqua-io/internal/core/crypto/constants"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"net/url"
	"strconv"
)

// getBalance returns initial balance of a ByBit account.
func (b *ByBit) getBalance() error {

	params := url.Values{
		"accountType": {"CONTRACT"},
		"coin":        {"USDT"},
	}

	resp, err := b.Do(http.MethodGet, betBalanceURL, params.Encode())
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.HandleERR(constants.ByBit, errors.HandleStatus(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response rest.Balance
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("GetBalance: %s", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if response.RetCode != 0 {
		return errors.HandleERR(constants.ByBit, fmt.Errorf("Getting Balance [%s]", response.RetMsg))
	} else {
		b.Info.Wallet.Value, err = strconv.ParseFloat(response.Result.List[0].Coin[0].Equity, 64)
		if err != nil {
			return err
		}
		return nil
	}
}

func (b *ByBit) setLeverage() error {

	var buf bytes.Buffer

	payload := map[string]interface{}{
		"category":     LinearType,
		"symbol":       b.Info.Coin.Symbol,
		"buyLeverage":  fmt.Sprint(b.Info.Leverage),
		"sellLeverage": fmt.Sprint(b.Info.Leverage),
	}

	/*
		if b.Trade.Side == "Buy" {
			payload["buyLeverage"] = fmt.Sprint(b.Trade.Leverage)
			payload["sellLeverage"] = fmt.Sprint(b.Trade.Leverage)
		} else {
			payload["sellLeverage"] = fmt.Sprintf("%f", b.Leverage.Value)
			payload["buyLeverage"] = fmt.Sprintf("%f", b.Leverage.Value)
		}
	*/

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := b.Do(http.MethodPost, setLeverageURL, buf.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Setting Leverage [%s] | [%s]", resp.Status, string(body))
	}

	var response rest.Leverage
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.RetCode == 0 {
		// success
		//b.Leverage.IsSet <- true
		return nil
	} else {
		return fmt.Errorf("ERROR Setting Leverage [%s]", response.RetMsg)

	}

}

// todo move to structs
func (b *ByBit) getCoinPrice() error {

	params := url.Values{
		"category": {"inverse"},
		"symbol":   {b.Trade.Info.Coin.Symbol},
	}

	resp, err := b.Do(http.MethodGet, getCoinPriceURL, params.Encode())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GetCoinPrice: ERROR Getting Price [%s]", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	list := response["result"].(map[string]interface{})["list"].([]interface{})
	markPrice := list[0].(map[string]interface{})["markPrice"].(string)
	b.Trade.Info.Coin.ExchangePrice, err = strconv.ParseFloat(markPrice, 64)
	if err != nil {
		return fmt.Errorf("getCoinPrice: ERR Casting to float64 markPrice")
	}

	return nil
}

func (b *ByBit) getPNL() error {
	//var buf bytes.Buffer
	var PnL string

	params := url.Values{
		"category": {LinearType},
		"symbol":   {b.Info.Coin.Symbol},
		"limit":    {"1"},
	}

	resp, err := b.Do(http.MethodGet, getPnLURL, params.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERR Getting PnL [%s]", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if resultList, ok := response["result"].(map[string]interface{})["list"].([]interface{}); ok {
		var valMap map[string]interface{}
		for _, val := range resultList {
			if valMap, ok = val.(map[string]interface{}); ok {
				if valMap["symbol"].(string) == b.Trade.Info.Coin.Symbol {
					var fPnL float64
					fPnL, err = strconv.ParseFloat(valMap["closedPnl"].(string), 64)
					if err != nil {
						return err
					}
					b.Trade.Info.PnL = int(fPnL)

					// was used to determine whether PnL was positive or negative
					/*
						if fPnL > 0 {
							b.Trade.Info.PnL = valMap["closedPnl"].(string)
							return nil
						} else {
							b.Trade.Info.PnL = valMap["closedPnl"].(string)
							return nil
						}
					*/
				}
			}
		}
	}

	if PnL == "" {
		PnL = "Unable to get PNL"
	}

	return nil
}
