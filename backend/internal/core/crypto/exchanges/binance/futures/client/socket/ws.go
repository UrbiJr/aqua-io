package socket

import (
	"allie/backend/modules/copytrade/crypto/constants"
	"allie/backend/modules/copytrade/crypto/exchanges/binance/futures/models/socket"
	"bytes"
	"encoding/json"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

var (
	websocketClient, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger())
)

const (
	getListenKey = "/fapi/v1/listenKey"
)

const (
	monitorBalanceMethod = "account.status"
	monitorOrdersMethod  = "allOrders"
)

type Socket struct {
	APIkey            string
	ListenKey         string
	BalanceChannel    chan<- string            // USDT Balance
	PositionsChannels map[string]chan<- string // Symbol:PositionOpen
	mainConn          *websocket.Conn
	Conns             map[string]*websocket.Conn //UUID: conn
}

func Connect(balanceChannel chan string, positionsChannel string) {

	socket := &Socket{}

	err := socket.generateListenKey()
	if err != nil {

	}

	go keepAliveListenKey(socket.ListenKey)

}

func (w *Socket) MonitorPosition(symbol string) {
	var buf bytes.Buffer
	payload := socket.PositionsPayload{
		ID:     uuid.NewString(),
		Method: monitorOrdersMethod,
		Params: map[string]interface{}{
			"symbol": symbol,
			"limit":  1,
		},
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {

	}

	err = w.mainConn.WriteJSON(&buf)
	if err != nil {

	}

	go func() {
		var msg []byte
		for {
			_, msg, err = w.mainConn.ReadMessage()
			if err != nil {
				w.BalanceChannel <- "ERR"
			}

			var resBalance socket.Positions
			err = json.Unmarshal(msg, &resBalance)
			if err != nil {
				w.BalanceChannel <- "ERR"
			}

		}
	}()
}

func (w *Socket) MonitorBalance() {
	var buf bytes.Buffer
	payload := socket.BalancePayload{
		ID:     uuid.NewString(),
		Method: monitorBalanceMethod,
		Params: map[string]interface{}{
			"apiKey":    w.APIkey,
			"signature": "",
			"timestamp": 1, // use "now" func
		},
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {

	}

	err = w.mainConn.WriteJSON(&buf)
	if err != nil {

	}

	go func() {
		var msg []byte
		for {
			_, msg, err = w.mainConn.ReadMessage()
			if err != nil {
				w.BalanceChannel <- "ERR"
			}

			var resBalance socket.Balances
			err = json.Unmarshal(msg, &resBalance)
			if err != nil {
				w.BalanceChannel <- "ERR"
			}

			for _, balance := range resBalance.Result.Balances {
				if balance.Asset == "USDT" {
					w.BalanceChannel <- balance.Free
				}
			}
		}
	}()
}

func (w *Socket) pingHandler() {
	err := w.mainConn.WriteJSON(map[string]interface{}{})
	if err != nil {
		//handle
	}
}

func (w *Socket) generateListenKey() error {

	URL, _ := url.Parse(constants.BinancebaseRestAPI + getListenKey)

	req := &http.Request{
		Method: http.MethodPost,
		URL:    URL,
	}

	resp, err := websocketClient.Do(req)
	if err != nil {
		return err
	}

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	w.ListenKey = response["listenKey"]

	return nil
}

func keepAliveListenKey(listenKey string) {
	ticker := time.NewTicker(58 * time.Minute)

	keepAlive := func() {
		URL, _ := url.Parse(constants.BinancebaseRestAPI + getListenKey)

		req := &http.Request{
			Method: http.MethodPut,
			URL:    URL,
		}

		resp, err := websocketClient.Do(req)
		if err != nil {

		}

		if resp.StatusCode != 200 {
			// log that there was an err keeping alive key
		}
	}

	for {
		select {
		case <-ticker.C:
			keepAlive()
		}
	}
}
