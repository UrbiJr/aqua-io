package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/UrbiJr/aqua-io/internal/core/crypto/common"
	"github.com/UrbiJr/aqua-io/internal/core/crypto/constants"
	http "github.com/bogdanfinn/fhttp"
	"golang.org/x/sync/errgroup"
	"io"
	"net/url"
	"strings"
)

type Binance struct {
	*common.Trade
}

func New(trade *common.Trade) *Binance {
	return &Binance{trade}
}

func (b *Binance) HandlePosition(event common.OrderType) error {
	var payload url.Values

	price, err := b.getCoinPrice()
	if price == -1 {
		return fmt.Errorf("ERROR Parsing %s Ticker Price [== -1]", b.Info.Coin.Symbol)
	} else if err != nil {
		return err
	}

	switch event {
	case common.OpenPos:
		grp1, _ := errgroup.WithContext(context.TODO())
		grp2, _ := errgroup.WithContext(context.TODO())

		b.assert()

		grp1.Go(func() error {
			ok := b.Trade.ComparePrice()
			if !ok {
				return fmt.Errorf("price difference is too high between exchanges")
			}
			return nil
		})

		grp1.Go(func() error {
			for _, coin := range config.CopyTradingCfg.Traders[b.Misc.Index].BlackListedCoins {
				if strings.ToUpper(coin) == b.Info.Coin.Symbol {
					return fmt.Errorf("blacklisted coin match [%s]", coin)
				}
			}
			return nil
		})

		grp1.Go(func() error {
			err = b.getBalance()
			return err
		})

		if err = grp1.Wait(); err != nil {
			return err
		}

		grp2.Go(func() error {
			b.Info.Quantity = b.Trade.CalculateQuantity()
			return nil
		})

		grp2.Go(func() error {
			err = b.setLeverage()
			return err
		})

		if err = grp2.Wait(); err != nil {
			return err
		}

		payload.Add("symbol", b.Info.Coin.Symbol)
		payload.Add("side", b.Info.Side)
		payload.Add("quantity", fmt.Sprint(b.Info.Quantity))
		payload.Add("price", fmt.Sprint(price))

	case common.ClosePos:
		payload.Add("closePosition", "true")
		payload.Add("symbol", b.Info.Coin.Symbol)
		payload.Add("side", b.Info.Side)
		payload.Add("quantity", fmt.Sprint(b.Info.Quantity))
		payload.Add("price", fmt.Sprint(price))

	case common.DecreasePos:
		qty := b.CalculateQuantity()

		payload.Add("symbol", b.Info.Coin.Symbol)
		payload.Add("side", b.Info.Side)
		payload.Add("price", fmt.Sprint(price))
		payload.Add("quantity", fmt.Sprint(qty))

		if b.Info.Side == constants.BINANCE_BUY {
			payload.Add("side", constants.BINANCE_SELL)
		} else {
			payload.Add("side", constants.BINANCE_BUY)
		}
		b.Info.Quantity -= qty

	case common.IncreasePos:
		grp1, _ := errgroup.WithContext(context.TODO())

		grp1.Go(func() error {
			var pnl int
			pnl, err = b.getPNL()
			if err != nil || pnl == -1 {
				return fmt.Errorf("unable to get PnL [%w]", err)
			}
			ok := b.BlockAddsIfPositive()
			if !ok {
				return fmt.Errorf("event aborted [blockAddsIfPositive ON & triggered]")
			}
			return nil
		})

		grp1.Go(func() error {
			return nil
		})

		qty := b.CalculateQuantity()

		payload.Add("quantity", fmt.Sprint(qty))
		b.Info.Quantity += qty
	}

	payload.Add("timestamp", now())

	resp, err := b.Do(http.MethodPost, payload, positionURL)
	if err != nil {
		return errors.HandleERR(constants.Binance, fmt.Errorf("[%s] Unknown CLIENT issue [%s]", event, err))
	}

	defer resp.Body.Close()

	err = errors.HandleStatus(resp.StatusCode)
	if err != nil {
		return errors.HandleERR(constants.Binance, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Info("BINANCE resp PosHandler", string(body))

	var position rest.Position
	err = json.Unmarshal(body, &position)
	if err != nil {
		return err
	}

	if position.Symbol == b.Info.Coin.Symbol {
		log.Success(constants.Binance, fmt.Sprintf("Success [%s]!", event))
		return nil
	} else {
		return fmt.Errorf("err doing stuf w pos LAST STEP CHECKING")
	}
}
