package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/common"
	"github.com/UrbiJr/aqua-io/backend/internal/core/crypto/constants"
	http "github.com/bogdanfinn/fhttp"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
	"io"
	"strings"
)

const (
	LimitType  = "Limit"
	LinearType = "linear"
)

type ByBit struct {
	*common.Trade
}

func New(trade *common.Trade) *ByBit {
	return &ByBit{trade}
}

func (b *ByBit) HandlePosition(event common.OrderType) error {

	var buf bytes.Buffer
	var payload = map[string]interface{}{
		"category":    LinearType,
		"symbol":      b.Trade.Info.Coin.Symbol,
		"side":        b.Trade.Info.Side,
		"orderType":   LimitType,
		"orderLinkId": uuid.NewString(),
	}

	switch event {
	case common.OpenPos:
		grp1, _ := errgroup.WithContext(context.TODO())
		grp2, _ := errgroup.WithContext(context.TODO())

		b.assert()

		grp1.Go(func() error {
			ok := b.Trade.ComparePrice()
			if !ok {
				return fmt.Errorf("price Difference is too high between exchanges")
			}
			return nil
		})

		grp1.Go(func() error {

			for _, coin := range config.CopyTradingCfg.Traders[b.Misc.Index].BlackListedCoins {
				if strings.ToUpper(coin) == b.Info.Coin.Symbol {
					return fmt.Errorf("blacklisted coin match [%s]", coin)
				}
			}
			log.Success("did black listed coin")
			return nil
		})

		grp1.Go(func() error {
			err := b.getBalance()
			return err
		})

		if err := grp1.Wait(); err != nil {
			return err
		}

		grp2.Go(func() error {
			b.Info.Quantity = b.Trade.CalculateQuantity()
			return nil
		})

		grp2.Go(func() error {
			err := b.setLeverage()
			return err
		})

		if err := grp2.Wait(); err != nil {
			return err
		}

		payload["qty"] = b.Info.Quantity
		payload["price"] = b.Info.Coin.ExchangePrice

		//todo add rest of handling
	case common.ClosePos:
	case common.IncreasePos:
	case common.DecreasePos:
	default:

	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := b.Do(http.MethodPost, positionURL, buf.String())
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var position rest.Position
	err = json.Unmarshal(body, &position)
	if err != nil {
		return err
	}

	if position.RetCode != 0 {
		return errors.HandleERR(constants.ByBit, fmt.Errorf("%s", position.RetMsg))
	} else {
		log.Success(constants.ByBit, fmt.Sprintf("Success [%s]!", event))
		err = discord.SendTradeWebhook(types.WebhookTrade{
			ExchangeName: constants.ByBit,
			Event:        string(event),
			TraderID:     b.Info.TraderID,
			TradeType:    types.TradeOpen,
			TraderName:   b.Info.TraderName,
			Coin:         b.Info.Coin.Symbol,
			Side:         b.Info.Side,
			Price:        fmt.Sprint(b.Info.Coin.ExchangePrice),
			PositionSize: fmt.Sprint(b.Info.Quantity),
			Leverage:     fmt.Sprintf("%d", int(b.Info.Leverage)),
		})
		if err != nil {
			log.Error(constants.ByBit, err)
		}
		return nil
	}

}

/*
func (b *ByBit) setTakeProfitStopLoss() error {

	var buf bytes.Buffer

	currentTP := b.TradeStrategy.AutoTakeProfit[b.CurrentTPIndex]["takeProfit"]
	currentSL := b.TradeStrategy.AutoStopLoss[b.CurrentTPIndex]["stopLoss"]

	payload := safemap[string]interface{}{
		"category":   LimitType,
		"symbol":     b.Coin.Symbol,
		"takeProfit": currentTP,
		"stopLoss":   currentSL,
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	tpslURI, err := url.Parse(baseURI + "position/trading-stop")
	if err != nil {
		return err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    tpslURI,
		Body:   io.NopCloser(&buf),
		Header: b.V5SignPOST(buf.Bytes()),
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Setting TP/SL [%s]", resp.Status)
	}

	body := bodyString(io.ReadAll(resp.Body))

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	retMsg := gjson.Get(body, "retMsg").String()

	if retMsg != "OK" {
		return fmt.Errorf("ERROR Setting TP/SL [%s]", retMsg)
	}

	log.Info(fmt.Sprintf("TP[%f]/SL[%f] Set for %s", currentTP, currentSL, b.Coin.Symbol))

	return nil
}
*/

// MoveToBreakEven partially closes a position and move
// the rest to break even by changing SL.
func (b *ByBit) MoveToBreakEven() error {

	var buf bytes.Buffer

	payload := map[string]interface{}{
		"category": LimitType,
		"symbol":   b.Info.Coin.Symbol,
		"stopLoss": b.Info.Coin.ExchangePrice,
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	resp, err := b.Do(http.MethodPost, moveToBreakEven, buf.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ERROR Moving Pos to Break Even [%s]", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	retMsg := gjson.Get(string(body), "retMsg").String()
	if retMsg == "OK" {
		log.Success(fmt.Sprintf("Successfully Moved [%s] Trade to Break Even", b.Trade.Info.Coin.Symbol))
		return nil
	} else {
		return fmt.Errorf("ERROR Moving Pos to Break Even [%s]", retMsg)
	}
}

/*
func (b *ByBit) monitorTPSL() {
	defer b.wg.Done()

	for {
		select {
		case tpIndex := <-b.tpReached:
			b.currentTPIndex = tpIndex
			err := b.setTPSL()
			if err != nil {
				fmt.Println("Error setting TP:", err)
			}
		case slIndex := <-b.slReached:
			b.currentSLIndex = slIndex
			err := b.setTPSL()
			if err != nil {
				fmt.Println("Error setting SL:", err)
			}
		}
	}
}
*/

// Leverage is an integer which we have to convert to a float.
/*
func (b *ByBit) convertLeverage() {
	leverage := config.CopyTradingCfg.Traders[b.Trade.Misc.Index].Leverage
	b.Trade.Info.Leverage = b.Trade.Info.Quantity * leverage / 100
}
*/

// initialTradeQuantity calculates the position's size.
/*
func (b *ByBit) initialTradeQuantity() {
	b.Trade.Info.Coin.Symbol
	notionalSize := config.CopyTradingCfg.Traders[b.Trade.Misc.Index].MaxCoinAllocation / 100 * b.Balance.Value
	b.Quantity = fmt.Sprintf("%f", notionalSize/b.Trade.Coin.Value)
}
*/

/*pas bon*/
// When a trader increases its position.
/*
func (b *ByBit) calculateMaxAddMultiplier() error {
	quantity := b.Binance.Trade.Amount
	maxAdd := quantity * config.GlobalConfig.Traders[b.Index].MaxAddMultiplier
	if quantity <= maxAdd {
		return fmt.Errorf("calculateMaxAddMultiplier: Quantity calculated is lower than maxAdd") // todo handle plus tard
	} else {
		b.MaxAddMultiplier = maxAdd / quantity

		log.Info(fmt.Sprintf("MaxAddMultipler: %f", b.MaxAddMultiplierResult)) // temp

		return nil
	}
}
*/

// shouldAddToPosition checks if, based on the user's choice, we should increase/decrease a position.
func (b *ByBit) shouldAddToPosition(addPreventionPercent float64, currentAvgPrice float64, currentPrice float64) bool {
	priceThreshold := currentAvgPrice * (1.0 - addPreventionPercent/100.0)
	if currentPrice < priceThreshold {
		return true
	} else {
		return false
	}
}

// TODO test why it fails sometimes
/*
func (b *ByBit) checkPriceDifference() (bool, float64, error) {

	var err error
	var bybitCurrency, binanceCurrency float64
	maxDifference := config.GlobalConfig.Traders[b.Index].MaxPriceDifferenceBetweenExchange
	binanceCurrency = b.Binance.Trade.EntryPrice

	err = b.GetCoinPrice()
	if err != nil {
		return false, -1, err
	}

	bybitCurrency, err = strconv.ParseFloat(b.Coin.ByBitPrice, 64)
	if err != nil {
		return false, -1, err
	}

	b.Coin.Value = bybitCurrency

	difference := math.Abs(binanceCurrency - bybitCurrency)
	if difference <= maxDifference {
		return true, difference, nil
	} else {
		return false, difference, nil
	}
}
*/

// OpenDelayBetweenPositions ✅
/*
func (b *ByBit) isOpenDelayMet() bool {
	b.OpenDelayBetweenPositions.Value = config.GlobalConfig.Traders[b.Index].OpenDelayBetweenPositions

	now := time.Now().Unix()
	if b.OpenDelayBetweenPositions.LastTimeAction == 0 {
		b.OpenDelayBetweenPositions.LastTimeAction = now
		return true
	}

	if now >= b.OpenDelayBetweenPositions.SleepUntil {
		b.OpenDelayBetweenPositions.LastTimeAction = now
		b.OpenDelayBetweenPositions.SleepUntil = now + int64(b.OpenDelayBetweenPositions.Value)
		return true
	}

	return false
}
*/

// TODO review if that func is useful
// add timestamp to the safemap to check ✅
/*
func (b *ByBit) addTimestampsToMap() {
	if _, ok := v5.PreviousOpenings[b.TraderID]; !ok {
		v5.PreviousOpenings[b.TraderID] = []*time.Time{}
	}
	timestamp := time.Unix(0, b.Binance.Data.UpdateTimeStamp*int64(time.Millisecond))
	b.LatestTradeTimestamp = &timestamp
	PreviousOpenings[b.TraderID] = append(PreviousOpenings[b.TraderID], &timestamp)
}

func (b *ByBit) removeTimestampFromMap() {
	for k, v := range PreviousOpenings {
		if k == b.TraderID {
			for _, val := range v {
				if val == b.LatestTradeTimestamp {
					// delete a value from a slice and not a safemap

					delete(PreviousOpenings[b.TraderID], val)
				}
			}
		}
	}
}
*/

/*
func (b *ByBit) addTradeStrategy() {

	var wg sync.WaitGroup
	wg.Add(2)

	// If the size of AutoTP/SL is lower or equal to 1, we return as we won't need to monitor the position and make a strategy.
	if len(config.GlobalConfig.Traders[b.Index].AutoTakeProfit) <= 0 && len(config.GlobalConfig.Traders[b.Index].AutoStopLoss) <= 0 {
		return
	}

	go func() {
		data := make(safemap[string]float64)
		for index, tp := range config.GlobalConfig.Traders[b.Index].AutoTakeProfit {
			b.TPmu.RLock()
			data["positionSize"] = tp.Strategy.PositionSize
			data["percentage"] = tp.Strategy.Percentage
			b.AutoTakeProfit[index] = data
			b.TPmu.Unlock()
		}
		wg.Done()
	}()

	go func() {
		data := make(safemap[string]float64)
		for index, sl := range config.GlobalConfig.Traders[b.Index].AutoStopLoss {
			b.SLmu.RLock()
			data["positionSize"] = sl.Strategy.PositionSize
			data["percentage"] = sl.Strategy.Percentage
			b.AutoStopLoss[index] = data
			b.SLmu.Unlock()
		}
		wg.Done()
	}()

	wg.Wait()

	b.organizeTradeStrategy()
}

func (b *ByBit) updateAction(action OrderType) {
	b.LatestAction = action
}

// TODO review if needed
// calculateTPSL1 calculates the first TP and SL that is set when opening a rest, then we let the WS handle that.
func (b *ByBit) calculateTPSL1() {

}

func bodyString(body []byte, _ error) string {
	return string(body)
}

func hashByte(body io.ReadCloser) []byte {
	result, _ := io.ReadAll(body)
	return result
}
*/
