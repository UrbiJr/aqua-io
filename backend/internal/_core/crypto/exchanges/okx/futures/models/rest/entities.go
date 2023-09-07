package rest

type Order struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		ClOrdId string `json:"clOrdId"`
		OrdId   string `json:"ordId"`
		Tag     string `json:"tag"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

type Close struct {
	Code string `json:"code"`
	Data []struct {
		ClOrdId string `json:"clOrdId"`
		InstId  string `json:"instId"`
		PosSide string `json:"posSide"`
		Tag     string `json:"tag"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Position struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Amt      string `json:"amt"`
		Ccy      string `json:"ccy"`
		InstId   string `json:"instId"`
		Leverage string `json:"leverage"`
		PosSide  string `json:"posSide"`
		Type     string `json:"type"`
	} `json:"data"`
}

type Cancel struct {
	Code string `json:"code"`
	Data []struct {
		ClOrdId string `json:"clOrdId"`
		OrdId   string `json:"ordId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type PositionsHistory struct {
	Code string `json:"code"`
	Data []struct {
		CTime         string `json:"cTime"`
		Ccy           string `json:"ccy"`
		CloseAvgPx    string `json:"closeAvgPx"`
		CloseTotalPos string `json:"closeTotalPos"`
		Direction     string `json:"direction"`
		InstId        string `json:"instId"`
		InstType      string `json:"instType"`
		Lever         string `json:"lever"`
		MgnMode       string `json:"mgnMode"`
		OpenAvgPx     string `json:"openAvgPx"`
		OpenMaxPos    string `json:"openMaxPos"`
		Pnl           string `json:"pnl"`
		PnlRatio      string `json:"pnlRatio"`
		PosId         string `json:"posId"`
		TriggerPx     string `json:"triggerPx"`
		Type          string `json:"type"`
		UTime         string `json:"uTime"`
		Uly           string `json:"uly"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Balance struct {
	Code string `json:"code"`
	Data []struct {
		AdjEq   string `json:"adjEq"`
		Details []struct {
			AvailBal      string `json:"availBal"`
			AvailEq       string `json:"availEq"`
			CashBal       string `json:"cashBal"`
			Ccy           string `json:"ccy"`
			CrossLiab     string `json:"crossLiab"`
			DisEq         string `json:"disEq"`
			Eq            string `json:"eq"`
			EqUsd         string `json:"eqUsd"`
			FrozenBal     string `json:"frozenBal"`
			Interest      string `json:"interest"`
			IsoEq         string `json:"isoEq"`
			IsoLiab       string `json:"isoLiab"`
			IsoUpl        string `json:"isoUpl"`
			Liab          string `json:"liab"`
			MaxLoan       string `json:"maxLoan"`
			MgnRatio      string `json:"mgnRatio"`
			NotionalLever string `json:"notionalLever"`
			OrdFrozen     string `json:"ordFrozen"`
			Twap          string `json:"twap"`
			UTime         string `json:"uTime"`
			Upl           string `json:"upl"`
			UplLiab       string `json:"uplLiab"`
			StgyEq        string `json:"stgyEq"`
			SpotInUseAmt  string `json:"spotInUseAmt"`
		} `json:"details"`
		Imr         string `json:"imr"`
		IsoEq       string `json:"isoEq"`
		MgnRatio    string `json:"mgnRatio"`
		Mmr         string `json:"mmr"`
		NotionalUsd string `json:"notionalUsd"`
		OrdFroz     string `json:"ordFroz"`
		TotalEq     string `json:"totalEq"`
		UTime       string `json:"uTime"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Ticker struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstType  string `json:"instType"`
		InstId    string `json:"instId"`
		Last      string `json:"last"`
		LastSz    string `json:"lastSz"`
		AskPx     string `json:"askPx"`
		AskSz     string `json:"askSz"`
		BidPx     string `json:"bidPx"`
		BidSz     string `json:"bidSz"`
		Open24H   string `json:"open24h"`
		High24H   string `json:"high24h"`
		Low24H    string `json:"low24h"`
		VolCcy24H string `json:"volCcy24h"`
		Vol24H    string `json:"vol24h"`
		SodUtc0   string `json:"sodUtc0"`
		SodUtc8   string `json:"sodUtc8"`
		Ts        string `json:"ts"`
	} `json:"data"`
}

type Leverage struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Lever   string `json:"lever"`
		MgnMode string `json:"mgnMode"`
		InstId  string `json:"instId"`
		PosSide string `json:"posSide"`
	} `json:"data"`
}

type Instruments struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstType     string `json:"instType"`
		InstId       string `json:"instId"`
		InstFamily   string `json:"instFamily"`
		Uly          string `json:"uly"`
		Category     string `json:"category"`
		BaseCcy      string `json:"baseCcy"`
		QuoteCcy     string `json:"quoteCcy"`
		SettleCcy    string `json:"settleCcy"`
		CtVal        string `json:"ctVal"`
		CtMult       string `json:"ctMult"`
		CtValCcy     string `json:"ctValCcy"`
		OptType      string `json:"optType"`
		Stk          string `json:"stk"`
		ListTime     string `json:"listTime"`
		ExpTime      string `json:"expTime"`
		Lever        string `json:"lever"`
		TickSz       string `json:"tickSz"`
		LotSz        string `json:"lotSz"`
		MinSz        string `json:"minSz"`
		CtType       string `json:"ctType"`
		Alias        string `json:"alias"`
		State        string `json:"state"`
		MaxLmtSz     string `json:"maxLmtSz"`
		MaxMktSz     string `json:"maxMktSz"`
		MaxTwapSz    string `json:"maxTwapSz"`
		MaxIcebergSz string `json:"maxIcebergSz"`
		MaxTriggerSz string `json:"maxTriggerSz"`
		MaxStopSz    string `json:"maxStopSz"`
	} `json:"data"`
}
