package nyx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	fhttp "github.com/bogdanfinn/fhttp"
)

type LeaderboardResponse struct {
	Code          any      `json:"code"`
	Message       any      `json:"message"`
	MessageDetail any      `json:"messageDetail"`
	Data          []Trader `json:"data"`
	Success       bool     `json:"success"`
}

type Trader struct {
	FutureUid      any     `json:"futureUid"`
	NickName       string  `json:"nickName"`
	UserPhotoUrl   string  `json:"userPhotoUrl"`
	Rank           int64   `json:"rank"`
	Pnl            float64 `json:"pnl"`
	Roi            float64 `json:"roi"`
	PositionShared bool    `json:"positionShared"`
	TwitterUrl     any     `json:"twitterUrl"`
	EncryptedUid   string  `json:"encryptedUid"`
	UpdateTime     int64   `json:"updateTime"`
	FolloweCount   int64   `json:"followerCount"`
	TwShared       string  `json:"-"`
	IsTwTrader     bool    `json:"isTwTrader"`
	OpenId         any     `json:"openId"`
}

func (app *Config) fetchTraders() ([]Trader, error) {

	binanceApi := "https://www.binance.com/bapi/futures/v3/public/future/leaderboard/getLeaderboardRank"

	postJson := make(map[string]interface{})
	var postData []byte
	postJson["tradeType"] = "PERPETUAL"
	postJson["statisticsType"] = "ROI"
	postJson["periodType"] = "WEEKLY"
	postJson["isShared"] = true
	postJson["isTrader"] = false
	postData, err := json.MarshalIndent(postJson, " ", "")
	if err != nil {
		return []Trader{}, err
	}
	req, err := fhttp.NewRequest("POST", binanceApi, bytes.NewReader(postData))
	if err != nil {
		return []Trader{}, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Origin", "https://www.binance.com")
	req.Header.Add("Referer", "https://www.binance.com/en/futures-activity/leaderboard/futures")
	req.Header.Add("lang", "en")
	req.Header.Add("clienttype", "web")
	req.Header.Add("content-type", "application/json")

	resp, err := app.TLSClient.Do(req)
	if err != nil {
		app.Logger.Error("Binance API request failed: " + err.Error())
		return []Trader{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Trader{}, err
	}

	var binanceApiResponse LeaderboardResponse
	err = json.Unmarshal(body, &binanceApiResponse)
	if err != nil {
		app.Logger.Error("Error decoding Binance API response " + err.Error())
		return []Trader{}, err
	}

	switch code := binanceApiResponse.Code.(type) {
	case float64:
		if code == 400 {
			app.Logger.Debug(fmt.Sprintf("failed connecting to binance API, status code: %.0f", code))
		}
	default:
		app.Logger.Debug(fmt.Sprintf("collected %d traders from binance API", len(binanceApiResponse.Data)))
	}

	return binanceApiResponse.Data, nil
}

func (app *Config) searchByNickname(nickname string) ([]Trader, error) {

	binanceApi := "https://www.binance.com/bapi/futures/v1/public/future/leaderboard/searchNickname"

	postJson := make(map[string]interface{})
	var postData []byte
	postJson["nickname"] = nickname

	postData, err := json.MarshalIndent(postJson, " ", "")
	if err != nil {
		return []Trader{}, err
	}
	req, err := fhttp.NewRequest("POST", binanceApi, bytes.NewReader(postData))
	if err != nil {
		return []Trader{}, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Origin", "https://www.binance.com")
	req.Header.Add("Referer", "https://www.binance.com/en/futures-activity/leaderboard/futures")
	req.Header.Add("lang", "en")
	req.Header.Add("clienttype", "web")
	req.Header.Add("content-type", "application/json")

	resp, err := app.TLSClient.Do(req)
	if err != nil {
		app.Logger.Error("Binance API request failed: " + err.Error())
		return []Trader{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Trader{}, err
	}

	var binanceApiResponse LeaderboardResponse
	err = json.Unmarshal(body, &binanceApiResponse)
	if err != nil {
		app.Logger.Error("Error decoding Binance API response " + err.Error())
		return []Trader{}, err
	}

	switch code := binanceApiResponse.Code.(type) {
	case float64:
		if code == 400 {
			app.Logger.Debug(fmt.Sprintf("failed connecting to binance API, status code: %.0f", code))
		}
	default:
		app.Logger.Debug(fmt.Sprintf("collected %d traders from binance API", len(binanceApiResponse.Data)))
	}

	return binanceApiResponse.Data, nil
}
