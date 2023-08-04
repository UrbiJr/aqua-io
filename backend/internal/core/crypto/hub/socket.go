package hub

import (
	"encoding/json"
	"fmt"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/common"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/constants"
	bit "github.com/UrbiJr/aqua-io/backend/internal/core/crypto/exchanges/bybit/futures/client/rest"
	okx "github.com/UrbiJr/aqua-io/backend/internal/core/crypto/exchanges/okx/futures/client/rest"
	"github.com/gorilla/websocket"
	"math"
	"net/http"
	"time"
)

var dialer = &websocket.Dialer{HandshakeTimeout: time.Second * 10}

const (
	binanceTrades string = "wss://api.allierobotics.com/crypto/socket/"
)

type Websocket struct {
	Exchange   string
	Connected  bool
	Connection *websocket.Conn
}

func connect() (*websocket.Conn, string, error) {

	var traders []string
	for _, trader := range config.CopyTradingCfg.Traders {
		for _, trd := range trader.Traders {
			traders = append(traders, trd)
		}
	}

	payloadBytes, err := json.Marshal(traders)
	if err != nil {
		return nil, "", err
	}

	payload := string(payloadBytes)

	conn, resp, err := dialer.Dial(binanceTrades, http.Header{
		"X-Api-Key":  {config.GlobalCfg.License},
		"X-Traders":  {payload},
		"X-Metadata": {hwid.Get()},
	})
	if err != nil {
		return nil, payload, err
	}

	if resp.StatusCode != 101 {
		return nil, payload, fmt.Errorf("ERROR Establishing Connection to the Server")
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, payload, err
	}

	return conn, payload, nil
}

func (ws *Websocket) disconnect() {
	err := ws.Connection.Close()
	if err != nil {
		panic(err)
	}
}

func handleTrades(exchange constants.Exchange) {
	conn, payload, err := connect()
	if err != nil {
		//log.Error("Socket_Broadcast", err)
		return
	}

	go func(conn *websocket.Conn) {
		/*
			if err = discord.SendNotificationWebhook(discord.WebhookNotification{
				Title:       types.BotConnected,
				Description: "__Successfully Connected__ to the Broadcasting Server.\n> ––Awaiting Trades <a:acheck:1016117530991014049>",
				Embeds: []types.EmbedFields{
					{
						Name:   "Exchange",
						Value:  string(exchange),
						Inline: true,
					},
					{
						Name:   "Traders",
						Value:  payload,
						Inline: true,
					},
				},
			}); err != nil {
				return
			}
		*/

		for {
			var msg common.ExchangeTrade

			err = conn.ReadJSON(&msg)
			if err != nil {
				time.Sleep(1500 * time.Millisecond)
				continue
			}

			if msg.Event == "ping" {
				err = conn.WriteJSON(map[string]string{"event": "pong"})
				if err != nil {
					time.Sleep(1500 * time.Millisecond)
					continue
				}
			} else {
				var ok bool
				switch exchange {
				/*⚠️🐳 OKX 🐳⚠️*/
				case constants.OKX:
					switch msg.Event {
					case common.OpenPos:
						trade := NewTrade(&msg)
						okxTrade := okx.New(trade)
						err = okxTrade.HandlePosition(common.OpenPos)
						if err != nil {
							//log.Error(constants.OKX, err)
							return
						}

						OKX_ORDERBOOK.AddTrade(msg.Trade.Symbol, okxTrade)

					case common.ClosePos:
						trade := OKX_ORDERBOOK.GetTrade(msg.Trade.Symbol)
						err = trade.HandlePosition(common.ClosePos)
						if err != nil {
							//log.Error(constants.OKX, err)
							return
						}

						OKX_ORDERBOOK.Trades.Delete(trade.Info.Coin.Symbol)

					case common.UpdatePos:
						var trade *okx.OKX
						trade, ok = OKX_ORDERBOOK.Trades.Get(msg.Trade.Symbol)
						if !ok {
							return
						}

						if math.Abs(msg.Trade.Amount) > trade.Trade.Info.Metrics.PositionChange[1] {
							trade.Event = common.IncreasePos
							trade.Info.Quantity = msg.Trade.Amount - trade.Info.Quantity
						} else {
							trade.Event = common.DecreasePos
							trade.Info.Quantity = trade.Info.Quantity - msg.Trade.Amount
						}

						err = trade.HandlePosition(common.UpdatePos)
						if err != nil {
							//log.Error(constants.OKX, err)
							return
						}
					default:
						return
					}
					/*⚠️🐳 BYBIT 🐳⚠️*/
				case constants.ByBit:
					switch msg.Event {
					case common.OpenPos:
						trade := NewTrade(&msg)
						bitTrade := bit.New(trade)
						err = bitTrade.HandlePosition(common.OpenPos)
						if err != nil {
							//log.Error(constants.ByBit, err)
							return
						}

						BIT_ORDERBOOK.AddTrade(msg.Trade.Symbol, bitTrade)

					case common.ClosePos:
						trade := BIT_ORDERBOOK.GetTrade(msg.Trade.Symbol)
						err = trade.HandlePosition(common.ClosePos)
						if err != nil {
							//log.Error(constants.ByBit, err)
							return
						}

						BIT_ORDERBOOK.Trades.Delete(trade.Info.Coin.Symbol)

					case common.UpdatePos:
						var trade *bit.ByBit
						trade, ok = BIT_ORDERBOOK.Trades.Get(msg.Trade.Symbol)
						if !ok {
							return
						}

						if math.Abs(msg.Trade.Amount) > trade.Trade.Info.Metrics.PositionChange[1] {
							trade.Event = common.IncreasePos
							trade.Info.Quantity = msg.Trade.Amount - trade.Info.Quantity
						} else {
							trade.Event = common.DecreasePos
							trade.Info.Quantity = trade.Info.Quantity - msg.Trade.Amount
						}

						err = trade.HandlePosition(common.UpdatePos)
						if err != nil {
							//log.Error(constants.ByBit, err)
							return
						}
					default:
						return
					}
					/*⚠️🐳 BINANCE 🐳⚠️*/
				case constants.Binance:
					//todo add support for binance
				default:
					//log.Error("Unsupported Exchange")
				}
			}
		}
	}(conn)

}
