package manager

import (
	"housekeeper/adaptation"
	"housekeeper/api"
	"log"
	"sync"
	"time"
)

type Account struct {
	account      string
	password     string
	futureid     string
	apitype      string
	certinfo     string
	ip           string
	port         string
	rules        map[string]*Rule
	positions    map[string][2]adaptation.PositionField
	reqpositions map[string][2]adaptation.PositionField
	tradaccount  adaptation.TradingAccount
	sessionid    int
	frontid      int
	client       api.Client
	proc         *Proc
	quoter       *Quoter
	wait         chan bool
	reqbond      float64
	ruleMutex    sync.Mutex
	sync.Mutex
}

func (a *Account) SetQuoter(q *Quoter) {
	a.quoter = q
}

func NewAccount(accName string, pwd string, fid string, apitype string, certinfo string, ip string, port string) *Account {
	var client api.Client
	var proc *Proc

	wchan := make(chan bool)
	acc := &Account{account: accName, password: pwd, futureid: fid, apitype: apitype, certinfo: certinfo, ip: ip, port: port, wait: wchan,
		positions: make(map[string][2]adaptation.PositionField), reqpositions: make(map[string][2]adaptation.PositionField), rules: make(map[string]*Rule)}

	if apitype == "CTP" {
		client = api.NewCTPClient(certinfo, accName, pwd, "tcp://"+ip+":"+port, "tcp://"+ip+":"+port, 0, 0)
		proc = NewProc("CTP", acc, wchan)
		adaptation.RegisProc(proc)
	}

	acc.proc = proc
	acc.client = client
	log.Println("new account success")
	log.Println("client ", client, "start!")
	if !client.StartTrader() {
		return nil
	}
	return acc
}

func (a *Account) SubPosition(pst *adaptation.PositionField) {
	a.Lock()
	a.Unlock()
	if a.quoter.GetQuote(pst.InstrumentID) == nil {
		p := []string{pst.InstrumentID}
		a.quoter.SubscribeMarketData(p)
	}
	if pst.OpenOrClose == '0' {
		if v, ok := a.positions[pst.InstrumentID]; ok {
			if pst.PosiDirection == '2' {
				if v[1].InstrumentID == "" {
					v[1] = *pst
				} else {
					v[1].Position += pst.Position
				}

			} else {
				if v[0].InstrumentID == "" {
					v[0] = *pst
				} else {
					v[0].Position += pst.Position
				}
			}
			a.positions[pst.InstrumentID] = v
		} else {
			if pst.PosiDirection == '2' {
				v[1] = *pst
			} else {
				v[0] = *pst
			}
			a.positions[pst.InstrumentID] = v
		}

	} else {
		if v, ok := a.positions[pst.InstrumentID]; ok {
			if pst.PosiDirection == '2' {
				v[1].Position -= pst.Position
			} else {
				v[0].Position -= pst.Position
			}
			a.positions[pst.InstrumentID] = v
		}

		if v2, ok := a.reqpositions[pst.InstrumentID]; ok {

			if pst.PosiDirection == '2' {
				if v2[1].InstrumentID != "" {
					if v2[1].Position <= pst.Position {
						v2[1].InstrumentID = ""
						v2[1].Position = 0
					} else {
						v2[1].Position -= pst.Position
					}
				}
			} else {
				if v2[0].InstrumentID != "" {
					if v2[0].Position <= pst.Position {
						v2[0].InstrumentID = ""
						v2[0].Position = 0
					} else {
						v2[0].Position -= pst.Position
					}
				}
			}
			a.reqpositions[pst.InstrumentID] = v2
		}
	}
	log.Println(a.account, "当前持仓信息", a.positions)
	log.Println(a.account, "当前报单信息", a.reqpositions)
}

func (a *Account) SubReqPosition(pst *adaptation.PositionField) {
	a.Lock()
	defer a.Unlock()
	if v, ok := a.reqpositions[pst.InstrumentID]; ok {

		if pst.PosiDirection == '2' {
			if v[1].InstrumentID != "" {
				if v[1].Position <= pst.Position {
					v[1].InstrumentID = ""
					v[1].Position = 0
				} else {
					v[1].Position -= pst.Position
				}
			}
		} else {
			if v[0].InstrumentID != "" {
				if v[0].Position <= pst.Position {
					v[0].InstrumentID = ""
					v[0].Position = 0
				} else {
					v[0].Position -= pst.Position
				}
			}
		}
		a.reqpositions[pst.InstrumentID] = v
	}
	log.Println(a.account, "当前报单信息", a.reqpositions)
}

func (a *Account) AddRule(urle *Rule) {
	a.ruleMutex.Lock()
	defer a.ruleMutex.Unlock()
	log.Println("in addRule", urle)
	a.rules[urle.id] = urle
}

func (a *Account) MyInit() {
	a.client.ReqQryInvestorPosition(a.client.GetTraderRequestID())
	<-a.wait
	a.client.ReqQryTradingAccount(a.client.GetTraderRequestID())
	<-a.wait
}

func (a *Account) GetBalance() float64 {
	return a.tradaccount.Balance
}

func (a *Account) SetBalance(b float64) {
	a.tradaccount.Balance = b
	log.Println("当前权益 : ", a.tradaccount.Balance)
}

func (a *Account) ClearPosition() {
	a.Lock()
	defer a.Unlock()
	a.positions = make(map[string][2]adaptation.PositionField)
}

func (a *Account) AddPosition(pst *adaptation.PositionField) {
	a.Lock()
	defer a.Unlock()
	if a.quoter.GetQuote(pst.InstrumentID) == nil {
		p := []string{pst.InstrumentID}
		a.quoter.SubscribeMarketData(p)
	}

	mr := a.positions[pst.InstrumentID]
	if pst.PosiDirection == '2' {
		mr[1] = *pst
	} else {
		mr[0] = *pst
	}
	a.positions[pst.InstrumentID] = mr
	log.Println(a.account, "当前持仓信息", a.positions)
	log.Println(a.account, "当前报单信息", a.reqpositions)
}

func (a *Account) AddReqPosition(pst *adaptation.PositionField) {
	a.Lock()
	defer a.Unlock()
	mr := a.reqpositions[pst.InstrumentID]
	if pst.PosiDirection == '2' {
		if mr[1].InstrumentID == "" {
			mr[1] = *pst
		} else {
			mr[1].Position += pst.Position
		}

	} else {
		if mr[0].InstrumentID == "" {
			mr[0] = *pst
		} else {
			mr[0].Position += pst.Position
		}
	}
	a.reqpositions[pst.InstrumentID] = mr
	log.Println(a.account, "当前持仓信息", a.positions)
	log.Println(a.account, "当前报单信息", a.reqpositions)
}

func (a *Account) ClosePosition(pst adaptation.PositionField) {

	a.client.GetTraderRequestID()

}

func (a *Account) GetPostionNum(pst adaptation.PositionField) (int, int) {
	a.Lock()
	defer a.Unlock()
	if v, ok := a.positions[pst.InstrumentID]; ok {
		if reqv, ok := a.reqpositions[pst.InstrumentID]; ok {
			if pst.PosiDirection == '2' {
				return v[1].Position - reqv[1].Position, v[1].YdPosition - reqv[1].YdPosition
			} else {
				return v[0].Position - reqv[0].Position, v[0].YdPosition - reqv[0].YdPosition
			}
		} else {
			if pst.PosiDirection == '2' {
				return v[1].Position - reqv[1].Position, v[1].YdPosition - reqv[1].YdPosition
			} else {
				return v[0].Position - reqv[0].Position, v[0].YdPosition - reqv[0].YdPosition
			}
		}

	}
	return 0, 0
}

func (a *Account) GetCurrMargin() float64 {
	return a.tradaccount.CurrMargin
}

func (a *Account) GetPositions() map[string][2]adaptation.PositionField {
	a.Lock()
	defer a.Unlock()
	return a.positions
}

func (a *Account) GetContct(name string) Contract {
	a.Lock()
	defer a.Unlock()
	return a.quoter.GetContct(name)
}

func (a *Account) GetReqPositions() map[string][2]adaptation.PositionField {
	a.Lock()
	defer a.Unlock()
	return a.reqpositions
}

func (a *Account) CloseAllPosition() {
	psts := a.GetPositions()
	reqpsts := a.GetReqPositions()
	for _, v := range psts {
		if len(v) > 0 && v[0].InstrumentID != "" && v[0].Position > 0 {
			if reqv, ok := reqpsts[v[0].InstrumentID]; ok {
				v[0].Position -= reqv[0].Position
				v[0].YdPosition -= reqv[0].YdPosition
			}
			if v[0].Position > 0 {
				q := a.quoter.GetQuote(v[0].InstrumentID)
				if q != nil {
					v[0].HighestPrice = q.HighestPrice
					v[0].LowestPrice = q.LowestPrice
					a.client.ClosePosition(v[0], a.client.GetTraderRequestID())
					a.AddReqPosition(&v[0])
				}

			}

		}
		if len(v) > 0 && v[1].InstrumentID != "" && v[1].Position > 0 {
			if reqv, ok := reqpsts[v[1].InstrumentID]; ok {
				v[1].Position -= reqv[1].Position
				v[1].YdPosition -= reqv[1].YdPosition
			}
			if v[1].Position > 0 {
				q := a.quoter.GetQuote(v[1].InstrumentID)
				if q != nil {
					v[1].HighestPrice = q.HighestPrice
					v[1].LowestPrice = q.LowestPrice
					a.client.ClosePosition(v[1], a.client.GetTraderRequestID())
					a.AddReqPosition(&v[1])
				}

			}

		}
	}
}

func (a *Account) CheckRule() {
	a.ruleMutex.Lock()
	defer a.ruleMutex.Unlock()

	for _, v := range a.rules {
		ts, err := time.Parse("15:04:05", v.starttime)
		if err != nil {
			log.Println("star time", v.starttime, "err :", err.Error())
		}
		te, err := time.Parse("15:04:05", v.endtime)
		if err != nil {
			log.Println("end time", v.endtime, "err :", err.Error())
		}
		tn, _ := time.Parse("15:04:05", time.Now().Format("15:04:05"))
		if v.rtype == "01" {
			if tn.Before(te) && tn.After(ts) && a.GetBalance() < v.fundlevel {
				log.Println("当前权益 : ", a.GetBalance(), "符合规则 : ", v)
				a.CloseAllPosition()
			}
		} else if v.rtype == "02" {
			if tn.Before(te) && tn.After(ts) && a.GetCurrMargin() < v.bondlevel {
				log.Println("当前保证金", a.GetCurrMargin(), "符合规则", v)
				needbond := v.bondlevel - a.GetCurrMargin()
				a.CloseBondPosition(needbond - a.GetReqBond())
			}
		}
	}

}

func (a *Account) GetReqBond() float64 {
	var bond float64
	reqpst := a.GetReqPositions()
	for _, v := range reqpst {

		if len(v) > 0 && v[0].InstrumentID != "" && v[0].Position > 0 {
			q := *a.quoter.GetQuote(v[0].InstrumentID)
			contract := a.GetContct(v[0].InstrumentID)
			bond += q.LastPrice * contract.money * 0.1 * float64(v[0].Position)
			if v[1].InstrumentID != "" && v[1].Position > 0 {
				contract := a.GetContct(v[1].InstrumentID)
				bond += q.LastPrice * contract.money * 0.1 * float64(v[1].Position)
			}
		}

	}
	return bond

}

func (a *Account) CloseBondPosition(bond float64) {
	psts := a.GetPositions()
	reqpsts := a.GetReqPositions()

	for _, v := range psts {
		if bond <= 0 {
			return
		}
		if len(v) > 0 && v[0].InstrumentID != "" && v[0].Position > 0 {

			if reqv, ok := reqpsts[v[0].InstrumentID]; ok {
				v[0].Position -= reqv[0].Position
				v[0].YdPosition -= reqv[0].YdPosition
			}
			if v[0].Position > 0 {
				q := a.quoter.GetQuote(v[0].InstrumentID)
				contract := a.GetContct(v[0].InstrumentID)
				if q != nil && contract.money > 0 {
					InstrumentIDBond := q.LastPrice * contract.money * 0.1
					if bond > InstrumentIDBond*float64(v[0].Position) {
						bond -= InstrumentIDBond * float64(v[0].Position)
					} else {
						num := (int(bond) % int((InstrumentIDBond * float64(v[0].Position)))) + 1
						bond = 0
						v[0].Position = num
					}
					v[0].HighestPrice = q.HighestPrice
					v[0].LowestPrice = q.LowestPrice
					a.client.ClosePosition(v[0], a.client.GetTraderRequestID())
					a.AddReqPosition(&v[0])
				}

			}

		}
		if len(v) > 0 && v[1].InstrumentID != "" && v[1].Position > 0 {
			if reqv, ok := reqpsts[v[1].InstrumentID]; ok {
				v[1].Position -= reqv[1].Position
				v[1].YdPosition -= reqv[1].YdPosition
			}
			if v[1].Position > 0 {
				contract := a.GetContct(v[1].InstrumentID)
				q := a.quoter.GetQuote(v[1].InstrumentID)
				if q != nil && contract.money > 0 {
					InstrumentIDBond := q.LastPrice * contract.money * 0.1
					if bond > InstrumentIDBond*float64(v[1].Position) {
						bond -= InstrumentIDBond * float64(v[1].Position)
					} else {
						num := (int(bond) % int((InstrumentIDBond * float64(v[1].Position)))) + 1
						bond = 0
						v[1].Position = num
					}
					v[1].HighestPrice = q.HighestPrice
					v[1].LowestPrice = q.LowestPrice
					a.client.ClosePosition(v[1], a.client.GetTraderRequestID())
					a.AddReqPosition(&v[1])
				}

			}

		}
	}
}
