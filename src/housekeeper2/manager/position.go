package manager

import (
	"fmt"
	"strconv"
)

type PostitionInfo struct {
	InstrumentID  string  `json:"instrumentid"`
	PosiDirection byte    `json:"direction,string"`
	YdPosition    int     `json:"ynum"`
	Position      int     `json:"snum"`
	UseMargin     float64 `json:"usemargin"`
}

type Position struct {
	contract    string
	direct      byte //0 表示空 1表示多
	num         int
	yestodayNum int
	useMargin   float64
	openAmount  float64
	closeProfit float64
	//markid      string
	//	marktime    string
}

/*
func (p *Position) SetMarkId(id string) {
	p.markid = id
}

func (p *Position) GetMarkId() string {
	return p.markid
}

func (p *Position) SetMarkTime(time string) {
	p.marktime = time
}

func (p *Position) GetMarkTime() string {
	return p.marktime
}*/

func (p Position) SubNum(tmp Position) Position {
	p.num = p.num - tmp.num
	p.yestodayNum = p.yestodayNum - tmp.yestodayNum
	return p
}

/*
func Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}*/

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func (p *Position) DirctIndex() int {
	if p.GetDirect() == 0x30 {
		return 0
	}
	return 1
}

func (p *Position) SetContract(contract string) {
	p.contract = contract
}

func (p *Position) GetContract() string {
	return p.contract
}

func (p *Position) SetDirect(dirct byte) {
	p.direct = dirct
}

func (p *Position) GetDirect() byte {
	return p.direct
}

func (p *Position) SetNum(num int) {
	if num < 0 {
		num = 0
	}
	p.num = num
}

func (p *Position) GetNum() int {
	return p.num
}

func (p *Position) SetYestodayNum(num int) {
	if num < 0 {
		num = 0
	}
	p.yestodayNum = num
}

func (p *Position) GetYestodayNum() int {
	return p.yestodayNum
}

func (p *Position) SetUseMargin(usemargin float64) {
	p.useMargin = Decimal(usemargin)
}

func (p *Position) GetUseMargin() float64 {
	return p.useMargin
}

func (p *Position) SetOpenAmount(openamount float64) {
	p.openAmount = Decimal(openamount)
}

func (p *Position) GetOpenAmount() float64 {
	return p.openAmount
}

func (p *Position) SetCloseProfit(closeprofit float64) {
	p.closeProfit = Decimal(closeprofit)
}

func (p *Position) GetCloseProfit() float64 {
	return p.closeProfit
}
