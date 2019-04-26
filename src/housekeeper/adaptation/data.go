package adaptation

type PositionField struct {
	InstrumentID  string
	PosiDirection byte
	YdPosition    int
	Position      int
	OpenOrClose   byte
	HighestPrice  float64
	LowestPrice   float64
}

type TradingAccount struct {
	Balance    float64
	CurrMargin float64
}

type Quote struct {
	TradingDay         string
	InstrumentID       string
	ExchangeId         string
	ExchangeInstId     string
	LastPrice          float64
	PreSettlementPrice float64
	PreClosePrice      float64
	PreOpenInterest    float64
	OpenPrice          float64
	HighestPrice       float64
	LowestPrice        float64
	Volume             int
}
