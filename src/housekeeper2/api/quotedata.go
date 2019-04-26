package api

type Quoter interface {
	SetVolume(int)
	GetVolume() int
	SetLowestPrice(float64)
	GetLowestPrice() float64
	SetHighestPrice(float64)
	GetHighestPrice() float64
	SetOpenPrice(float64)
	GetOpenPrice() float64
	SetPreOpenInterest(float64)
	GetPreOpenInterest() float64
	SetPreClosePrice(float64)
	GetPreClosePrice() float64
	SetPreSettlementPrice(float64)
	GetPreSettlementPrice() float64
	SetLastPrice(float64)
	GetLastPrice() float64
	SetExchangeInstId(string)
	GetExchangeInstId() string
	SetExchangeId(string)
	GetExchangeId() string
	SetInstrumentID(string)
	GetInstrumentID() string
	SetTradingDay(string)
	GetTradingDay() string
}

type MyQuoter struct {
	tradingDay         string
	instrumentID       string
	exchangeId         string
	exchangeInstId     string
	lastPrice          float64
	preSettlementPrice float64
	preClosePrice      float64
	preOpenInterest    float64
	openPrice          float64
	highestPrice       float64
	lowestPrice        float64
	volume             int
}

func NewQuoter() Quoter {
	return &MyQuoter{}
}

func (m *MyQuoter) SetVolume(d int) {
	m.volume = d
}

func (m *MyQuoter) GetVolume() int {
	return m.volume
}

func (m *MyQuoter) SetLowestPrice(d float64) {
	m.lowestPrice = d
}

func (m *MyQuoter) GetLowestPrice() float64 {
	return m.lowestPrice
}

func (m *MyQuoter) SetHighestPrice(d float64) {
	m.highestPrice = d
}

func (m *MyQuoter) GetHighestPrice() float64 {
	return m.highestPrice
}

func (m *MyQuoter) SetOpenPrice(d float64) {
	m.openPrice = d
}

func (m *MyQuoter) GetOpenPrice() float64 {
	return m.openPrice
}

func (m *MyQuoter) SetPreOpenInterest(d float64) {
	m.preOpenInterest = d
}

func (m *MyQuoter) GetPreOpenInterest() float64 {
	return m.preOpenInterest
}

func (m *MyQuoter) SetPreClosePrice(d float64) {
	m.preClosePrice = d
}

func (m *MyQuoter) GetPreClosePrice() float64 {
	return m.preClosePrice
}

func (m *MyQuoter) SetPreSettlementPrice(d float64) {
	m.preSettlementPrice = d
}

func (m *MyQuoter) GetPreSettlementPrice() float64 {
	return m.preSettlementPrice
}

func (m *MyQuoter) SetLastPrice(d float64) {
	m.lastPrice = d
}

func (m *MyQuoter) GetLastPrice() float64 {
	return m.lastPrice
}

func (m *MyQuoter) SetExchangeInstId(d string) {
	m.exchangeInstId = d
}

func (m *MyQuoter) GetExchangeInstId() string {
	return m.exchangeInstId
}

func (m *MyQuoter) SetExchangeId(d string) {
	m.exchangeId = d
}

func (m *MyQuoter) GetExchangeId() string {
	return m.exchangeId
}

func (m *MyQuoter) SetInstrumentID(d string) {
	m.instrumentID = d
}

func (m *MyQuoter) GetInstrumentID() string {
	return m.instrumentID
}

func (m *MyQuoter) SetTradingDay(d string) {
	m.tradingDay = d
}

func (m *MyQuoter) GetTradingDay() string {
	return m.tradingDay
}
