package api

type Response interface {
	SetFrontID(int)
	SetSessionID(int)
	SetMaxOrderRef(string)
	SetInstrumentID(string)
	SetPosiDirection(byte)
	SetYdPosition(int)
	SetPosition(int)
	SetUseMargin(float64)
	SetOpenAmount(float64)
	SetCloseAmount(float64)
	SetPositionCost(float64)
	SetCloseProfit(float64)
	SetPositionProfit(float64)
	SetOpenCost(float64)
	SetExchangeMargin(float64)
	SetAccountMargin(float64)
	SetBalance(float64)

	GetFrontID() int
	GetSessionID() int
	GetMaxOrderRef() string
	GetInstrumentID() string
	GetPosiDirection() byte
	GetYdPosition() int
	GetPosition() int
	GetUseMargin() float64
	GetOpenAmount() float64
	GetCloseAmount() float64
	GetPositionCost() float64
	GetCloseProfit() float64
	GetPositionProfit() float64
	GetOpenCost() float64
	GetExchangeMargin() float64
	GetAccountMargin() float64
	GetBalance() float64
	SetOpenOrClose(byte)
	GetOpenOrClose() byte
}

type RespInfo interface {
	GetErrorID() int
	GetErrorMsg() string
}

type RespData struct {
	frontid        int
	sessionid      int
	maxorderRef    string
	instrumentid   string
	posidirection  byte
	ydposition     int
	position       int
	usemargin      float64
	openamount     float64
	closeamount    float64
	positioncost   float64
	closeprofit    float64
	positionprofit float64
	opencost       float64
	exchangemargin float64
	accountmargin  float64
	balance        float64
	openorclose    byte
}

func NewRespData() Response {
	return &RespData{}
}

func (resp *RespData) SetOpenOrClose(b byte) {
	resp.openorclose = b
}

func (resp *RespData) GetOpenOrClose() byte {
	return resp.openorclose
}

func (resp *RespData) SetFrontID(id int) {
	resp.frontid = id
}

func (resp *RespData) GetFrontID() int {
	return resp.frontid
}

func (resp *RespData) SetSessionID(id int) {
	resp.sessionid = id
}

func (resp *RespData) GetSessionID() int {
	return resp.sessionid
}

func (resp *RespData) SetMaxOrderRef(ref string) {
	resp.maxorderRef = ref
}

func (resp *RespData) GetMaxOrderRef() string {
	return resp.maxorderRef
}

func (resp *RespData) SetInstrumentID(ins string) {
	resp.instrumentid = ins
}

func (resp *RespData) GetInstrumentID() string {
	return resp.instrumentid
}

func (resp *RespData) SetPosiDirection(b byte) {
	resp.posidirection = b
}

func (resp *RespData) GetPosiDirection() byte {
	return resp.posidirection
}

func (resp *RespData) SetYdPosition(n int) {
	resp.ydposition = n
}

func (resp *RespData) GetYdPosition() int {
	return resp.ydposition
}

func (resp *RespData) SetPosition(n int) {
	resp.position = n
}

func (resp *RespData) GetPosition() int {
	return resp.position
}

func (resp *RespData) SetUseMargin(f float64) {
	resp.usemargin = f
}

func (resp *RespData) GetUseMargin() float64 {
	return resp.usemargin
}

func (resp *RespData) SetOpenAmount(f float64) {
	resp.openamount = f
}

func (resp *RespData) GetOpenAmount() float64 {
	return resp.openamount
}

func (resp *RespData) SetCloseAmount(f float64) {
	resp.closeamount = f
}

func (resp *RespData) GetCloseAmount() float64 {
	return resp.closeamount
}

func (resp *RespData) SetPositionCost(f float64) {
	resp.positioncost = f
}

func (resp *RespData) GetPositionCost() float64 {
	return resp.positioncost
}

func (resp *RespData) SetCloseProfit(f float64) {
	resp.closeprofit = f
}

func (resp *RespData) GetCloseProfit() float64 {
	return resp.closeprofit
}

func (resp *RespData) SetPositionProfit(f float64) {
	resp.positionprofit = f
}

func (resp *RespData) GetPositionProfit() float64 {
	return resp.positionprofit
}

func (resp *RespData) SetOpenCost(f float64) {
	resp.opencost = f
}

func (resp *RespData) GetOpenCost() float64 {
	return resp.opencost
}

func (resp *RespData) SetExchangeMargin(f float64) {
	resp.exchangemargin = f
}

func (resp *RespData) GetExchangeMargin() float64 {
	return resp.exchangemargin
}

func (resp *RespData) SetAccountMargin(f float64) {
	resp.accountmargin = f
}

func (resp *RespData) GetAccountMargin() float64 {
	return resp.accountmargin
}

func (resp *RespData) SetBalance(f float64) {
	resp.balance = f
}

func (resp *RespData) GetBalance() float64 {
	return resp.balance
}
