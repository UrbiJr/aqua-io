package ws

type Position struct {
	Arg struct {
		Channel  string `json:"channel"`
		Uid      string `json:"uid"`
		InstType string `json:"instType"`
	} `json:"arg"`
	Data []struct {
		Adl            string `json:"adl"`
		AvailPos       string `json:"availPos"`
		AvgPx          string `json:"avgPx"`
		CTime          string `json:"cTime"`
		Ccy            string `json:"ccy"`
		DeltaBS        string `json:"deltaBS"`
		DeltaPA        string `json:"deltaPA"`
		GammaBS        string `json:"gammaBS"`
		GammaPA        string `json:"gammaPA"`
		Imr            string `json:"imr"`
		InstId         string `json:"instId"`
		InstType       string `json:"instType"`
		Interest       string `json:"interest"`
		Last           string `json:"last"`
		UsdPx          string `json:"usdPx"`
		Lever          string `json:"lever"`
		Liab           string `json:"liab"`
		LiabCcy        string `json:"liabCcy"`
		LiqPx          string `json:"liqPx"`
		MarkPx         string `json:"markPx"`
		Margin         string `json:"margin"`
		MgnMode        string `json:"mgnMode"`
		MgnRatio       string `json:"mgnRatio"`
		Mmr            string `json:"mmr"`
		NotionalUsd    string `json:"notionalUsd"`
		OptVal         string `json:"optVal"`
		PTime          string `json:"pTime"`
		Pos            string `json:"pos"`
		BaseBorrowed   string `json:"baseBorrowed"`
		BaseInterest   string `json:"baseInterest"`
		QuoteBorrowed  string `json:"quoteBorrowed"`
		QuoteInterest  string `json:"quoteInterest"`
		PosCcy         string `json:"posCcy"`
		PosId          string `json:"posId"`
		PosSide        string `json:"posSide"`
		SpotInUseAmt   string `json:"spotInUseAmt"`
		SpotInUseCcy   string `json:"spotInUseCcy"`
		BizRefId       string `json:"bizRefId"`
		BizRefType     string `json:"bizRefType"`
		ThetaBS        string `json:"thetaBS"`
		ThetaPA        string `json:"thetaPA"`
		TradeId        string `json:"tradeId"`
		UTime          string `json:"uTime"`
		Upl            string `json:"upl"`
		UplRatio       string `json:"uplRatio"`
		VegaBS         string `json:"vegaBS"`
		VegaPA         string `json:"vegaPA"`
		CloseOrderAlgo []struct {
			AlgoId          string `json:"algoId"`
			SlTriggerPx     string `json:"slTriggerPx"`
			SlTriggerPxType string `json:"slTriggerPxType"`
			TpTriggerPx     string `json:"tpTriggerPx"`
			TpTriggerPxType string `json:"tpTriggerPxType"`
			CloseFraction   string `json:"closeFraction"`
		} `json:"closeOrderAlgo"`
	} `json:"data"`
}
