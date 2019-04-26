package manager

import "housekeeper2/Tools"

type RespOrder interface {
	GetInstrumentID() string
	GetOpenOrClose() byte
	GetPosiDirection() byte
	GetPosition() int
}

type Order interface {
	DirctIndex() int
	SetDirect(byte)
	GetDirect() byte
	ConvPosition(*Position) error
	SetInstrumentId(string)
	GetInstrumentId() string
	SetNum(int)
	GetNum() int
	SetYdOrTd(byte)
	GetYdOrTd() byte
	SetExchangId(string)
	GetExchangId() string
	SetOrderPriceType(byte)
	GetOrderPriceType() byte
	SetLimitPrice(float64)
	GetLimitPrice() float64
	SetCombOffsetFlag(string)
	GetCombOffsetFlag() string
	SetTimeCondition(byte)
	GetTimeCondition() byte
	SetCombHedgeFlag(string)
	GetCombHedgeFlag() string
}

type MyOrder struct {
	direct         byte //0 表示空 1表示多
	instrumentID   string
	num            int
	ydortd         byte //1 昨仓 2今仓 0 全部
	exchangId      string
	orderPriceType byte
	limitPrice     float64
	combOffsetFlag string
	combhedgeflag  string
	timeCondition  byte
}

func NewOrder(flag byte) Order {
	return &MyOrder{ydortd: flag}
}

func (o *MyOrder) ConvPosition(pst *Position) error {
	o.SetInstrumentId(pst.GetContract())
	if o.ydortd == 1 {
		o.SetNum(pst.GetYestodayNum())
	} else if o.ydortd == 2 {
		o.SetNum(pst.GetNum() - pst.GetYestodayNum())
	} else {
		o.SetNum(pst.GetNum())
	}
	o.SetDirect(pst.GetDirect())
	return nil
}

func (o *MyOrder) SetCombHedgeFlag(b string) {
	o.combhedgeflag = b
}

func (o *MyOrder) GetCombHedgeFlag() string {
	return o.combhedgeflag
}

func (o *MyOrder) SetDirect(d byte) {
	o.direct = d
}

func (o *MyOrder) GetDirect() byte {
	return o.direct
}

func (o *MyOrder) DirctIndex() int {
	if o.GetDirect() == 0x30 {
		return 0
	}
	return 1
}

func (o *MyOrder) SetInstrumentId(contract string) {
	o.instrumentID = contract
}

func (o *MyOrder) GetInstrumentId() string {
	return o.instrumentID
}

func (o *MyOrder) SetNum(num int) {
	o.num = num
}

func (o *MyOrder) GetNum() int {
	return o.num
}

func (o *MyOrder) SetYdOrTd(b byte) {
	o.ydortd = b
}

func (o *MyOrder) GetYdOrTd() byte {
	return o.ydortd
}

func (o *MyOrder) SetExchangId(id string) {
	o.exchangId = id
}

func (o *MyOrder) GetExchangId() string {
	return o.exchangId
}

func (o *MyOrder) SetOrderPriceType(b byte) {
	o.orderPriceType = b
}

func (o *MyOrder) GetOrderPriceType() byte {
	return o.orderPriceType
}

func (o *MyOrder) SetLimitPrice(price float64) {
	o.limitPrice = Tools.Decimal(price)
}

func (o *MyOrder) GetLimitPrice() float64 {
	return o.limitPrice
}

func (o *MyOrder) SetCombOffsetFlag(flag string) {
	o.combOffsetFlag = flag
}

func (o *MyOrder) GetCombOffsetFlag() string {
	return o.combOffsetFlag
}

func (o *MyOrder) SetTimeCondition(b byte) {
	o.timeCondition = b
}

func (o *MyOrder) GetTimeCondition() byte {
	return o.timeCondition
}
