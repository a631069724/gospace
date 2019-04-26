package api

type Request interface {
	SetBrokeId(string)
	BrokeId() string
	SetName(string)
	UserName() string
	SetPwd(string)
	Pwd() string
	SetId(int)
	ReqId() int
	SetInstrumentID(string)
	InstrumentID() string
	SetCur(string)
	Cur() string
	SetExchangeId(string)
	GetExchangeId() string
	SetCombHedgeFlag(string)
	GetCombHedgeFlag() string
	SetOrderPriceType(byte)
	GetOrderPriceType() byte
	SetDirection(byte)
	GetDirection() byte
	SetCombOffSetFlag(string)
	GetCombOffSetFlag() string
	SetNum(int)
	GetNum() int
	SetLimitPrice(float64)
	GetLimitPrice() float64
}

type ReqData struct {
	brokeid        string
	name           string
	pwd            string
	id             int
	instrument     string
	cur            string
	exchangeID     string
	orderpricetype byte
	direction      byte
	comboffsetflag string
	combhedgeflag  string
	num            int
	limitprice     float64
}

func (r *ReqData) SetLimitPrice(n float64) {
	r.limitprice = n
}

func (r *ReqData) GetLimitPrice() float64 {
	return r.limitprice
}

func (r *ReqData) SetNum(n int) {
	r.num = n
}

func (r *ReqData) GetNum() int {
	return r.num
}

func (r *ReqData) SetDirection(b byte) {
	r.direction = b
}

func (r *ReqData) GetDirection() byte {
	return r.direction
}

func (r *ReqData) SetCombHedgeFlag(b string) {
	r.combhedgeflag = b
}

func (r *ReqData) GetCombHedgeFlag() string {
	return r.combhedgeflag
}

func (r *ReqData) SetCombOffSetFlag(b string) {
	r.comboffsetflag = b
}

func (r *ReqData) GetCombOffSetFlag() string {
	return r.comboffsetflag
}

func (r *ReqData) SetOrderPriceType(b byte) {
	r.orderpricetype = b
}

func (r *ReqData) GetOrderPriceType() byte {
	return r.orderpricetype
}

func (r *ReqData) SetExchangeId(id string) {
	r.exchangeID = id
}

func (r *ReqData) GetExchangeId() string {
	return r.exchangeID
}

func (r *ReqData) SetInstrumentID(instrument string) {
	r.instrument = instrument
}

func (r *ReqData) InstrumentID() string {
	return r.instrument
}

func (r *ReqData) SetBrokeId(brokeid string) {
	r.brokeid = brokeid
}

func (r *ReqData) BrokeId() string {
	return r.brokeid
}

func (r *ReqData) SetName(name string) {
	r.name = name
}

func (r *ReqData) UserName() string {
	return r.name
}

func (r *ReqData) SetPwd(pwd string) {
	r.pwd = pwd
}

func (r *ReqData) Pwd() string {
	return r.pwd
}

func (r *ReqData) SetId(id int) {
	r.id = id
}

func (r *ReqData) ReqId() int {
	return r.id
}

func (r *ReqData) SetCur(cur string) {
	r.cur = cur
}

func (r *ReqData) Cur() string {
	return r.cur
}

func NewReqData() Request {
	return &ReqData{}
}
