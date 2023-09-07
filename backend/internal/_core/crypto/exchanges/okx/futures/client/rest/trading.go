package rest

import (
	"bytes"
	"context"
	"encoding/json"
	errorsn "errors"
	"fmt"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/common"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/constants"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/exchanges/okx/futures/models/rest"
	http "github.com/bogdanfinn/fhttp"
	"golang.org/x/sync/errgroup"
	"io"
	"strings"
)

var (
	upgradeAccountMode = errorsn.New("ERROR: Please Upgrade Your OKX Account to Single Currency Margin [ALL UPCOMING TRADES WILL NOT BE EXECUTED ON THIS ACCOUNT]")
	invalidPassPhrase  = errorsn.New("ERROR: Your PassPhrase is invalid – Check its correctness or create a new API [ALL UPCOMING TRADES WILL FAIL ON THIS ACCOUNT]")
	busyServers        = errorsn.New("ERROR: Systems are busy! [Please try again later]")
)

type OKX struct {
	*common.Trade
}

func New(trade *common.Trade) *OKX {
	return &OKX{trade}
}

func (o *OKX) HandlePosition(event common.OrderType) error {
	var requestURL string
	var ok bool
	var buf bytes.Buffer
	var payload = map[string]interface{}{}

	err := o.getInstrumentsLeverage()
	if err != nil {
		return err
	}

	switch event {
	case common.OpenPos:
		grp1, _ := errgroup.WithContext(context.TODO())
		grp2, _ := errgroup.WithContext(context.TODO())

		o.assert()

		grp1.Go(func() error {
			ok = o.Trade.ComparePrice()
			if !ok {
				return fmt.Errorf("price Difference is too high between exchanges")
			}
			return nil
		})

		grp1.Go(func() error {
			for _, coin := range config.CopyTradingCfg.Traders[o.Misc.Index].BlackListedCoins {
				if strings.ToUpper(coin) == o.Info.Coin.Symbol {
					return fmt.Errorf("blacklisted coin match [%s]", coin)
				}
			}
			return nil
		})

		grp1.Go(func() error {
			err = o.getBalance()
			return err
		})

		if err = grp1.Wait(); err != nil {
			return err
		}

		grp2.Go(func() error {
			o.Trade.CalculateQuantity()
			return nil
		})

		grp2.Go(func() error {
			err = o.setLeverage()
			return err
		})

		if err = grp2.Wait(); err != nil {
			return err
		}

		payload = map[string]interface{}{
			"tdMode":  o.Info.TradeMode,
			"instId":  o.Info.Coin.Symbol,
			"side":    o.Info.Side,
			"ccy":     "USDT",
			"ordType": "limit",
			"sz":      o.Info.Quantity,
			"px":      fmt.Sprint(o.Info.Coin.ExchangePrice),
		}

		requestURL = openPositionURL

	case common.UpdatePos: // we alr deducted if it was a reduce or

		payload = map[string]interface{}{
			"instId":  o.Trade.Info.Coin.Symbol,
			"posSide": o.Trade.Info.Side,
			"amt":     o.Trade.Info.Quantity,
		}

		// todo: review
		if o.Trade.Event == common.IncreasePos {
			payload["type"] = "add"
		} else {
			payload["type"] = "reduce"
		}

		grp1, _ := errgroup.WithContext(context.TODO())

		grp1.Go(func() error {
			ok = o.Trade.BlockAddsIfPositive()
			if ok {
				return fmt.Errorf("ERROR Aborting DecreasePos > BlockPositionAdds=true")
			}
			return nil
		})

		grp1.Go(func() error {
			err = o.getCoinPrice()
			return err
		})

		if err = grp1.Wait(); err != nil {
			return err
		}

		requestURL = openPositionURL

	case common.ClosePos:

		payload = map[string]interface{}{
			"instId":  o.Trade.Info.Coin.Symbol,
			"mgnMode": o.Trade.Info.TradeMode,
			"autoCxl": true,
			"ccy":     "USDT",
		}

		err = o.ClosePosition()
		return err
	default:
		return fmt.Errorf("OrderType: invalid value, want [openPos, closePos, increasePos, decreasePos]")
	}

	err = json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := o.Do(http.MethodPost, requestURL, buf.String(), true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.HandleERR(constants.OKX, errors.HandleStatus(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var position rest.Order
	err = json.Unmarshal(body, &position)
	if err != nil {
		return err
	}

	if position.Code == "0" && position.Msg == "" {
		if position.Msg == "Request unsupported under current account mode" {
			return upgradeAccountMode
		}
		if requestURL == closePositionURL {
			err = o.GetPnL()
			if err != nil {
				return err
			}

			if position.Msg == "Position doesn't exist" {
				log.Info(constants.OKX, fmt.Sprintf("Order [%s | %s] not Filled – Cancelling...", o.Info.Coin.Symbol, o.Info.Side))
				err = o.cancelOrder()
				if err != nil {
					return err
				}
				return nil
			}

			log.Success(constants.OKX, fmt.Sprintf("Position Closed [p:%d]", o.Trade.Info.PnL))
			return nil
		}

		log.Success(constants.OKX, fmt.Sprintf("Success [%s | %s | %s]!", event, o.Info.Coin.Symbol, o.Info.Side))
		err = discord.SendTradeWebhook(types.WebhookTrade{
			ExchangeName: constants.OKX,
			Event:        "event",
			TraderID:     o.Info.TraderID,
			TradeType:    types.TradeOpen,
			TraderName:   o.Info.TraderName,
			Coin:         o.Info.Coin.Symbol,
			Side:         o.Info.Side,
			Price:        fmt.Sprint(o.Info.Coin.ExchangePrice),
			PositionSize: fmt.Sprint(o.Info.Quantity),
			Leverage:     fmt.Sprintf("%d", int(o.Info.Leverage)),
		})
		return nil
	} else {
		return errors.HandleERR(constants.OKX, fmt.Errorf("ERROR Opening Position [%s]", position.Data[0].SMsg))
	}
}

func (o *OKX) OpenPosition() error {

	var buf bytes.Buffer
	var openPositionURI = "/api/v5/trade/order"

	payload := map[string]interface{}{
		"tdMode":  strings.ToLower(config.CopyTradingCfg.Traders[o.Trade.Misc.Index].TradeMode),
		"instId":  o.Trade.Info.Coin.Symbol,
		"side":    o.Trade.Info.Side,
		"ccy":     "USDT",
		"ordType": "limit",
		"sz":      o.Trade.Info.Quantity,
		"px":      fmt.Sprint(o.Trade.Info.Coin.ExchangePrice),
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := o.Do(http.MethodPost, openPositionURI, buf.String(), true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Opening Position [%s]", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var position rest.Order
	err = json.Unmarshal(body, &position)
	if err != nil {
		return err
	}

	if position.Code == "0" && position.Msg == "" {
		o.Trade.Info.TradeID = position.Data[0].OrdId
		log.Success(constants.OKX, fmt.Sprintf("Position Open [ccy:%s|qty:%s|s:%s]", o.Trade.Info.Coin.Symbol, o.Trade.Info.Quantity, o.Trade.Info.Side))
		return nil
	} else if position.Msg == "Request unsupported under current account mode" {
		return upgradeAccountMode
	} else {
		return fmt.Errorf("ERROR Opening Position [%s]", position.Data[0].SMsg)
	}
}

func (o *OKX) ClosePosition() error {

	var buf bytes.Buffer
	var closePositionURI = "/api/v5/trade/close-position"

	payload := map[string]interface{}{
		"instId":  o.Trade.Info.Coin.Symbol,
		"mgnMode": o.Trade.Info.TradeMode,
		"autoCxl": true,
		"ccy":     "USDT",
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := o.Do(http.MethodPost, closePositionURI, buf.String(), true)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Closing Position [%s]", resp.Status)
	}

	var position rest.Close
	err = json.NewDecoder(resp.Body).Decode(&position)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if position.Code == "0" && position.Msg == "" {
		err = o.GetPnL()
		if err != nil {
			return err
		}
		log.Success(fmt.Sprintf("Position Closed [%d]", o.Trade.Info.PnL))
	} else if position.Msg == "Position doesn't exist" {
		log.Info("Position not Filled – Cancelling...")
		err = o.cancelOrder()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("ERROR Closing Position [%s]", position.Msg)
	}
	return nil
}

func (o *OKX) UpdatePosition() error {

	var buf bytes.Buffer

	payload := map[string]interface{}{
		"instId":  o.Trade.Info.Coin.Symbol,
		"posSide": o.Trade.Info.Side,
		"amt":     o.Trade.Info.Quantity,
	}

	if o.Trade.Event == common.IncreasePos {
		payload["type"] = "add"
	} else {
		payload["type"] = "reduce"
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := o.Do(http.MethodPost, updatePositionMargin, buf.String(), true)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Inc/Dec Position [%s]", resp.Status)
	}

	var position rest.Position
	err = json.NewDecoder(resp.Body).Decode(&position)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	if position.Code == "0" && position.Msg == "" {
		log.Success(fmt.Sprintf("Position Updated [%s]", payload["type"]))
		return nil
	} else if position.Msg == "Request unsupported under current account mode" {
		return upgradeAccountMode
	} else {
		return fmt.Errorf("ERROR Updating Position [%s | %s]", payload["type"], o.Trade.Info.Coin.Symbol)
	}
}
