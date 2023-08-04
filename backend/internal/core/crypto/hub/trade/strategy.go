package tradeTypes

import (
	"sort"
	"sync"
)

type TradeStrategy struct {
	StopTradingControl bool
	AutoTakeProfit     map[int]map[string]float64 // index – key: val (position size & percentage)
	AutoStopLoss       map[int]map[string]float64 // index – key: val (position size & percentage)
	TPmu               sync.RWMutex
	SLmu               sync.RWMutex

	CurrentTPIndex int
	CurrentSLIndex int
	TPReached      chan int
	SLReached      chan int
}

func New(isStopTrading bool) *TradeStrategy {
	return &TradeStrategy{
		StopTradingControl: isStopTrading,
		AutoStopLoss:       make(map[int]map[string]float64),
	}
}

func (b *TradeStrategy) CreateStrategy() {

}

func (b *TradeStrategy) organizeTradeStrategy() {

	type IndexedStrategy struct {
		Index        int
		PositionSize float64
		Percentage   float64
	}

	validStrategies := []IndexedStrategy{}

	for index, strategy := range b.AutoTakeProfit {
		posSize := strategy["positionSize"]
		percentage := strategy["percentage"]

		if posSize >= 0 && posSize <= 100 && percentage >= 0 && percentage <= 100 {
			validStrategies = append(validStrategies, IndexedStrategy{Index: index, PositionSize: posSize, Percentage: percentage})
		}
	}

	sort.Slice(validStrategies, func(i, j int) bool {
		return validStrategies[i].PositionSize < validStrategies[j].PositionSize
	})

	organizedStrategies := make(map[int]map[string]float64)
	for _, strategy := range validStrategies {
		organizedStrategies[strategy.Index] = map[string]float64{
			"positionSize": strategy.PositionSize,
			"percentage":   strategy.Percentage,
		}
	}

	// either return a new safemap or delete the current one and add the correct strategy

	b.handleStrategy()
}

func (b *TradeStrategy) setStrategy() {

}

// handleStrategy establishes a WS connection and
// updates the TP/SL of a rest.
func (b *TradeStrategy) handleStrategy() {
	//go b.TradeStrategy.WebsocketConn.SubscribePosition()
}

func (b *TradeStrategy) subscribe(sub Subscription, socketUrl string) error {
	switch sub {
	case Balance:

	}
	return nil
}

// handleBalance subscribes to the balance of a
// ByBit Account and closes if it reaches a threshold.
func (b *TradeStrategy) monitorBalance() {
	//go b.TradeStrategy.WebsocketConn.SubscribeBalance(b.Balance.Value, config.GlobalConfig.ByBit[b.Index].StopIfFallUnder, b.TraderID)
}
