package copy_io

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/UrbiJr/aqua-io/internal/user"
	fhttp "github.com/bogdanfinn/fhttp"
)

type GetOtherPositionResponse struct {
	Code          any            `json:"code"`
	Message       any            `json:"message"`
	MessageDetail any            `json:"messageDetail"`
	Data          map[string]any `json:"data"`
	Success       bool           `json:"success"`
}

type LeaderboardResponse struct {
	Code          any           `json:"code"`
	Message       any           `json:"message"`
	MessageDetail any           `json:"messageDetail"`
	Data          []user.Trader `json:"data"`
	Success       bool          `json:"success"`
}

type TraderResponse struct {
	Code          any         `json:"code"`
	Message       any         `json:"message"`
	MessageDetail any         `json:"messageDetail"`
	Data          user.Trader `json:"data"`
	Success       bool        `json:"success"`
}

func (app *Config) fetchTraderByUid(encryptedUid string) (*user.Trader, error) {

	binanceApi := "https://www.binance.com/bapi/futures/v2/public/future/leaderboard/getOtherLeaderboardBaseInfo"

	postJson := make(map[string]any)
	var postData []byte
	postJson["encryptedUid"] = encryptedUid
	postData, err := json.MarshalIndent(postJson, " ", "")
	if err != nil {
		return nil, err
	}
	req, err := fhttp.NewRequest("POST", binanceApi, bytes.NewReader(postData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Origin", "https://www.binance.com")
	req.Header.Add("Referer", "https://www.binance.com/en/futures-activity/leaderboard/user/um?encryptedUid="+encryptedUid)
	req.Header.Add("lang", "en")
	req.Header.Add("clienttype", "web")
	req.Header.Add("content-type", "application/json")

	resp, err := app.TLSClient.Do(req)
	if err != nil {
		app.Logger.Error("Binance API request failed: " + err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var binanceApiResponse TraderResponse
	err = json.Unmarshal(body, &binanceApiResponse)
	if err != nil {
		app.Logger.Error("Error decoding Binance API response " + err.Error())
		return nil, err
	}

	switch code := binanceApiResponse.Code.(type) {
	case float64:
		if code == 400 {
			err := errors.New(fmt.Sprintf("failed connecting to binance API, status code: %.0f", code))
			app.Logger.Error(err)
			return nil, err
		}
	default:
		app.Logger.Debug(fmt.Sprintf("fetched trader with ID from binance API", encryptedUid))
	}

	binanceApiResponse.Data.EncryptedUid = encryptedUid
	return &binanceApiResponse.Data, nil
}

func (app *Config) fetchTraders(statisticsType, periodType string) ([]user.Trader, error) {

	if periodType == "TOTAL" {
		periodType = "ALL"
	}
	binanceApi := "https://www.binance.com/bapi/futures/v3/public/future/leaderboard/getLeaderboardRank"

	postJson := make(map[string]any)
	var postData []byte
	postJson["tradeType"] = "PERPETUAL"
	postJson["statisticsType"] = statisticsType
	postJson["periodType"] = periodType
	postJson["isShared"] = true
	postJson["isTrader"] = false
	postData, err := json.MarshalIndent(postJson, " ", "")
	if err != nil {
		return []user.Trader{}, err
	}
	req, err := fhttp.NewRequest("POST", binanceApi, bytes.NewReader(postData))
	if err != nil {
		return []user.Trader{}, err
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
		return []user.Trader{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []user.Trader{}, err
	}

	var binanceApiResponse LeaderboardResponse
	err = json.Unmarshal(body, &binanceApiResponse)
	if err != nil {
		app.Logger.Error("Error decoding Binance API response " + err.Error())
		return []user.Trader{}, err
	}

	switch code := binanceApiResponse.Code.(type) {
	case float64:
		if code == 400 {
			err := errors.New(fmt.Sprintf("failed connecting to binance API, status code: %.0f", code))
			app.Logger.Error(err)
			return nil, err
		}
	default:
		app.Logger.Debug(fmt.Sprintf("collected %d traders from binance API", len(binanceApiResponse.Data)))
	}

	return binanceApiResponse.Data, nil
}

func (app *Config) fetchTraderPositions(uid string) ([]user.Position, error) {

	binanceApi := "https://www.binance.com/bapi/futures/v1/public/future/leaderboard/getOtherPosition"

	postJson := make(map[string]any)
	var postData []byte
	postJson["encryptedUid"] = uid
	postJson["tradeType"] = "PERPETUAL"
	postData, err := json.MarshalIndent(postJson, " ", "")
	if err != nil {
		return []user.Position{}, err
	}
	req, err := fhttp.NewRequest("POST", binanceApi, bytes.NewReader(postData))
	if err != nil {
		return []user.Position{}, err
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
		return []user.Position{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []user.Position{}, err
	}

	var binanceApiResponse GetOtherPositionResponse
	err = json.Unmarshal(body, &binanceApiResponse)
	if err != nil {
		app.Logger.Error("Error decoding Binance API response " + err.Error())
		return []user.Position{}, err
	}

	switch code := binanceApiResponse.Code.(type) {
	case float64:
		if code == 400 {
			err := errors.New(fmt.Sprintf("failed connecting to binance API, status code: %.0f", code))
			app.Logger.Error(err)
			return nil, err
		}
	default:
		app.Logger.Debug(fmt.Sprintf("collected %d traders from binance API", len(binanceApiResponse.Data)))
	}

	_, ok := binanceApiResponse.Data["otherPositionRetList"]
	// If the key exists
	if ok {
		var positions []user.Position
		byteData, _ := json.Marshal(binanceApiResponse.Data["otherPositionRetList"])
		err = json.Unmarshal(byteData, &positions)
		if err != nil {
			app.Logger.Error("Error decoding Binance API response " + err.Error())
			return []user.Position{}, err
		}
		return positions, nil
	} else {
		err = fmt.Errorf("failed parsing positions for trader %s: key otherPositionRetList does not exist", uid)
		app.Logger.Debug(err)
		return []user.Position{}, err
	}

}

func (app *Config) searchByNickname(nickname string) ([]user.Trader, error) {

	binanceApi := "https://www.binance.com/bapi/futures/v1/public/future/leaderboard/searchNickname"

	postJson := make(map[string]any)
	var postData []byte
	postJson["nickname"] = nickname

	postData, err := json.MarshalIndent(postJson, " ", "")
	if err != nil {
		return []user.Trader{}, err
	}
	req, err := fhttp.NewRequest("POST", binanceApi, bytes.NewReader(postData))
	if err != nil {
		return []user.Trader{}, err
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
		return []user.Trader{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []user.Trader{}, err
	}

	var binanceApiResponse LeaderboardResponse
	err = json.Unmarshal(body, &binanceApiResponse)
	if err != nil {
		app.Logger.Error("Error decoding Binance API response " + err.Error())
		return []user.Trader{}, err
	}

	switch code := binanceApiResponse.Code.(type) {
	case float64:
		if code == 400 {
			err := errors.New(fmt.Sprintf("failed connecting to binance API, status code: %.0f", code))
			app.Logger.Error(err)
			return nil, err
		}
	default:
		app.Logger.Debug(fmt.Sprintf("collected %d traders from binance API", len(binanceApiResponse.Data)))
	}

	return binanceApiResponse.Data, nil
}
