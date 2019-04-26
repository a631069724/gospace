package manager

import (
	"database/sql"
	"encoding/json"
	"errors"
	"housekeeper2/Tools"
	"housekeeper2/api"
	"housekeeper2/rule"
	"log"
	"math"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type AccountInfo struct {
	Account    string          `json:"account"`
	Balance    float64         `json:"balance"`
	CurrMargin float64         `json:"margin"`
	Positions  []PostitionInfo `json:"positions"`
}

func (a *Account) GetAccountInfo() string {
	var acin AccountInfo
	acin.Account = a.Name()
	acin.Balance = a.GetBalance()
	acin.CurrMargin = a.GetMargin()
	for _, v := range a.positions {

		if v[0].GetContract() != "" && v[0].GetNum() > 0 {
			var pst PostitionInfo
			pst.InstrumentID = v[0].GetContract()
			pst.PosiDirection = v[0].GetDirect()
			pst.Position = v[0].GetNum()
			pst.YdPosition = v[0].GetYestodayNum()
			pst.UseMargin = v[0].GetUseMargin()
			acin.Positions = append(acin.Positions, pst)
		}
		if v[1].GetContract() != "" && v[1].GetNum() > 0 {
			var pst PostitionInfo
			pst.InstrumentID = v[1].GetContract()
			pst.PosiDirection = v[1].GetDirect()
			pst.Position = v[1].GetNum()
			pst.YdPosition = v[1].GetYestodayNum()
			pst.UseMargin = v[1].GetUseMargin()
			acin.Positions = append(acin.Positions, pst)
		}
	}
	jb, err := json.Marshal(acin)
	if err != nil {
		log.Println("json marshal err: ", err)
		return ""
	}
	log.Println(string(jb))
	return string(jb)
}

type Account struct {
	isused        bool
	isexpir       bool
	running       bool
	relogin       bool
	certinfo      string
	accName       string
	pwd           string
	ip            string
	port          string
	expriration   string
	_type         string
	reqid         int
	balance       float64
	margin        float64
	traderapi     api.TraderApi
	positions     map[string][2]Position
	inspositions  map[string][2]Position
	result        chan error
	params        rule.Params
	rulers        map[string]rule.Ruler
	priofund      float64
	afterfund     float64
	bmult         float64
	forceclose    float64
	premin        string
	udpect        float64
	contractTable ContractTable
	qtclient      QuoterClient
	priorclose    byte
	db            *sql.DB
	sync.Mutex
}

func (a *Account) SetDber(db *sql.DB) {
	a.db = db
}

func (a *Account) SetQuoteClient(q QuoterClient) {
	a.qtclient = q
}

func (a *Account) DelRuler(_type string) {
	a.Lock()
	defer a.Unlock()
	delete(a.rulers, _type)
}

func (a *Account) AddRuler(_type string, r rule.Ruler) {
	a.Lock()
	defer a.Unlock()
	a.rulers[_type] = r
}

func (a *Account) SetBalance(b float64) {
	a.balance = Decimal(b)
}

func (a *Account) SetMargin(m float64) {
	a.margin = Decimal(m)
}

func (a *Account) ShowPositions() {
	log.Println(a.positions)
}

func (a *Account) ShowAccount() {
	log.Println(a)
}

func (a *Account) AppendPositions(pst Position) {
	a.Lock()
	defer a.Unlock()
	contract := pst.GetContract()
	tmp := a.positions[contract]
	tmp[pst.DirctIndex()] = pst
	a.positions[contract] = tmp
}

func (a *Account) SetReLogin(b bool) {
	a.relogin = b
}

func (a *Account) ReLogin() bool {
	return a.relogin
}

func (a *Account) AccType() string {
	return a._type
}

func (a *Account) GetRequestID() int {
	a.reqid += 1
	if a.reqid >= 999999 {
		a.reqid = 0
	}
	return a.reqid
}

func (a *Account) Name() string {
	return a.accName
}

func (a *Account) Pwd() string {
	return a.pwd
}

func (a *Account) Certinfo() string {
	return a.certinfo
}

func (a *Account) Result() error {
	select {
	case b := <-a.result:
		return b
	case <-time.After(5 * time.Second):
		return errors.New("over time")
	}
}

func (a *Account) Return(err error) bool {
	select {
	case <-time.After(5 * time.Second):
		return false
	case a.result <- err:
		return true
	}

}

func NewAccount(used bool, certinfo string, accName string, pwd string, ip string, port string, expriration string, _type string,
	priofund float64, afterfund float64, bmult float64, udpect float64, permin string, fclose float64, ct ContractTable,
	closeflag byte) (*Account, error) {
	ctpapi, err := api.NewTraderApi(_type)
	if err != nil {
		return nil, err
	}
	return &Account{
		isused:        used,
		ip:            ip,
		port:          port,
		_type:         _type,
		traderapi:     ctpapi,
		expriration:   expriration,
		certinfo:      certinfo,
		accName:       accName,
		pwd:           pwd,
		result:        make(chan error),
		positions:     make(map[string][2]Position),
		inspositions:  make(map[string][2]Position),
		rulers:        make(map[string]rule.Ruler),
		params:        rule.NewParams(),
		priofund:      priofund,
		afterfund:     afterfund,
		bmult:         bmult,
		udpect:        udpect,
		premin:        permin,
		forceclose:    fclose,
		priorclose:    closeflag,
		contractTable: ct,
	}, nil

}

func (a *Account) RegistSpi(spi api.TraderSpi) error {
	if spi == nil {
		return errors.New("The spi is NULL")
	}
	a.traderapi.RegistSpi(api.NewTraderSpi(spi, a.AccType()))
	return nil
}

func (a *Account) Connect() error {
	a.traderapi.Connect(a.ip, a.port)
	return a.Result()
}

func (a *Account) OnConnect() {
	log.Println("account.OnConnect")
	if a.ReLogin() {
		a.ReqUserLogin()
		return
	}
	a.Return(nil)

}

func (a *Account) ReqUserLogin() error {
	reqdata := api.NewReqData()
	reqdata.SetName(a.Name())
	reqdata.SetPwd(a.Pwd())
	reqdata.SetId(a.GetRequestID())
	reqdata.SetBrokeId(a.Certinfo())
	a.traderapi.ReqUserLogin(reqdata)
	if a.ReLogin() {
		return nil
	}
	return a.Result()

}

func (a *Account) OnRspUserLogin(resp api.Response, info api.RespInfo, reqid int, islast bool) {
	log.Println("account.OnRspUserLogin")
	a.Return(nil)
}

func (a *Account) ReqUserLogout() {
	log.Println("account.ReqUserLogout")
	reqdata := api.NewReqData()
	reqdata.SetName(a.Name())
	reqdata.SetId(a.GetRequestID())
	reqdata.SetBrokeId(a.Certinfo())
	a.traderapi.ReqUserLogout(reqdata)
}

func (a *Account) OnRspUserLogout() {
	log.Println("account.OnRspUserLogout")
	a.isused = false
}

func (a *Account) ReqQryInvestorPosition() {
	reqdata := api.NewReqData()
	reqdata.SetName(a.Name())
	reqdata.SetBrokeId(a.Certinfo())
	reqdata.SetInstrumentID("")
	reqdata.SetId(a.GetRequestID())
	a.traderapi.ReqQryInvestorPosition(reqdata)
	a.Result()
}

func (a *Account) OnRspQryInvestorPosition(resp api.Response, info api.RespInfo, reqid int, islast bool) {
	log.Println("account.OnRspQryInvestorPosition")

	var pst Position
	instrumentid := resp.GetInstrumentID()
	pst.SetCloseProfit(resp.GetCloseProfit())
	pst.SetDirect(resp.GetPosiDirection())
	pst.SetContract(instrumentid)
	pst.SetNum(resp.GetPosition())
	pst.SetYestodayNum(resp.GetYdPosition())
	pst.SetOpenAmount(resp.GetOpenAmount())
	pst.SetUseMargin(resp.GetUseMargin())
	cnt := a.contractTable.GetContract(instrumentid)
	cnt.SetMargin(instrumentid, resp.GetUseMargin()/float64(resp.GetPosition()))
	a.AppendPositions(pst)
	p := []string{instrumentid}
	if !a.qtclient.QuoterExist(instrumentid) {
		a.qtclient.SubscribeMarketData(p)
	}
	if islast {
		a.Return(nil)
	}
}

func (a *Account) ReqQryTradingAccount() {
	reqdata := api.NewReqData()
	reqdata.SetBrokeId(a.Certinfo())
	reqdata.SetName(a.Name())
	reqdata.SetCur("CNY")
	reqdata.SetId(a.GetRequestID())
	a.traderapi.ReqQryTradingAccount(reqdata)
	a.Result()
}

func (a *Account) OnRspQryTradingAccount(resp api.Response, info api.RespInfo, reqid int, islast bool) {
	log.Println("account.OnRspQryTradingAccount")
	a.SetBalance(resp.GetBalance())
	a.SetMargin(resp.GetAccountMargin())
	a.Return(nil)
}

func (a *Account) ReqOrderInsert(reqorder Order) {
	reqdata := api.NewReqData()
	reqdata.SetBrokeId(a.Certinfo())
	reqdata.SetName(a.Name())
	reqdata.SetExchangeId(reqorder.GetCombOffsetFlag())
	reqdata.SetInstrumentID(reqorder.GetInstrumentId())
	reqdata.SetOrderPriceType(reqorder.GetOrderPriceType())
	reqdata.SetCombOffSetFlag(reqorder.GetCombOffsetFlag())
	reqdata.SetCombHedgeFlag(reqorder.GetCombHedgeFlag())
	reqdata.SetId(a.GetRequestID())
	reqdata.SetDirection(reqorder.GetDirect())
	reqdata.SetNum(reqorder.GetNum())
	reqdata.SetLimitPrice(reqorder.GetLimitPrice())
	a.traderapi.ReqOrderInsert(reqdata)
}

func (a *Account) Run() error {
	a.RegistSpi(a)
	if err := a.Connect(); err != nil {
		log.Println("acount: ", a.Name(), "connect error: ", err.Error())
		return err
	}
	if err := a.ReqUserLogin(); err != nil {
		log.Println("acount: ", a.Name(), "login error: ", err.Error())
		return err
	}
	a.SetReLogin(true)
	//a.ReqQryInvestorPosition()
	//a.ReqQryTradingAccount()
	a.InitParams()
	//a.ShowPositions()
	//a.ShowAccount()
	for a.IsUsed() && a.IsExpir() {
		a.ReqQryInvestorPosition()
		a.ReqQryTradingAccount()
		log.Println("begin check rule")
		a.CheckRule()
		time.Sleep(30 * time.Second)
	}
	return nil
}

func (a *Account) Exit() {
	a.SetReLogin(false)
	a.ReqUserLogout()
}

func (a *Account) SetUsed(b bool) {
	a.isused = b
}

func (a *Account) IsUsed() bool {
	return a.isused
}

func (a *Account) SetExpir(time string) {
	a.expriration = time
}

func (a *Account) IsExpir() bool {
	te, err := time.Parse("20060102150405", a.expriration)
	if err != nil {
		log.Println("exprirtion: ", te, "err :", err.Error())
		return false
	}
	tn, _ := time.Parse("20060102150405", time.Now().Format("20060102150405"))
	log.Println(te, tn)
	return tn.Before(te)
}

func (a *Account) GetRulers() map[string]rule.Ruler {
	return a.rulers
}

func (a *Account) InitParams() {
	a.params.SetParam("Balance", a.GetBalance())
	a.params.SetParam("ForceClose", a.GetForceClose())
	a.params.SetParam("Margin", a.GetMargin())
	a.params.SetParam("PrioFund", a.GetPrioritypFund())
	a.params.SetParam("BondMult", a.GetBondMult())
	a.params.SetParam("PreMin", a.GetBondMult())
	a.params.SetParam("UDPercent", a.GetBondMult())
	a.params.SetParam("ContractTable", a.contractTable)
}

func (a *Account) SetUDPercent(f float64) {
	a.udpect = f
}

func (a *Account) GetUDPercent() float64 {
	return a.udpect
}

func (a *Account) GetBalance() float64 {
	return a.balance
}

func (a *Account) SetForceClose(f float64) {
	a.forceclose = f
}

func (a *Account) GetForceClose() float64 {
	return a.forceclose
}

func (a *Account) GetMargin() float64 {
	return a.margin
}

func (a *Account) SetPrioritypFund(f float64) {
	a.priofund = f
}

func (a *Account) GetPrioritypFund() float64 {
	return a.priofund
}

func (a *Account) SetBondMult(f float64) {
	a.bmult = f
}

func (a *Account) GetBondMult() float64 {
	return a.bmult
}

func (a *Account) SetPreMin(s string) {
	a.premin = s
}

func (a *Account) GetPreMin() string {
	return a.premin
}

func (a *Account) GetParams() rule.Params {
	return a.params
}

func (a *Account) GetPriorClose() byte {
	return a.priorclose
}

func (a *Account) GetValidPositions() map[string][2]Position {
	var mpst map[string][2]Position
	mpst = make(map[string][2]Position)
	for k, v := range a.positions {
		var pst [2]Position
		tmp := a.inspositions[k]
		pst[0] = v[0].SubNum(tmp[0])
		pst[1] = v[1].SubNum(tmp[1])
		mpst[k] = pst
	}
	return mpst

}

func (a *Account) InsReqTable(order Order) {
	pst := a.inspositions[order.GetInstrumentId()]
	index := order.DirctIndex()
	if order.GetYdOrTd() == '1' {
		pst[index].SetYestodayNum(pst[index].GetYestodayNum() + order.GetNum())
	}
	pst[index].SetNum(pst[index].GetNum() + order.GetNum())
}

func (a *Account) RespModifyTable(respOrder RespOrder) {
	pst := a.inspositions[respOrder.GetInstrumentID()]
	index := 0
	if respOrder.GetPosiDirection() == '1' {
		index = 1
	}
	ynum := pst[index].GetYestodayNum()
	tnum := pst[index].GetNum()
	num := respOrder.GetPosition()
	flag := respOrder.GetOpenOrClose()
	if flag == '1' || flag == '3' || flag == '4' { //平仓
		pst[index].SetNum(tnum - num)
	}
	if respOrder.GetOpenOrClose() == '3' { //平昨
		pst[index].SetYestodayNum(ynum - num)
	}

}

func (a *Account) ClosePostion(flag byte, p *Position, idstr string, timestr string) {
	if p.GetNum() > 0 {
		insOrder := NewOrder(flag)
		insOrder.ConvPosition(p)
		if insOrder.GetNum() <= 0 {
			return
		}
		insOrder.SetOrderPriceType('2') //3代表最优价 1代表任意价 2代表限价单 4代表最新价
		if flag == 1 {
			insOrder.SetCombOffsetFlag("4") //4 代表平zuo仓
		} else if flag == 2 {
			insOrder.SetCombOffsetFlag("3") //3 代表平jin仓
		} else {
			insOrder.SetCombOffsetFlag("1") //1 代表平仓
		}

		insOrder.SetCombHedgeFlag("1") //1 代表套保
		quoter := a.qtclient.GetQuoter(p.GetContract())
		if quoter == nil {
			return
		}
		if p.GetDirect() == '0' {
			insOrder.SetLimitPrice(quoter.GetHighestPrice())
		} else if p.GetDirect() == '1' {
			insOrder.SetLimitPrice(quoter.GetLowestPrice())
		}
		a.SaveReqOrderInsert(insOrder, idstr, timestr)
		a.ReqOrderInsert(insOrder)
		a.InsReqTable(insOrder)
	}
}

func (a *Account) SaveReqOrderInsert(od Order, idstr string, timestr string) {
	dblocker.Lock()
	defer dblocker.Unlock()
	a.db.Exec(`insert into ZHGJ_REQ_ORDER_LS (ID,TIME,CONTRACT,COMBOFFSETFLAG,
		NUM,DIRECT,LIMITPRICE) values ('%s','%s','%s',
		'%s',%d,'%s',%f)`,
		idstr, timestr, od.GetInstrumentId(), od.GetCombOffsetFlag(),
		od.GetNum(), od.GetDirect(), od.GetLimitPrice())
}

func (a *Account) CloseAllPositions(idstr string, timestr string) {
	psts := a.GetValidPositions()
	//flag := a.GetPriorClose()
	for name, pst := range psts {
		contract := a.contractTable.GetContract(name)
		if contract.IsSHFE() {
			a.ClosePostion(1, &pst[0], idstr, timestr)
			a.ClosePostion(2, &pst[0], idstr, timestr)
			a.ClosePostion(1, &pst[1], idstr, timestr)
			a.ClosePostion(2, &pst[1], idstr, timestr)
		} else {
			a.ClosePostion(0, &pst[0], idstr, timestr)
			a.ClosePostion(0, &pst[1], idstr, timestr)
		}

	}

}

//func GetCloseBSNum(bond float64)

func (a *Account) CloseBondPositions(bond float64, idstr string, timestr string) {
	psts := a.GetValidPositions()
	log.Println(psts)
	for name, pst := range psts {
		if bond <= 0 {
			return
		}

		contract := a.contractTable.GetContract(name)
		if Tools.InTheTimeMins(a.GetPreMin(), contract.GetCloseTime()) {
			useMargin := contract.GetMargin(name)
			psellnum := pst[0].GetNum()
			pbuynum := pst[1].GetNum()
			pmaxnum := psellnum
			pflag := 's'
			if psellnum < pbuynum {
				pflag = 'b'
				pmaxnum = pbuynum
			}
			if float64(pmaxnum)*useMargin < bond {
				if contract.IsSHFE() {
					a.ClosePostion(1, &pst[0], idstr, timestr)
					a.ClosePostion(2, &pst[0], idstr, timestr)
					a.ClosePostion(1, &pst[1], idstr, timestr)
					a.ClosePostion(2, &pst[1], idstr, timestr)
				} else {
					a.ClosePostion(0, &pst[0], idstr, timestr)
					a.ClosePostion(0, &pst[1], idstr, timestr)
				}

				bond = bond - float64(pmaxnum)*useMargin
			} else {
				var ponum int
				mm := bond / useMargin
				pnum := int(math.Ceil(mm))
				bond = 0
				if pflag == 's' {
					ponum = pnum + pst[1].GetNum() - pst[0].GetNum()
					if contract.IsSHFE() {
						if pnum > pst[0].GetYestodayNum() {
							a.ClosePostion(1, &pst[0], idstr, timestr)
							pst[0].SetNum(pnum)
							a.ClosePostion(2, &pst[0], idstr, timestr)
						} else {
							pst[0].SetYestodayNum(pnum)
							a.ClosePostion(1, &pst[0], idstr, timestr)
						}
						if ponum > pst[1].GetYestodayNum() {
							a.ClosePostion(1, &pst[1], idstr, timestr)
							pst[1].SetNum(ponum)
							a.ClosePostion(2, &pst[1], idstr, timestr)
						} else {
							pst[1].SetYestodayNum(ponum)
							a.ClosePostion(1, &pst[1], idstr, timestr)
						}
					} else {
						pst[0].SetNum(pnum)
						a.ClosePostion(0, &pst[0], idstr, timestr)
						pst[1].SetNum(ponum)
						a.ClosePostion(0, &pst[1], idstr, timestr)
					}
				} else {
					ponum = pnum + pst[0].GetNum() - pst[1].GetNum()
					if contract.IsSHFE() {
						if pnum > pst[1].GetYestodayNum() {
							a.ClosePostion(1, &pst[1], idstr, timestr)
							pst[1].SetNum(pnum)
							a.ClosePostion(2, &pst[1], idstr, timestr)
						} else {
							pst[1].SetYestodayNum(pnum)
							a.ClosePostion(1, &pst[1], idstr, timestr)
						}
						if ponum > pst[0].GetYestodayNum() {
							a.ClosePostion(1, &pst[0], idstr, timestr)
							pst[0].SetNum(ponum)
							a.ClosePostion(2, &pst[0], idstr, timestr)
						} else {
							pst[0].SetYestodayNum(ponum)
							a.ClosePostion(1, &pst[0], idstr, timestr)
						}
					} else {
						pst[1].SetNum(pnum)
						a.ClosePostion(0, &pst[1], idstr, timestr)
						pst[0].SetNum(ponum)
						a.ClosePostion(0, &pst[0], idstr, timestr)
					}
				}

			}
		}
	}
}

func (a *Account) CheckRule() {
	a.Lock()
	defer a.Unlock()
	rulers := a.GetRulers()
	params := a.GetParams()

	params.SetParam("Balance", a.GetBalance())
	params.SetParam("ForceClose", a.GetForceClose())
	params.SetParam("Margin", a.GetMargin())
	params.SetParam("Positions", a.GetValidPositions())
	params.SetParam("QuoterCln", a.qtclient)
	params.SetParam("UDPercent", a.GetUDPercent())
	for rtype, ruler := range rulers {
		ret, err := ruler.Judge(params)
		if err != nil {
			log.Println("judge type: ", rtype, "result: ", err.Error())
			continue
		}
		idstr, timestr := a.SaveRule(rtype)
		if rtype == "01" && ret.(bool) {
			log.Println("rule type: ", rtype, "judge return: ", ret.(bool))
			a.CloseAllPositions(idstr, timestr)
		}
		if rtype == "02" && ret.(float64) > 0 {
			log.Println("rule type: ", rtype, "judge return: ", ret.(float64))
			a.CloseBondPositions(ret.(float64), idstr, timestr)
		}
		if rtype == "03" && err == nil {
			log.Println("rule type: ", rtype, "judge return: ", ret.([]Position))
			for _, pst := range ret.([]Position) {
				contract := a.contractTable.GetContract(pst.GetContract())
				if contract.IsSHFE() {
					a.ClosePostion(1, &pst, idstr, timestr)
					a.ClosePostion(2, &pst, idstr, timestr)

				} else {
					a.ClosePostion(0, &pst, idstr, timestr)
				}
			}
		}
	}
}

func (a *Account) SaveRule(ttype string) (string, string) {
	dblocker.Lock()
	defer dblocker.Unlock()
	id, _ := uuid.NewV4()
	idstr := id.String()
	timestr := time.Now().Format("2006-01-02 03:04:05")
	a.db.Exec(`insert into ZHGJ_RULE_LS (ID,TIME,RULE_TYPE,
		BALANCE,FORCECLOSE,MARGIN,PRIOFUND,
		BONDMULT,UDPERCENT) VALUES ('%s','%s','%s',%f,%f,
	%f,%f,%f)`,
		idstr, timestr, ttype,
		a.params.GetParam("Balance").(float64),
		a.params.GetParam("ForceClose").(float64),
		a.params.GetParam("Margin").(float64),
		a.params.GetParam("PrioFund").(float64),
		a.params.GetParam("BondMult").(float64),
		a.params.GetParam("UDPercent").(float64))
	return idstr, timestr
}

func (a *Account) OnRtnTrade(resp api.Response) {
	a.Lock()
	defer a.Unlock()
	log.Println("account.OnRtnTrade")
	a.RespModifyTable(resp)

}
