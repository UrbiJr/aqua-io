package copy_io

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/utils"
)

func (app *Config) createOrder(p *user.Profile, symbol, side, orderType string, amount, price, TP, SL float64) (*user.OpenedPosition, error) {
	var url string

	if utils.Contains(p.BlacklistCoins, symbol) {
		return nil, fmt.Errorf("cannot create order: symbol %s is in user blacklisted coins", symbol)
	}

	if p.TestMode {
		url = "https://api-testnet.bybit.com/v5/order/create"
	}
	method := "POST"

	var timeInForce string
	if orderType == "Market" {
		timeInForce = "IOC"
	} else {
		timeInForce = "GTC"
	}

	var stopLossStr, takeProfitStr string
	if SL == price {
		stopLossStr = ""
	} else {
		stopLossStr = fmt.Sprintf("%f", SL)
	}
	if TP == price {
		takeProfitStr = ""
	} else {
		takeProfitStr = fmt.Sprintf("%f", TP)
	}

	postData := fmt.Sprintf(`{
		"category": "spot",
		"symbol": "%s",
		"side": "%s",
		"orderType": "%s",
		"qty": "%f",
		"price": "%s",
		"timeInForce": "%s",
		"isLeverage": "%d",
		"orderFilter": "Order",
		"takeProfit": "%s",
		"stopLoss": "%s"
	}`, symbol, side, orderType, math.Abs(amount), fmt.Sprintf("%.2f", price), timeInForce, p.Leverage, takeProfitStr, stopLossStr)

	req, err := http.NewRequest(method, url, strings.NewReader(postData))

	if err != nil {
		app.Logger.Error(err)
		return nil, err
	}

	// create a time variable
	now := time.Now()
	// convert to unix time in milliseconds
	unixMilli := now.UnixMilli()

	// generate hmac for X-BAPI-SIGN
	str_to_sign := fmt.Sprintf("%d%s%s%s", unixMilli, p.BybitApiKey, "5000", postData)

	hm := hmac.New(sha256.New, []byte(p.BybitApiSecret))
	hm.Write([]byte(str_to_sign))
	HMAC := hex.EncodeToString(hm.Sum(nil))

	req.Header.Add("X-BAPI-API-KEY", p.BybitApiKey)
	req.Header.Add("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", unixMilli))
	req.Header.Add("X-BAPI-RECV-WINDOW", "5000")
	req.Header.Add("X-BAPI-SIGN", HMAC)
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Client.Do(req)
	if err != nil {
		app.Logger.Error(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error(err)
		return nil, err
	}

	var parsed map[string]any

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {
		if key == "retMsg" {
			if parsed["retMsg"].(string) == "OK" {
				result := parsed["result"].(map[string]any)
				for key, _ = range result {
					if key == "orderId" {
						app.Logger.Debug(fmt.Sprintf("successfully created bybit %s order with id %s", side, result[key].(string)))
						return &user.OpenedPosition{OrderID: result[key].(string), Symbol: symbol, ProfileID: p.ID}, nil
					}
				}
			} else {
				if strings.Contains(parsed["retMsg"].(string), "Timestamp for this request is outside of the recvWindow.") {
					// send notification to adjust system time
					return nil, errors.New("Timestamp not synchronized: please sync your system time and try again")
				}
				return nil, errors.New(parsed["retMsg"].(string))
			}
		}
	}

	return nil, fmt.Errorf("create order failed: order id not found")
}

func (app *Config) cancelOrder(p *user.Profile, category, orderID, symbol string) error {
	var url string

	if p.TestMode {
		url = "https://api-testnet.bybit.com/v5/order/cancel"
	}
	method := "POST"

	postData := fmt.Sprintf(`{
		"category": "%s",
		"symbol": "%s",
		"orderId": "%s",
	}`, category, symbol, orderID)

	req, err := http.NewRequest(method, url, strings.NewReader(postData))

	if err != nil {
		app.Logger.Error(err)
		return err
	}

	// create a time variable
	now := time.Now()
	// convert to unix time in milliseconds
	unixMilli := now.UnixMilli()

	// generate hmac for X-BAPI-SIGN
	str_to_sign := fmt.Sprintf("%d%s%s%s", unixMilli, p.BybitApiKey, "5000", postData)

	hm := hmac.New(sha256.New, []byte(p.BybitApiSecret))
	hm.Write([]byte(str_to_sign))
	HMAC := hex.EncodeToString(hm.Sum(nil))

	req.Header.Add("X-BAPI-API-KEY", p.BybitApiKey)
	req.Header.Add("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", unixMilli))
	req.Header.Add("X-BAPI-RECV-WINDOW", "5000")
	req.Header.Add("X-BAPI-SIGN", HMAC)
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Client.Do(req)
	if err != nil {
		app.Logger.Error(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error(err)
		return err
	}

	var parsed map[string]any

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {
		if key == "retMsg" {
			if parsed["retMsg"].(string) == "OK" {
				result := parsed["result"].(map[string]any)
				for key, _ = range result {
					if key == "orderId" {
						app.Logger.Debug(fmt.Sprintf("successfully cancelled %s order with ID %s", symbol, orderID))
						return nil
					}
				}
			} else {
				if strings.Contains(parsed["retMsg"].(string), "Timestamp for this request is outside of the recvWindow.") {
					// send notification to adjust system time
					return errors.New("Timestamp not synchronized: please sync your system time and try again")
				}
				return errors.New(parsed["retMsg"].(string))
			}
		}
	}

	return fmt.Errorf("cancel order failed")
}

func (app *Config) getBybitTransactions(p user.Profile) []user.Transaction {
	var url string

	if p.TestMode {
		url = "https://api-testnet.bybit.com/v5/account/transaction-log"
	}
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		app.Logger.Error(err)
		return nil
	}

	// create a time variable
	now := time.Now()
	// convert to unix time in milliseconds
	unixMilli := now.UnixMilli()

	// generate hmac for X-BAPI-SIGN
	str_to_sign := fmt.Sprintf("%d%s%s%s", unixMilli, p.BybitApiKey, "20000", "")

	hm := hmac.New(sha256.New, []byte(p.BybitApiSecret))
	hm.Write([]byte(str_to_sign))
	HMAC := hex.EncodeToString(hm.Sum(nil))

	req.Header.Add("X-BAPI-API-KEY", p.BybitApiKey)
	req.Header.Add("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", unixMilli))
	req.Header.Add("X-BAPI-RECV-WINDOW", "20000")
	req.Header.Add("X-BAPI-SIGN", HMAC)

	res, err := app.Client.Do(req)
	if err != nil {
		app.Logger.Error(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error(err)
		return nil
	}

	var parsed map[string]any
	var transactions []user.Transaction

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {

		if key == "result" {
			result := parsed["result"].(map[string]any)
			for key, _ := range result {
				if key == "list" && parsed["list"] != nil {
					list := parsed["list"].([]map[string]any)
					for _, x := range list {
						tradePrice, err := strconv.ParseFloat(x["tradePrice"].(string), 64)
						if err != nil {
							continue
						}
						qty, err := strconv.ParseFloat(x["qty"].(string), 64)
						if err != nil {
							continue
						}
						size, err := strconv.ParseFloat(x["size"].(string), 64)
						if err != nil {
							continue
						}
						funding, err := strconv.ParseFloat(x["funding"].(string), 64)
						if err != nil {
							continue
						}
						transactionTime, err := strconv.ParseInt(x["transactionTime"].(string), 10, 64)
						if err != nil {
							continue
						}
						transactions = append(transactions, user.Transaction{
							ProfileID:       p.ID,
							OrderID:         x["orderId"].(string),
							TradeID:         x["tradeId"].(string),
							Symbol:          x["symbol"].(string),
							Funding:         funding,
							Currency:        x["currency"].(string),
							TradePrice:      tradePrice,
							Qty:             qty,
							Size:            size,
							Side:            x["side"].(string),
							TransactionTime: transactionTime,
						})
					}
				}
			}
		}
	}

	app.Logger.Debug(fmt.Sprintf("Fetched %d transactions from bybit", len(transactions)))

	return transactions
}

func (app *Config) fetchOrderHistory(category, orderFilter string, p user.Profile) ([]user.Order, error) {
	var url string

	if orderFilter == "" {
		orderFilter = "Order"
	}
	if p.TestMode {
		url = "https://api-testnet.bybit.com/v5/order/history?category=" + category + "&orderFilter=" + orderFilter
	}
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		app.Logger.Error(err)
		return nil, err
	}

	// create a time variable
	now := time.Now()
	// convert to unix time in milliseconds
	unixMilli := now.UnixMilli()

	// generate hmac for X-BAPI-SIGN
	str_to_sign := fmt.Sprintf("%d%s%s%s", unixMilli, p.BybitApiKey, "20000", "category="+category+"&orderFilter="+orderFilter)

	hm := hmac.New(sha256.New, []byte(p.BybitApiSecret))
	hm.Write([]byte(str_to_sign))
	HMAC := hex.EncodeToString(hm.Sum(nil))

	req.Header.Add("X-BAPI-API-KEY", p.BybitApiKey)
	req.Header.Add("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", unixMilli))
	req.Header.Add("X-BAPI-RECV-WINDOW", "20000")
	req.Header.Add("X-BAPI-SIGN", HMAC)

	res, err := app.Client.Do(req)
	if err != nil {
		app.Logger.Error(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error(err)
		return nil, err
	}

	if strings.Contains(string(body), "Timestamp for this request is outside of the recvWindow.") {
		// send notification to adjust system time
		return nil, errors.New("Timestamp not synchronized: please sync your system time and try again")
	}

	var parsed map[string]any
	var orders []user.Order

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {

		if key == "result" {
			result := parsed["result"].(map[string]any)
			for key := range result {
				if key == "list" && result["list"] != nil {
					list := result["list"].([]any)
					for _, x := range list {
						switch o := x.(type) {
						case map[string]any:
							price, err := strconv.ParseFloat(o["price"].(string), 64)
							if err != nil {
								app.Logger.Error(err)
							}
							avgPrice, err := strconv.ParseFloat(o["avgPrice"].(string), 64)
							if err != nil {
								app.Logger.Error(err)
							}
							triggerPrice, err := strconv.ParseFloat(o["triggerPrice"].(string), 64)
							if err != nil {
								app.Logger.Error(err)
							}
							takeProfit, err := strconv.ParseFloat(o["takeProfit"].(string), 64)
							if err != nil {
								app.Logger.Error(err)
							}
							stopLoss, err := strconv.ParseFloat(o["stopLoss"].(string), 64)
							if err != nil {
								app.Logger.Error(err)
							}
							qty, err := strconv.ParseFloat(o["qty"].(string), 64)
							if err != nil {
								app.Logger.Error(err)
							}
							filledQty, err := strconv.ParseFloat(o["cumExecQty"].(string), 64)
							if err != nil {
								app.Logger.Error(err)
							}
							isLeverage, err := strconv.ParseInt(o["isLeverage"].(string), 10, 64)
							if err != nil {
								app.Logger.Error(err)
							}
							createdTime, err := strconv.ParseInt(o["createdTime"].(string), 10, 64)
							if err != nil {
								app.Logger.Error(err)
							}
							orders = append(orders, user.Order{
								ProfileID:        p.ID,
								Symbol:           o["symbol"].(string),
								OrderID:          o["orderId"].(string),
								OrderLinkID:      o["orderLinkId"].(string),
								OrderStatus:      o["orderStatus"].(string),
								OrderType:        o["orderType"].(string),
								AvgFilledPrice:   avgPrice,
								Price:            price,
								TakeProfit:       takeProfit,
								StopLoss:         stopLoss,
								TriggerPrice:     triggerPrice,
								TriggerDirection: o["triggerDirection"].(float64),
								CreatedTime:      createdTime,
								FilledQty:        filledQty,
								Qty:              qty,
								Side:             o["side"].(string),
								IsLeverage:       isLeverage,
							})
						}
					}
				}
			}
		}
	}

	app.Logger.Debug(fmt.Sprintf("Fetched %d orders from bybit", len(orders)))

	return orders, nil
}

func (app *Config) fetchOpenOrders(category string, p user.Profile) []user.Order {
	var url string

	if p.TestMode {
		url = "https://api-testnet.bybit.com/v5/order/realtime?category=" + category
	}
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		app.Logger.Error(err)
		return nil
	}

	// create a time variable
	now := time.Now()
	// convert to unix time in milliseconds
	unixMilli := now.UnixMilli()

	// generate hmac for X-BAPI-SIGN
	str_to_sign := fmt.Sprintf("%d%s%s%s", unixMilli, p.BybitApiKey, "5000", "category="+category)

	hm := hmac.New(sha256.New, []byte(p.BybitApiSecret))
	hm.Write([]byte(str_to_sign))
	HMAC := hex.EncodeToString(hm.Sum(nil))

	req.Header.Add("X-BAPI-API-KEY", p.BybitApiKey)
	req.Header.Add("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", unixMilli))
	req.Header.Add("X-BAPI-RECV-WINDOW", "5000")
	req.Header.Add("X-BAPI-SIGN", HMAC)

	res, err := app.Client.Do(req)
	if err != nil {
		app.Logger.Error(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error(err)
		return nil
	}

	var parsed map[string]any
	var orders []user.Order

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {

		if key == "result" {
			result := parsed["result"].(map[string]any)
			for key := range result {
				if key == "list" && result["list"] != nil {
					list := result["list"].([]any)
					for _, x := range list {
						switch o := x.(type) {
						case map[string]any:
							price, err := strconv.ParseFloat(o["price"].(string), 64)
							if err != nil {
								continue
							}
							triggerPrice, err := strconv.ParseFloat(o["triggerPrice"].(string), 64)
							if err != nil {
								continue
							}
							qty, err := strconv.ParseFloat(o["qty"].(string), 64)
							if err != nil {
								continue
							}
							isLeverage, err := strconv.ParseInt(o["isLeverage"].(string), 10, 64)
							if err != nil {
								continue
							}
							createdTime, err := strconv.ParseInt(o["createdTime"].(string), 10, 64)
							if err != nil {
								continue
							}
							orders = append(orders, user.Order{
								ProfileID:    p.ID,
								Symbol:       o["symbol"].(string),
								OrderID:      o["orderId"].(string),
								OrderLinkID:  o["orderLinkId"].(string),
								OrderStatus:  o["orderStatus"].(string),
								OrderType:    o["orderType"].(string),
								Price:        price,
								TriggerPrice: triggerPrice,
								CreatedTime:  createdTime,
								Qty:          qty,
								Side:         o["side"].(string),
								IsLeverage:   isLeverage,
							})
						}
					}
				}
			}
		}
	}

	app.Logger.Debug(fmt.Sprintf("Fetched %d open orders from bybit", len(orders)))

	return orders
}

func (app *Config) getPositionInfo(category, symbol, orderID string, p user.Profile) []user.PositionInfo {
	var url string

	if p.TestMode {
		url = "https://api-testnet.bybit.com/v5/position/list?category=" + category + "&symbol=" + symbol
	}
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		app.Logger.Error(err)
		return nil
	}

	// create a time variable
	now := time.Now()
	// convert to unix time in milliseconds
	unixMilli := now.UnixMilli()

	// generate hmac for X-BAPI-SIGN
	str_to_sign := fmt.Sprintf("%d%s%s%s", unixMilli, p.BybitApiKey, "5000", "category="+category+"&symbol="+symbol)

	hm := hmac.New(sha256.New, []byte(p.BybitApiSecret))
	hm.Write([]byte(str_to_sign))
	HMAC := hex.EncodeToString(hm.Sum(nil))

	req.Header.Add("X-BAPI-API-KEY", p.BybitApiKey)
	req.Header.Add("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", unixMilli))
	req.Header.Add("X-BAPI-RECV-WINDOW", "5000")
	req.Header.Add("X-BAPI-SIGN", HMAC)

	res, err := app.Client.Do(req)
	if err != nil {
		app.Logger.Error(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		app.Logger.Error(err)
		return nil
	}

	var parsed map[string]any
	var positionInfoArr []user.PositionInfo

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {

		if key == "result" {
			result := parsed["result"].(map[string]any)
			for key := range result {
				if key == "list" && result["list"] != nil {
					list := result["list"].([]any)
					for _, x := range list {
						switch o := x.(type) {
						case map[string]any:
							leverage, err := strconv.ParseInt(o["leverage"].(string), 10, 64)
							if err != nil {
								continue
							}
							avgPrice, err := strconv.ParseFloat(o["avgPrice"].(string), 64)
							if err != nil {
								continue
							}
							liqPrice, err := strconv.ParseFloat(o["liqPrice"].(string), 64)
							if err != nil {
								continue
							}
							positionValue, err := strconv.ParseFloat(o["positionValue"].(string), 64)
							if err != nil {
								continue
							}
							unrealisedPnl, err := strconv.ParseFloat(o["unrealisedPnl"].(string), 64)
							if err != nil {
								continue
							}
							cumRealisedPnl, err := strconv.ParseFloat(o["cumRealisedPnl"].(string), 64)
							if err != nil {
								continue
							}
							markPrice, err := strconv.ParseFloat(o["markPrice"].(string), 64)
							if err != nil {
								continue
							}
							createdTime, err := strconv.ParseInt(o["createdTime"].(string), 10, 64)
							if err != nil {
								continue
							}
							updatedTime, err := strconv.ParseInt(o["updatedTime"].(string), 10, 64)
							if err != nil {
								continue
							}
							positionInfoArr = append(positionInfoArr, user.PositionInfo{
								OrderID:        orderID,
								PositionIdx:    o["positionIdx"].(float64),
								Symbol:         o["symbol"].(string),
								Leverage:       leverage,
								AvgPrice:       avgPrice,
								LiqPrice:       liqPrice,
								TakeProfit:     o["takeProfit"],
								StopLoss:       o["stopLoss"],
								PositionValue:  positionValue,
								UnrealisedPnl:  unrealisedPnl,
								CumRealisedPnl: cumRealisedPnl,
								MarkPrice:      markPrice,
								CreatedTime:    createdTime,
								UpdatedTime:    updatedTime,
								Side:           o["side"].(string),
								PositionStatus: o["positionStatus"].(string),
							})
						}
					}
				}
			}
		}
	}

	app.Logger.Debug(fmt.Sprintf("Fetched %d positions from bybit", len(positionInfoArr)))

	return positionInfoArr
}
