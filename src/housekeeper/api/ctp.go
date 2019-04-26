/*
 * @Description: CTP接口
 * @Author: gdd
 * @Date: 2018-11-07 09:21:29
 * @LastEditTime: 2018-11-30 14:24:56
 * @LastEditors: Please set LastEditors
 */

package api

import (
	"housekeeper/adaptation"
	"log"
	"time"

	"github.com/qerio/goctp"
)

type GoCTPClient struct {
	BrokerID        string
	InvestorID      string
	Password        string
	MdFront         string
	MdApi           goctp.CThostFtdcMdApi
	TraderFront     string
	TraderApi       goctp.CThostFtdcTraderApi
	MdRequestID     int
	TraderRequestID int
	login           chan bool
	connect         chan bool
}

func (g *GoCTPClient) StartTrader() bool {
	go func() {
		g.TraderApi.Join()
		g.TraderApi.Release()
	}()
	if <-g.connect {
		g.ReqUserLogin()
	} else {
		return false
	}
	if <-g.login {
		return true
	}
	return false
}

func (g *GoCTPClient) StartMd() {
	go func() {
		g.MdApi.Join()
		g.MdApi.Release()
	}()
}

func (g *GoCTPClient) GetMdRequestID() int {
	g.MdRequestID += 1
	if g.MdRequestID >= 999999 {
		g.MdRequestID = 0
	}
	return g.MdRequestID
}

func (g *GoCTPClient) GetTraderRequestID() int {
	g.TraderRequestID += 1
	if g.TraderRequestID >= 999999 {
		g.TraderRequestID = 0
	}
	return g.TraderRequestID
}

func NewDirectorCThostFtdcTraderSpi(v interface{}) goctp.CThostFtdcTraderSpi {
	return goctp.NewDirectorCThostFtdcTraderSpi(v)
}

type GoCThostFtdcMdSpi struct {
	client *GoCTPClient
	quote  adaptation.QuoteSpi
}

func (p *GoCThostFtdcMdSpi) OnRspError(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcMdSpi.OnRspError.")
	log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
}

func (p *GoCThostFtdcMdSpi) OnFrontDisconnected(nReason int) {
	log.Printf("GoCThostFtdcMdSpi.OnFrontDisconnected: %#v\n", nReason)
}

func (p *GoCThostFtdcMdSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Printf("GoCThostFtdcMdSpi.OnHeartBeatWarning: %v", nTimeLapse)
}

func (p *GoCThostFtdcMdSpi) OnFrontConnected() {
	log.Println("GoCThostFtdcMdSpi.OnFrontConnected.")
	p.client.ReqUserMdLogin()
}

func (p *GoCTPClient) ReqUserMdLogin() {
	log.Println("GoCThostFtdcMdSpi.ReqUserLogin.")

	req := goctp.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(p.BrokerID)
	req.SetUserID(p.InvestorID)
	req.SetPassword(p.Password)

	iResult := p.MdApi.ReqUserLogin(req, p.GetMdRequestID())

	if iResult != 0 {
		log.Println("发送用户登录请求: 失败.")
	} else {
		log.Println("发送用户登录请求: 成功.")
	}
}

func (p *GoCThostFtdcMdSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast {

		log.Printf("获取当前版本信息: %#v\n", goctp.CThostFtdcTraderApiGetApiVersion())
		log.Printf("获取当前交易日期: %#v\n", p.client.MdApi.GetTradingDay())
		log.Printf("获取用户登录信息: %#v %#v %#v\n", pRspUserLogin.GetLoginTime(), pRspUserLogin.GetSystemName(), pRspUserLogin.GetSessionID())

		//ppInstrumentID := []string{"cu1610", "cu1611", "cu1612", "cu1701", "cu1702", "cu1703", "cu1704", "cu1705", "cu1706"}
		//	ppInstrumentID := []string{"cu1810", "cu1811", "cu1812"}

		//	p.SubscribeMarketData(ppInstrumentID)
		//	p.SubscribeForQuoteRsp(ppInstrumentID)
		p.quote.OnRspUserLogin()
	}
}

func (p *GoCTPClient) SubscribeMarketData(name []string) {
	var tmp []string
	copy(tmp, name)
	iResult := p.MdApi.SubscribeMarketData(tmp)

	if iResult != 0 {
		log.Println("发送行情订阅请求: 失败.")
	} else {
		log.Println("发送行情订阅请求: 成功.")
	}
}

func (p *GoCThostFtdcMdSpi) OnRspSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspSubMarketData: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
}

func (p *GoCThostFtdcMdSpi) OnRtnDepthMarketData(pDepthMarketData goctp.CThostFtdcDepthMarketDataField) {

	p.quote.OnRtnDepthMarketData(&adaptation.Quote{
		TradingDay:         pDepthMarketData.GetTradingDay(),
		InstrumentID:       pDepthMarketData.GetInstrumentID(),
		ExchangeId:         pDepthMarketData.GetExchangeID(),
		ExchangeInstId:     pDepthMarketData.GetExchangeInstID(),
		LastPrice:          pDepthMarketData.GetLastPrice(),
		PreSettlementPrice: pDepthMarketData.GetPreSettlementPrice(),
		PreClosePrice:      pDepthMarketData.GetPreClosePrice(),
		PreOpenInterest:    pDepthMarketData.GetPreOpenInterest(),
		OpenPrice:          pDepthMarketData.GetOpenPrice(),
		HighestPrice:       pDepthMarketData.GetHighestPrice(),
		LowestPrice:        pDepthMarketData.GetLowestPrice(),
		Volume:             pDepthMarketData.GetVolume()})

}

type GoCThostFtdcTraderSpi struct {
	proc    adaptation.Proc
	login   chan bool
	connect chan bool
}

func (p *GoCThostFtdcTraderSpi) GetName() string {
	return "CTP"
}

func (p *GoCThostFtdcTraderSpi) SetProc(proc adaptation.Proc) {
	p.proc = proc
}

func (p *GoCThostFtdcTraderSpi) OnRspError(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspError.")
	p.IsErrorRspInfo(pRspInfo)
}

func (p *GoCThostFtdcTraderSpi) OnFrontDisconnected(nReason int) {
	log.Printf("GoCThostFtdcTraderSpi.OnFrontDisconnected: %#v", nReason)
	p.connect <- false
}

func (p *GoCThostFtdcTraderSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Printf("GoCThostFtdcTraderSpi.OnHeartBeatWarning: %#v", nTimeLapse)
}

func (p *GoCThostFtdcTraderSpi) OnFrontConnected() {
	log.Println("GoCThostFtdcTraderSpi.OnFrontConnected.")
	p.connect <- true
}

func (p *GoCTPClient) ReqUserLogin() {
	log.Println("GoCThostFtdcTraderSpi.ReqUserLogin.")

	req := goctp.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(p.BrokerID)
	req.SetUserID(p.InvestorID)
	req.SetPassword(p.Password)

	iResult := p.TraderApi.ReqUserLogin(req, p.GetTraderRequestID())

	if iResult != 0 {
		log.Println("发送用户登录请求: 失败.")
	} else {
		log.Println("发送用户登录请求: 成功.")
	}
}

func (p *GoCTPClient) IsFlowControl(iResult int) bool {
	return ((iResult == -2) || (iResult == -3))
}

func (p *GoCThostFtdcTraderSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {
	// 如果ErrorID != 0, 说明收到了错误的响应
	bResult := (pRspInfo.GetErrorID() != 0)
	if bResult {
		log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
	}
	return bResult
}

func (p *GoCThostFtdcTraderSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("GoCThostFtdcTraderSpi.OnRspUserLogin.")
	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		//log.Printf("获取当前交易日  : %#v\n", p.Client.TraderApi.GetTradingDay())
		log.Printf("获取用户登录信息: %#v %#v %#v\n", pRspUserLogin.GetFrontID(), pRspUserLogin.GetSessionID(), pRspUserLogin.GetMaxOrderRef())

		// // 保存会话参数
		// FRONT_ID = pRspUserLogin->FrontID;
		// SESSION_ID = pRspUserLogin->SessionID;
		// int iNextOrderRef = atoi(pRspUserLogin->MaxOrderRef);
		// iNextOrderRef++;
		// sprintf(ORDER_REF, "%d", iNextOrderRef);
		// sprintf(EXECORDER_REF, "%d", 1);
		// sprintf(FORQUOTE_REF, "%d", 1);
		// sprintf(QUOTE_REF, "%d", 1);
		// ///获取当前交易日
		// cerr << "获取当前交易日 = " << pMdApi->GetTradingDay() << endl;
		// ///投资者结算结果确认
		//p.ReqSettlementInfoConfirm()
		p.proc.SetSessionId(pRspUserLogin.GetSessionID())
		p.proc.SetFrontId(pRspUserLogin.GetFrontID())
		p.login <- true
	} else {
		p.login <- false
	}
}
func (p *GoCThostFtdcTraderSpi) OnRtnOrder(pOrder goctp.CThostFtdcOrderField) {

}

func (p *GoCThostFtdcTraderSpi) OnRtnTrade(pTrade goctp.CThostFtdcTradeField) {
	var pst adaptation.PositionField
	pst.InstrumentID = pTrade.GetInstrumentID()
	log.Println("return direction : ", pTrade.GetDirection(), "OpenOrClose", pTrade.GetOffsetFlag())
	if pTrade.GetOffsetFlag() == '0' {
		if pTrade.GetDirection() == '0' {
			pst.PosiDirection = '2'
		} else {
			pst.PosiDirection = '3'
		}
	} else {
		if pTrade.GetDirection() == '1' {
			pst.PosiDirection = '2'
		} else {
			pst.PosiDirection = '3'
		}
	}

	pst.Position = pTrade.GetVolume()
	pst.OpenOrClose = pTrade.GetOffsetFlag()
	p.proc.OnRtnTrade(&pst)
}

func (p *GoCThostFtdcTraderSpi) OnRspOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspOrderInsert.")
	log.Println("错误信息 : ", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
	if pRspInfo != nil && pRspInfo.GetErrorID() != 0 {
		var pst adaptation.PositionField
		pst.InstrumentID = pInputOrder.GetInstrumentID()
		if pInputOrder.GetDirection() == '1' {
			pst.PosiDirection = '2'
		} else {
			pst.PosiDirection = '3'
		}
		pst.Position = pInputOrder.GetVolumeTotalOriginal()
		pst.OpenOrClose = '1'
		p.proc.OnRspOrderInsert(&pst, false)
	}

}

func (p *GoCThostFtdcTraderSpi) OnErrRtnOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, pRspInfo goctp.CThostFtdcRspInfoField) {
	log.Println("GoCThostFtdcTraderSpi.OnErrRtnOrderInsert.")
	log.Println("错误信息 : ", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())

}

/*
func (p *GoCThostFtdcTraderSpi) ReqSettlementInfoConfirm() {
	req := goctp.NewCThostFtdcSettlementInfoConfirmField()

	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	iResult := p.Client.TraderApi.ReqSettlementInfoConfirm(req, p.Client.GetTraderRequestID())

	if iResult != 0 {
		log.Println("投资者结算结果确认: 失败.")
	} else {
		log.Println("投资者结算结果确认: 成功.")
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspSettlementInfoConfirm(pSettlementInfoConfirm goctp.CThostFtdcSettlementInfoConfirmField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	//cerr << "--->>> " << "OnRspSettlementInfoConfirm" << endl
	log.Println("GoCThostFtdcTraderSpi.OnRspSettlementInfoConfirm.")
	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {
		///请求查询合约
		p.ReqQryInstrument()
	}
}

func (p *GoCThostFtdcTraderSpi) ReqQryInstrument() {
	req := goctp.NewCThostFtdcQryInstrumentField()

	var id string = "cu1612"
	req.SetInstrumentID(id)

	for {
		iResult := p.Client.TraderApi.ReqQryInstrument(req, p.Client.GetTraderRequestID())

		if !p.IsFlowControl(iResult) {
			log.Printf("--->>> 请求查询合约: 成功 %#v\n", iResult)
			//break
		} else {
			log.Printf("--->>> 请求查询合约: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryInstrument(pInstrument goctp.CThostFtdcInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryInstrument: ", pInstrument.GetInstrumentID(), pInstrument.GetExchangeID(),
		pInstrument.GetInstrumentName(), pInstrument.GetExchangeInstID(), pInstrument.GetProductID(), pInstrument.GetProductClass(),
		pInstrument.GetDeliveryYear(), pInstrument.GetDeliveryMonth(), pInstrument.GetMaxMarketOrderVolume(), pInstrument.GetMinMarketOrderVolume(),
		pInstrument.GetMaxLimitOrderVolume(), pInstrument.GetMinLimitOrderVolume(), pInstrument.GetVolumeMultiple(), pInstrument.GetPriceTick(),
		pInstrument.GetCreateDate(), pInstrument.GetOpenDate(), pInstrument.GetExpireDate(), pInstrument.GetStartDelivDate(), pInstrument.GetEndDelivDate(),
		nRequestID, bIsLast)
	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		///请求查询合约
		p.ReqQryTradingAccount()
	}
}




*/
func (p *GoCTPClient) ReqQryInvestorPosition(id int) {

	req := goctp.NewCThostFtdcQryInvestorPositionField()
	req.SetBrokerID(p.BrokerID)
	req.SetInvestorID(p.InvestorID)

	for {
		iResult := p.TraderApi.ReqQryInvestorPosition(req, id)

		if !p.IsFlowControl(iResult) {
			log.Printf("--->>> 请求查询投资者持仓: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询投资者持仓: 受到流控 %#v\n", iResult)
			time.Sleep(10 * time.Second)
		}
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryInvestorPosition(pInvestorPosition goctp.CThostFtdcInvestorPositionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryInvestorPosition.")

	var pst adaptation.PositionField
	pst.InstrumentID = pInvestorPosition.GetInstrumentID()
	pst.PosiDirection = pInvestorPosition.GetPosiDirection()
	pst.YdPosition = pInvestorPosition.GetYdPosition()
	pst.Position = pInvestorPosition.GetPosition()
	p.proc.OnRspQryInvestorPosition(&pst, bIsLast)

}

func (p *GoCTPClient) ReqQryTradingAccount(id int) {
	req := goctp.NewCThostFtdcQryTradingAccountField()
	req.SetBrokerID(p.BrokerID)
	req.SetInvestorID(p.InvestorID)

	for {
		iResult := p.TraderApi.ReqQryTradingAccount(req, id)

		if !p.IsFlowControl(iResult) {
			log.Printf("--->>> 请求查询资金账户: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询资金账户: 受到流控 %#v\n", iResult)
		}
		time.Sleep(1 * time.Second)
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryTradingAccount(pTradingAccount goctp.CThostFtdcTradingAccountField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("GoCThostFtdcTraderSpi.OnRspQryTradingAccount.")
	log.Println("冻结的保证金", pTradingAccount.GetFrozenMargin(), "当前保证金", pTradingAccount.GetCurrMargin(),
		"可用资金", pTradingAccount.GetAvailable(), "交易所保证金", pTradingAccount.GetExchangeMargin(), "权益", pTradingAccount.GetBalance())
	log.Println("动态权益", pTradingAccount.GetAvailable()+pTradingAccount.GetCurrMargin())
	var ta adaptation.TradingAccount
	ta.CurrMargin = pTradingAccount.GetCurrMargin()
	ta.Balance = pTradingAccount.GetAvailable() + pTradingAccount.GetCurrMargin()
	p.proc.OnRspQryTradingAccount(&ta)

}

func (p *GoCTPClient) ReqOrderInsert(id int) {
	req := goctp.NewCThostFtdcInputOrderField()
	req.SetBrokerID(p.BrokerID)
	req.SetInvestorID(p.InvestorID)

}

func (p *GoCTPClient) ClearPosition() {

}

func (p *GoCTPClient) ClosePosition(pst adaptation.PositionField, id int) {
	log.Println("GoCThostFtdcTraderSpi.ClosePosition.")
	req := goctp.NewCThostFtdcInputOrderField()
	req.SetBrokerID(p.BrokerID)
	req.SetInvestorID(p.InvestorID)
	req.SetInstrumentID(pst.InstrumentID)
	req.SetUserID(p.InvestorID)
	req.SetOrderPriceType('2')
	if pst.PosiDirection == '2' {
		req.SetDirection('1')
		req.SetLimitPrice(pst.LowestPrice)
	} else {
		req.SetDirection('0')
		req.SetLimitPrice(pst.HighestPrice)
	}

	req.SetCombOffsetFlag("1")

	req.SetVolumeTotalOriginal(pst.Position)
	req.SetCombHedgeFlag("1")
	req.SetVolumeCondition('1')
	req.SetMinVolume(1)
	req.SetContingentCondition('1')
	req.SetForceCloseReason('0')
	req.SetIsAutoSuspend(0)
	req.SetUserForceClose(0)
	req.SetTimeCondition('3')
	p.TraderApi.ReqOrderInsert(req, id)

}

func NewCTPQuoteClient(BrokerID string,
	InvestorID string,
	Password string,
	MdFront string,
	TraderFront string,
	MdRequestID int,
	TraderRequestID int,
	quoteSpi adaptation.QuoteSpi) Client {

	CTP := GoCTPClient{BrokerID: BrokerID,
		InvestorID:      InvestorID,
		Password:        Password,
		MdFront:         MdFront,
		TraderFront:     TraderFront,
		MdRequestID:     MdRequestID,
		TraderRequestID: TraderRequestID}

	CTP.MdApi = goctp.CThostFtdcMdApiCreateFtdcMdApi()
	mdspi := &GoCThostFtdcMdSpi{quote: quoteSpi, client: &CTP}
	pMdSpi := goctp.NewDirectorCThostFtdcMdSpi(mdspi)

	CTP.MdApi.RegisterSpi(pMdSpi)
	CTP.MdApi.RegisterFront(CTP.MdFront)
	CTP.MdApi.Init()

	return &CTP

}

func NewCTPClient(BrokerID string,
	InvestorID string,
	Password string,
	MdFront string,
	TraderFront string,
	MdRequestID int,
	TraderRequestID int) Client {

	CTP := GoCTPClient{BrokerID: BrokerID,
		InvestorID:      InvestorID,
		Password:        Password,
		MdFront:         MdFront,
		TraderFront:     TraderFront,
		MdRequestID:     MdRequestID,
		TraderRequestID: TraderRequestID}

	CTP.TraderApi = goctp.CThostFtdcTraderApiCreateFtdcTraderApi()
	CTP.login = make(chan bool)
	CTP.connect = make(chan bool)
	var spiapi = &GoCThostFtdcTraderSpi{proc: nil, login: CTP.login, connect: CTP.connect}
	adaptation.RegistApi(spiapi)
	pTraderSpi := goctp.NewDirectorCThostFtdcTraderSpi(spiapi)

	CTP.TraderApi.RegisterSpi(pTraderSpi)                         // 注册事件类
	CTP.TraderApi.SubscribePublicTopic(0 /*THOST_TERT_RESTART*/)  // 注册公有流
	CTP.TraderApi.SubscribePrivateTopic(0 /*THOST_TERT_RESTART*/) // 注册私有流
	CTP.TraderApi.RegisterFront(CTP.TraderFront)
	CTP.TraderApi.Init()

	return &CTP

}
