package copy_io

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/UrbiJr/nyx/internal/user"
)

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

	var parsed map[string]interface{}
	var transactions []user.Transaction

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {

		if key == "result" {
			result := parsed["result"].(map[string]interface{})
			for key, _ := range result {
				if key == "list" && parsed["list"] != nil {
					list := parsed["list"].([]map[string]interface{})
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
							ProfileGroupID:  p.GroupID,
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
