package socket

import (
	"allie/backend/pkg/config"
	"allie/backend/pkg/discord"
	"allie/backend/pkg/discord/types"
	"allie/backend/pkg/map"
	"allie/backend/tui/log"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	privateSocket = "wss://stream.bybit.com/v5/private"
	positionURI   = "position"
)

type SubscriptionType string

const (
	trades SubscriptionType = "trades"
	wallet SubscriptionType = "wallet"
)

type ByBitWebsocket struct {
	CoinSubscription CoinSubscription
	WebsocketConns   *safemap.safemap[string, *WebsocketConnection]
}

type WebsocketConnection struct {
	PublicAPI       string
	SecretAPI       string
	WebsocketConn   *websocket.Conn
	Balance         float64
	ConnID          string
	IsConnected     bool
	IsTradeStrategy bool
}

type CoinSubscription struct {
	PublicAPI      string
	SecretAPI      string
	WebsocketConn  *websocket.Conn
	CoinsMonitored *safemap.SafeMap[string, bool] //symbol: isMonitored(bool)
	CoinChan       chan string
}

func InitByBitWebsocket() {
	ws := newByBitWebsocket()
	errors := ws.CreateAllWebsocketConnections()
	if len(errors) != 0 {
		for _, err := range errors {
			discordErr := discord.SendNotificationWebhook(discord.WebhookNotification{
				Title:       types.Critical,
				Description: fmt.Sprintf("%s", err.Error()),
			})
			if discordErr != nil {
				log.Error(err.Error())
			}
		}
	}
}

func newByBitWebsocket() *ByBitWebsocket {
	return &ByBitWebsocket{
		WebsocketConns:   safemap.New[string, *WebsocketConnection](),
		CoinSubscription: CoinSubscription{
			//PublicAPI: config.GlobalConfig.ByBit[0].PublicAPI,
			//SecretAPI: config.GlobalConfig.ByBit[0].SecretAPI,
		},
	}
}

// CreateAllWebsocketConnections initializes a websocket connection on
// all the ByBit account stored in the config file.
func (b *ByBitWebsocket) CreateAllWebsocketConnections() []error {
	var errors []error
	for _, bybit := range config.CopyTradingCfg.ByBit {
		ws := &WebsocketConnection{
			PublicAPI: bybit.PublicAPI,
			SecretAPI: bybit.SecretAPI,
		}

		err := ws.Connect()
		if err != nil {
			errors = append(errors, err)
		} else if err == nil {
			b.WebsocketConns.Set(ws.PublicAPI, ws)
		}
	}
	return errors
}

func (ws *WebsocketConnection) Connect() error {
	var err error
	var resp *http.Response
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(signMessage(ws.PublicAPI, ws.SecretAPI))
	if err != nil {
		return err
	}

	ws.WebsocketConn, resp, err = websocket.DefaultDialer.Dial("wss://stream.bybit.com/v5/private", nil)
	if err != nil {
		ws.IsConnected = false
		return err
	}

	if resp.StatusCode != 101 {
		return fmt.Errorf("ERROR Connecting to ByBit Private Socket [%s]", resp.Status)
	}

	defer ws.WebsocketConn.Close()

	err = ws.WebsocketConn.WriteMessage(websocket.TextMessage, buf.Bytes())
	if err != nil {
		return err
	}

	_, message, err := ws.WebsocketConn.ReadMessage()
	if err != nil {
		return err
	}

	var response map[string]interface{}
	err = json.Unmarshal(message, &response)
	if err != nil {
		return err
	}

	if response["success"].(bool) {
		ws.ConnID = response["conn_id"].(string)
		ws.IsConnected = true
		fmt.Println()
		log.Success(fmt.Sprintf("WebsocketConnection: %t [%s]", ws.IsConnected, ws.ConnID))
		return nil
	} else {
		ws.IsConnected = false
		return fmt.Errorf("WebsocketConnection.connect: ERROR Connecting [%s]", response["ret_msg"].(string))
	}

}

func (ws *WebsocketConnection) subscribe(subType SubscriptionType) error {
	var err error
	var payload = make(map[string]interface{})

	payload["op"] = "subscribe"

	switch subType {
	case trades:

		payload["args"] = []SubscriptionType{trades}
		err = ws.WebsocketConn.WriteJSON(payload)
		if err != nil {
			return err
		}

		ws.SubscribePosition()
	case wallet:

		payload["args"] = []SubscriptionType{wallet}
		err = ws.WebsocketConn.WriteJSON(payload)
		if err != nil {
			return err
		}

		ws.SubscribeBalance(1, 2, "")
	}

	return nil
}

// SubscribeBalance subscribes to the USDT value of the trading balance.
func (ws *WebsocketConnection) SubscribeBalance(balance float64, maxLoss float64, traderUID string) {

	type Balance struct {
		ID           string `json:"id"`
		Topic        string `json:"topic"`
		CreationTime int64  `json:"creationTime"`
		Data         []struct {
			AccountIMRate          string `json:"accountIMRate"`
			AccountMMRate          string `json:"accountMMRate"`
			TotalEquity            string `json:"totalEquity"`
			TotalWalletBalance     string `json:"totalWalletBalance"`
			TotalMarginBalance     string `json:"totalMarginBalance"`
			TotalAvailableBalance  string `json:"totalAvailableBalance"`
			TotalPerpUPL           string `json:"totalPerpUPL"`
			TotalInitialMargin     string `json:"totalInitialMargin"`
			TotalMaintenanceMargin string `json:"totalMaintenanceMargin"`
			Coin                   []struct {
				Coin                string `json:"coin"`
				Equity              string `json:"equity"`
				UsdValue            string `json:"usdValue"`
				WalletBalance       string `json:"walletBalance"`
				AvailableToWithdraw string `json:"availableToWithdraw"`
				AvailableToBorrow   string `json:"availableToBorrow"`
				BorrowAmount        string `json:"borrowAmount"`
				AccruedInterest     string `json:"accruedInterest"`
				TotalOrderIM        string `json:"totalOrderIM"`
				TotalPositionIM     string `json:"totalPositionIM"`
				TotalPositionMM     string `json:"totalPositionMM"`
				UnrealisedPnl       string `json:"unrealisedPnl"`
				CumRealisedPnl      string `json:"cumRealisedPnl"`
				Bonus               string `json:"bonus"`
			} `json:"coin"`
			AccountType string `json:"accountType"`
			AccountLTV  string `json:"accountLTV"`
		} `json:"data"`
	}

	var initialBalance = balance

	for {

		_, msg, err := ws.WebsocketConn.ReadMessage()
		if err != nil {
			continue
		}

		var response Balance
		if err = json.Unmarshal(msg, &response); err != nil {
			continue
		}

		for _, coin := range response.Data[0].Coin {
			if strings.ToUpper(coin.Coin) == "USDT" {

				var usdValue float64

				usdValue, err = strconv.ParseFloat(coin.UsdValue, 64)
				if err != nil {
					log.Error(err.Error())
					continue
				}

				percentageChange := (usdValue - initialBalance) / initialBalance * 100

				if percentageChange <= maxLoss {
					//handler.OrderBook.CloseAllTraderPositions(traderUID)

					desc := fmt.Sprintf("All trades related to this ByBit account have been closed as the balance has reached %f.", usdValue)
					discordErr := discord.SendNotificationWebhook(discord.WebhookNotification{
						Title:       types.Critical,
						Description: desc,
						Color:       types.Red,
						IsEmbed:     true,
						Embeds: []types.EmbedFields{
							{
								Name:   "Current Balance",
								Value:  fmt.Sprintf("%f USDT", usdValue),
								Inline: true,
							},
							{
								Name:   "Initial Balance",
								Value:  strconv.FormatFloat(initialBalance, 'f', 0, 64) + "USDT",
								Inline: true,
							},
							{
								Name:   "Balance Setting 'StopIfFallUnder'",
								Value:  strconv.FormatFloat(maxLoss, 'f', 0, 64) + "%",
								Inline: true,
							},
						},
					})

					if discordErr != nil {
						log.Error(discordErr.Error())
					}

					os.Exit(0)
				}
			}
		}
	}
}

func (ws *WebsocketConnection) SubscribePosition() {

	type Position struct {
		ID           string `json:"id"`
		Topic        string `json:"topic"`
		CreationTime int64  `json:"creationTime"`
		Data         []struct {
			PositionIdx      int    `json:"positionIdx"`
			TradeMode        int    `json:"tradeMode"`
			RiskID           int    `json:"riskId"`
			RiskLimitValue   string `json:"riskLimitValue"`
			Symbol           string `json:"symbol"`
			Side             string `json:"side"`
			Size             string `json:"size"`
			EntryPrice       string `json:"entryPrice"`
			Leverage         string `json:"leverage"`
			PositionValue    string `json:"positionValue"`
			MarkPrice        string `json:"markPrice"`
			PositionIM       string `json:"positionIM"`
			PositionMM       string `json:"positionMM"`
			TakeProfit       string `json:"takeProfit"`
			StopLoss         string `json:"stopLoss"`
			TrailingStop     string `json:"trailingStop"`
			UnrealisedPnl    string `json:"unrealisedPnl"`
			CumRealisedPnl   string `json:"cumRealisedPnl"`
			CreatedTime      string `json:"createdTime"`
			UpdatedTime      string `json:"updatedTime"`
			TpslMode         string `json:"tpslMode"`
			LiqPrice         string `json:"liqPrice"`
			BustPrice        string `json:"bustPrice"`
			Category         string `json:"category"`
			PositionStatus   string `json:"positionStatus"`
			AdlRankIndicator int    `json:"adlRankIndicator"`
		} `json:"data"`
	}

	for {

		_, message, err := ws.WebsocketConn.ReadMessage()
		if err != nil {
			log.Error(fmt.Sprintf("SubscribePosition: %s", err.Error()))
			continue
		}

		var response Position
		if err = json.Unmarshal(message, &response); err != nil {
			log.Error(fmt.Sprintf("SubscribePosition: %s", err.Error()))
			continue
		}

		if response.Topic != "position" {
			continue
		}

		//SL := response.Data[0].StopLoss
		//TP := response.Data[0].TakeProfit

	}

}

// TODO en gros dès qu'on reçoit une func tradecrypto on check si on monitor une currency, si elle existe pas on commence à la moni
func (b *CoinSubscription) SubscribeCoin() {

	for {
		coin := <-b.CoinChan

		coinURI := fmt.Sprintf("tickers.%s", coin)

		err := b.WebsocketConn.WriteMessage(1, []byte(coinURI))
		if err != nil {
			//log.Error(err)

		}

		// todo move to a chan

		//go ReadCoinPrice()
	}

}

// TODO review
func (b *ByBitWebsocket) checkCoin() bool {

	coin := <-b.CoinSubscription.CoinChan
	_, ok := b.CoinSubscription.CoinsMonitored.Get(coin)
	return ok

	//b.CoinSubscription.CoinMu.RLock()
	//defer b.CoinSubscription.CoinMu.RUnlock()
	//coin := <-b.CoinSubscription.CoinChan
	//_, ok := b.CoinSubscription.CoinsMonitored[coin]
	//return ok
}

func (ws *WebsocketConnection) ping() bool {

	var response map[string]interface{}

	pingMessage := map[string]string{
		"op": "ping",
	}

	err := ws.WebsocketConn.WriteJSON(pingMessage)
	if err != nil {
		return false
	}

	_, message, err := ws.WebsocketConn.ReadMessage()
	if err != nil {
		return false
	}

	if err = json.Unmarshal(message, &response); err != nil {
		return false
	}

	if response["op"] == "pong" {
		return true
	} else {
		return false
	}
}

func signMessage(publicAPI, secretAPI string) map[string]interface{} {
	expires := int64(time.Now().UnixNano()/1e6) + 1000
	message := fmt.Sprintf("GET/realtime%d", expires)
	signature := hmac.New(sha256.New, []byte(secretAPI))
	signature.Write([]byte(message))
	sign := hex.EncodeToString(signature.Sum(nil))
	authMessage := map[string]interface{}{
		"op":   "auth",
		"args": []interface{}{publicAPI, expires, sign},
	}
	return authMessage
}

type Wallet struct {
	ID           string `json:"id"`
	Topic        string `json:"topic"`
	CreationTime int64  `json:"creationTime"`
	Data         []struct {
		AccountIMRate          string `json:"accountIMRate"`
		AccountMMRate          string `json:"accountMMRate"`
		TotalEquity            string `json:"totalEquity"`
		TotalWalletBalance     string `json:"totalWalletBalance"`
		TotalMarginBalance     string `json:"totalMarginBalance"`
		TotalAvailableBalance  string `json:"totalAvailableBalance"`
		TotalPerpUPL           string `json:"totalPerpUPL"`
		TotalInitialMargin     string `json:"totalInitialMargin"`
		TotalMaintenanceMargin string `json:"totalMaintenanceMargin"`
		Coin                   []struct {
			Coin                string `json:"coin"`
			Equity              string `json:"equity"`
			UsdValue            string `json:"usdValue"`
			WalletBalance       string `json:"walletBalance"`
			AvailableToWithdraw string `json:"availableToWithdraw"`
			AvailableToBorrow   string `json:"availableToBorrow"`
			BorrowAmount        string `json:"borrowAmount"`
			AccruedInterest     string `json:"accruedInterest"`
			TotalOrderIM        string `json:"totalOrderIM"`
			TotalPositionIM     string `json:"totalPositionIM"`
			TotalPositionMM     string `json:"totalPositionMM"`
			UnrealisedPnl       string `json:"unrealisedPnl"`
			CumRealisedPnl      string `json:"cumRealisedPnl"`
			Bonus               string `json:"bonus"`
		} `json:"coin"`
		AccountType string `json:"accountType"`
	} `json:"data"`
}

type Ticker struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Data  struct {
		Symbol            string `json:"symbol"`
		TickDirection     string `json:"tickDirection"`
		Price24HPcnt      string `json:"price24hPcnt"`
		LastPrice         string `json:"lastPrice"`
		PrevPrice24H      string `json:"prevPrice24h"`
		HighPrice24H      string `json:"highPrice24h"`
		LowPrice24H       string `json:"lowPrice24h"`
		PrevPrice1H       string `json:"prevPrice1h"`
		MarkPrice         string `json:"markPrice"`
		IndexPrice        string `json:"indexPrice"`
		OpenInterest      string `json:"openInterest"`
		OpenInterestValue string `json:"openInterestValue"`
		Turnover24H       string `json:"turnover24h"`
		Volume24H         string `json:"volume24h"`
		NextFundingTime   string `json:"nextFundingTime"`
		FundingRate       string `json:"fundingRate"`
		Bid1Price         string `json:"bid1Price"`
		Bid1Size          string `json:"bid1Size"`
		Ask1Price         string `json:"ask1Price"`
		Ask1Size          string `json:"ask1Size"`
	} `json:"data"`
	Cs int64 `json:"cs"`
	Ts int64 `json:"ts"`
}
