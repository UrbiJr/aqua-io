package nyx

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

func (app *Config) getBybitPositions(bybitApiKey, bybitApiSecret string, testMode bool) []user.Position {
	var url string

	if testMode {
		url = "https://api-testnet.bybit.com/v5/position/list"
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
	str_to_sign := fmt.Sprintf("%d%s%s%s", unixMilli, bybitApiKey, "20000", "")

	hm := hmac.New(sha256.New, []byte(bybitApiSecret))
	hm.Write([]byte(str_to_sign))
	HMAC := hex.EncodeToString(hm.Sum(nil))

	req.Header.Add("X-BAPI-API-KEY", bybitApiKey)
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
	var positions []user.Position

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &parsed)

	for key, _ := range parsed {

		if key == "result" {
			result := parsed["result"].(map[string]interface{})
			for key, _ := range result {
				if key == "list" {
					list := parsed["list"].([]map[string]interface{})
					for _, x := range list {
						leverage, err := strconv.ParseInt(x["leverage"].(string), 10, 64)
						if err != nil {
							continue
						}
						markPrice, err := strconv.ParseFloat(x["markPrice"].(string), 64)
						if err != nil {
							continue
						}
						unrealisedPnl, err := strconv.ParseFloat(x["unrealisedPnl"].(string), 64)
						if err != nil {
							continue
						}
						updateTimestamp, err := strconv.ParseInt(x["updatedTime"].(string), 10, 64)
						if err != nil {
							continue
						}
						positions = append(positions, user.Position{
							Symbol:          x["symbol"].(string),
							Leverage:        leverage,
							MarkPrice:       markPrice,
							Pnl:             unrealisedPnl,
							UpdateTimestamp: updateTimestamp,
						})
					}
				}
			}
		}
	}

	app.Logger.Debug(fmt.Sprintf("Fetched %d positions from bybit", len(positions)))

	return positions
}
