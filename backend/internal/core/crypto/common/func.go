package common

import (
	"math"
	"time"
)

// CalculateQuantity function calculates the quantity of a cryptocurrency that should be traded based on the maximum coin allocation specified in the trade settings and the current value of the user's wallet.
func (t *Trade) CalculateQuantity() float64 {
	notionalSize := t.Info.TradeSettings.MaxCoinAllocation / 100 * t.Info.Wallet.Value
	return notionalSize / t.Info.Coin.ExchangePrice
}

// ComparePrice function compares the current price of a cryptocurrency on both the Binance exchange and the user's exchange to determine if the price difference is within the maximum price difference specified in the trade settings.
func (t *Trade) ComparePrice() bool {
	difference := math.Abs(t.Info.Coin.BinancePrice - t.Info.Coin.ExchangePrice)
	if difference <= t.Info.TradeSettings.MaxPriceDifferenceBetweenExchange {
		return true
	} else {
		return false
	}
}

// BlockAddsIfPositive function checks if the profit and loss (PnL) of the current trade is positive, and if so, returns true to prevent adding a new position.
func (t *Trade) BlockAddsIfPositive() bool {
	if t.Info.PnL > 0 {
		return true
	} else {
		return false
	}
}

// IsOpenDelayMet function checks if the delay between opening new positions has been met, and if so, returns true to allow adding a new position.
func (t *Trade) IsOpenDelayMet() bool {
	t.Info.TradeSettings.OpenDelayBetweenPositions.Value = config.CopyTradingCfg.Traders[t.Misc.Index].OpenDelayBetweenPositions
	now := time.Now().Unix()
	if t.Info.TradeSettings.OpenDelayBetweenPositions.LastTimeAction == 0 {
		t.Info.TradeSettings.OpenDelayBetweenPositions.LastTimeAction = now
		return true
	}
	if now >= t.Info.TradeSettings.OpenDelayBetweenPositions.SleepUntil {
		t.Info.TradeSettings.OpenDelayBetweenPositions.LastTimeAction = now
		t.Info.TradeSettings.OpenDelayBetweenPositions.SleepUntil = now + int64(t.Info.TradeSettings.OpenDelayBetweenPositions.Value)
		return true
	}
	return false
}

// ShouldAddPosition function checks if the current price of the cryptocurrency is below a certain price threshold, and if so, returns true to allow adding a new position.
func (t *Trade) ShouldAddPosition(avgEntryPrice float64) bool {
	preventionPrice := avgEntryPrice * (1 - t.Info.TradeSettings.AddPreventionPercent/100)
	if t.Info.Coin.ExchangePrice < preventionPrice {
		return true
	} else {
		return false
	}
}

// CalculateMaxAddAmount function calculates the maximum amount that can be added to the current position based on the current position size and the maximum add multiplier specified in the trade settings.
func (t *Trade) CalculateMaxAddAmount(currentPosition float64, addAmount float64) float64 {
	maxAddAmount := currentPosition * t.Info.TradeSettings.MaxAddMultiplier.ConfigVal
	if addAmount > maxAddAmount {
		return maxAddAmount
	} else {
		return addAmount
	}
}
