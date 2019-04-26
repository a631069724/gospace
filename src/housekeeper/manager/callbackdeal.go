package manager

import (
	"housekeeper/adaptation"
	"log"
)

type Proc struct {
	apiname string
	acc     *Account
	wait    chan bool
}

type QuoteSpi struct {
	acc *Quoter
}

func (q *QuoteSpi) OnRspUserLogin() {
	//q.acc.SubscribeMarketData([]string{"c1901"})
}

func (q *QuoteSpi) OnRtnDepthMarketData(quote *adaptation.Quote) {
	log.Println("myRtnDepthMarketData")
	q.acc.SaveQuote(quote)
}

func (p *Proc) ApiName() string {
	return p.apiname
}

func NewProc(apiName string, acc *Account, w chan bool) *Proc {
	return &Proc{apiname: apiName, acc: acc, wait: w}
}

func (p *Proc) OnRspQryInvestorPosition(r *adaptation.PositionField, last bool) {
	log.Println("my OnRspQryInvestorPosition")

	p.acc.AddPosition(r)
	if last {
		p.wait <- true
	}

}

func (p *Proc) OnRspQryTradingAccount(r *adaptation.TradingAccount) {
	log.Println("my OnRspQryTradingAccount")
	p.acc.SetBalance(r.Balance)
	p.wait <- true
}

func (p *Proc) SetSessionId(id int) {
	p.acc.sessionid = id
}

func (p *Proc) SetFrontId(id int) {
	p.acc.frontid = id
}

func (p *Proc) OnRtnTrade(pst *adaptation.PositionField) {
	log.Println("my RtnTrade", pst)
	p.acc.SubPosition(pst)
}

func (p *Proc) OnRtnOrder() {
}

func (p *Proc) OnRspOrderInsert(pst *adaptation.PositionField, ret bool) {
	log.Println("my OnRspOrderInsert")
	if !ret {
		p.acc.SubReqPosition(pst)
	}

}
