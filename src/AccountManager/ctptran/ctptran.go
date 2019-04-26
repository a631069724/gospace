package ctptran

import (
	"log"
	"time"

	"github.com/qerio/goctp"
	"gopkg.in/ini.v1"
)

type CtpTran struct {
	TraderClient GoCTPClient
	MdClient     GoCTPClient
	traderlogin  chan bool
	mdlogin      chan bool
}

type GoCTPClient struct {
	BrokerID   string `ini:"CTP_BrokerId"`
	InvestorID string `ini:"CTP_InvertorId"`
	Password   string `ini:"CTP_Password"`

	MdFront string `ini:"CTP_MdFront"`
	MdApi   goctp.CThostFtdcMdApi

	TraderFront string `ini:"CTP_TraderFront"`
	TraderApi   goctp.CThostFtdcTraderApi

	MdRequestID     int `ini:"CTP_MdRequestId`
	TraderRequestID int `ini:"CTP_TraderRequestId'`
	login           chan bool
}

func (c *CtpTran) Start() bool {
	/*	go func() {
		c.MdClient.MdApi.Join()
		c.MdClient.MdApi.Release()
	}()*/
	go func() {
		c.TraderClient.TraderApi.Join()
		c.TraderClient.TraderApi.Release()
	}()
	defer log.Println("start ok")
	//return (<-c.TraderClient.login && <-c.MdClient.login)
	return <-c.TraderClient.login
}

func (c *CtpTran) SubscribeMarketData(name []string) {

	iResult := c.MdClient.MdApi.SubscribeMarketData(name)

	if iResult != 0 {
		log.Println("发送行情订阅请求: 失败.")
	} else {
		log.Println("发送行情订阅请求: 成功.")
	}
}

func (c *CtpTran) ReqQryInvestorPosition() {
	req := goctp.NewCThostFtdcQryInvestorPositionField()
	req.SetBrokerID(c.TraderClient.BrokerID)
	req.SetInvestorID(c.TraderClient.InvestorID)

	iResult := c.TraderClient.TraderApi.ReqQryInvestorPosition(req, c.TraderClient.GetTraderRequestID())

	if !c.IsFlowControl(iResult) {
		log.Printf("--->>> 请求查询投资者持仓: 成功 %#v\n", iResult)
		//break;
	} else {
		log.Printf("--->>> 请求查询投资者持仓: 受到流控 %#v\n", iResult)

	}

}
func (g *GoCTPClient) GetMdRequestID() int {
	g.MdRequestID += 1
	return g.MdRequestID
}

func (g *GoCTPClient) GetTraderRequestID() int {
	g.TraderRequestID += 1
	return g.TraderRequestID
}

func NewDirectorCThostFtdcTraderSpi(v interface{}) goctp.CThostFtdcTraderSpi {
	return goctp.NewDirectorCThostFtdcTraderSpi(v)
}

type GoCThostFtdcTraderSpi struct {
	Client GoCTPClient
}

func (p *GoCThostFtdcTraderSpi) OnRspError(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("GoCThostFtdcTraderSpi.OnRspError.")
	p.IsErrorRspInfo(pRspInfo)
}

func (p *GoCThostFtdcTraderSpi) OnFrontDisconnected(nReason int) {
	log.Printf("GoCThostFtdcTraderSpi.OnFrontDisconnected: %#v", nReason)
}

func (p *GoCThostFtdcTraderSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Printf("GoCThostFtdcTraderSpi.OnHeartBeatWarning: %#v", nTimeLapse)
}

func (p *GoCThostFtdcTraderSpi) OnFrontConnected() {
	log.Println("GoCThostFtdcTraderSpi.OnFrontConnected.")
	p.ReqUserLogin()
}

func (p *GoCThostFtdcTraderSpi) ReqUserLogin() {
	log.Println("GoCThostFtdcTraderSpi.ReqUserLogin.")

	req := goctp.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetUserID(p.Client.InvestorID)
	req.SetPassword(p.Client.Password)

	iResult := p.Client.TraderApi.ReqUserLogin(req, p.Client.GetTraderRequestID())

	if iResult != 0 {
		log.Println("发送用户登录请求: 失败.")
	} else {
		log.Println("发送用户登录请求: 成功.")
	}
}

func (p *CtpTran) IsFlowControl(iResult int) bool {
	return ((iResult == -2) || (iResult == -3))
}

func (p *GoCThostFtdcTraderSpi) IsFlowControl(iResult int) bool {
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

		log.Printf("获取当前交易日  : %#v\n", p.Client.TraderApi.GetTradingDay())
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
		p.Client.login <- true
	} else {
		p.Client.login <- false
	}
}

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

	var id string = "c1901"
	req.SetInstrumentID(id)

	for {
		iResult := p.Client.TraderApi.ReqQryInstrument(req, p.Client.GetTraderRequestID())

		if !p.IsFlowControl(iResult) {
			log.Printf("--->>> 请求查询合约: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询合约: 受到流控 %#v\n", iResult)

		}
		time.Sleep(10 * time.Second)
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryInstrument(pInstrument goctp.CThostFtdcInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryInstrument: ", pInstrument.GetInstrumentID(), pInstrument.GetExchangeID(),
		pInstrument.GetInstrumentName(), pInstrument.GetExchangeInstID(), pInstrument.GetProductID(), pInstrument.GetProductClass(),
		pInstrument.GetDeliveryYear(), pInstrument.GetDeliveryMonth(), pInstrument.GetMaxMarketOrderVolume(), pInstrument.GetMinMarketOrderVolume(),
		pInstrument.GetMaxLimitOrderVolume(), pInstrument.GetMinLimitOrderVolume(), pInstrument.GetVolumeMultiple(), pInstrument.GetPriceTick(),
		pInstrument.GetCreateDate(), pInstrument.GetOpenDate(), pInstrument.GetExpireDate(), pInstrument.GetStartDelivDate(), pInstrument.GetEndDelivDate(),
		nRequestID, bIsLast)

	if bIsLast {

		///请求查询合约
		p.ReqQryTradingAccount()
	}
}

func (p *GoCThostFtdcTraderSpi) ReqQryTradingAccount() {

	req := goctp.NewCThostFtdcQryTradingAccountField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	for {
		iResult := p.Client.TraderApi.ReqQryTradingAccount(req, p.Client.GetTraderRequestID())

		if !p.IsFlowControl(iResult) {
			log.Printf("--->>> 请求查询资金账户: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询资金账户: 受到流控 %#v\n", iResult)

		}
		time.Sleep(10 * time.Second)
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryTradingAccount(pTradingAccount goctp.CThostFtdcTradingAccountField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("GoCThostFtdcTraderSpi.OnRspQryTradingAccount.")

	if bIsLast {
		///请求查询投资者持仓
		p.ReqQryInvestorPosition()
	}
}

func (p *GoCThostFtdcTraderSpi) ReqQryInvestorPosition() {

	req := goctp.NewCThostFtdcQryInvestorPositionField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetInvestorID(p.Client.InvestorID)

	for {
		iResult := p.Client.TraderApi.ReqQryInvestorPosition(req, p.Client.GetTraderRequestID())

		if !p.IsFlowControl(iResult) {
			log.Printf("--->>> 请求查询投资者持仓: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询投资者持仓: 受到流控 %#v\n", iResult)

		}
		time.Sleep(10 * time.Second)
	}
}

func (p *GoCThostFtdcTraderSpi) OnRspQryInvestorPosition(pInvestorPosition goctp.CThostFtdcInvestorPositionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcTraderSpi.OnRspQryInvestorPosition.")
	log.Println("position:", pInvestorPosition.GetInstrumentID(), pInvestorPosition.GetPosiDirection(),
		pInvestorPosition.GetPositionDate(), pInvestorPosition.GetOpenVolume(),
		pInvestorPosition.GetCloseVolume, pInvestorPosition.GetOpenAmount(), pInvestorPosition.GetCloseAmount(),
		pInvestorPosition.GetPositionCost(), pInvestorPosition.GetUseMargin())

	if bIsLast {
		// ///报单录入请求
		// ReqOrderInsert();
		// //执行宣告录入请求
		// ReqExecOrderInsert();
		// //询价录入
		// ReqForQuoteInsert();
		// //做市商报价录入
		// ReqQuoteInsert();
	}
}

func NewDirectorCThostFtdcMdSpi(v interface{}) goctp.CThostFtdcMdSpi {

	return goctp.NewDirectorCThostFtdcMdSpi(v)
}

type GoCThostFtdcMdSpi struct {
	Client GoCTPClient
}

func (p *GoCThostFtdcMdSpi) OnRspError(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("GoCThostFtdcMdSpi.OnRspError.")
	p.IsErrorRspInfo(pRspInfo)
}

func (p *GoCThostFtdcMdSpi) OnFrontDisconnected(nReason int) {
	log.Printf("GoCThostFtdcMdSpi.OnFrontDisconnected: %#v\n", nReason)
}

func (p *GoCThostFtdcMdSpi) OnHeartBeatWarning(nTimeLapse int) {
	log.Printf("GoCThostFtdcMdSpi.OnHeartBeatWarning: %v", nTimeLapse)
}

func (p *GoCThostFtdcMdSpi) OnFrontConnected() {
	log.Println("GoCThostFtdcMdSpi.OnFrontConnected.")
	p.ReqUserLogin()
}

func (p *GoCThostFtdcMdSpi) ReqUserLogin() {
	log.Println("GoCThostFtdcMdSpi.ReqUserLogin.")

	req := goctp.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(p.Client.BrokerID)
	req.SetUserID(p.Client.InvestorID)
	req.SetPassword(p.Client.Password)

	iResult := p.Client.MdApi.ReqUserLogin(req, p.Client.GetMdRequestID())

	if iResult != 0 {
		log.Println("发送用户登录请求: 失败.")
	} else {
		log.Println("发送用户登录请求: 成功.")
	}
}

func (p *GoCThostFtdcMdSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {
	// 如果ErrorID != 0, 说明收到了错误的响应
	bResult := (pRspInfo.GetErrorID() != 0)
	if bResult {
		log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
	}
	return bResult
}

func (p *GoCThostFtdcMdSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

		log.Printf("获取当前版本信息: %#v\n", goctp.CThostFtdcTraderApiGetApiVersion())
		log.Printf("获取当前交易日期: %#v\n", p.Client.MdApi.GetTradingDay())
		log.Printf("获取用户登录信息: %#v %#v %#v\n", pRspUserLogin.GetLoginTime(), pRspUserLogin.GetSystemName(), pRspUserLogin.GetSessionID())

		//ppInstrumentID := []string{"cu1610", "cu1611", "cu1612", "cu1701", "cu1702", "cu1703", "cu1704", "cu1705", "cu1706"}
		//	ppInstrumentID := []string{"cu1810", "cu1811", "cu1812"}

		//	p.SubscribeMarketData(ppInstrumentID)
		//	p.SubscribeForQuoteRsp(ppInstrumentID)
		p.Client.login <- true
	} else {
		p.Client.login <- false
	}
}

/*
func (p *GoCThostFtdcMdSpi) SubscribeMarketData(name []string) {

	iResult := p.Client.MdApi.SubscribeMarketData(name)

	if iResult != 0 {
		log.Println("发送行情订阅请求: 失败.")
	} else {
		log.Println("发送行情订阅请求: 成功.")
	}
}*/

func (p *GoCThostFtdcMdSpi) SubscribeForQuoteRsp(name []string) {

	iResult := p.Client.MdApi.SubscribeForQuoteRsp(name)

	if iResult != 0 {
		log.Println("发送询价订阅请求: 失败.")
	} else {
		log.Println("发送询价订阅请求: 成功.")
	}
}

func (p *GoCThostFtdcMdSpi) OnRspSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspSubMarketData: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
	p.IsErrorRspInfo(pRspInfo)
}

func (p *GoCThostFtdcMdSpi) OnRspSubForQuoteRsp(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspSubForQuoteRsp: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
	p.IsErrorRspInfo(pRspInfo)
}

func (p *GoCThostFtdcMdSpi) OnRspUnSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspUnSubMarketData: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
	p.IsErrorRspInfo(pRspInfo)
}

func (p *GoCThostFtdcMdSpi) OnRspUnSubForQuoteRsp(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspUnSubForQuoteRsp: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
	p.IsErrorRspInfo(pRspInfo)
}

func (p *GoCThostFtdcMdSpi) OnRtnDepthMarketData(pDepthMarketData goctp.CThostFtdcDepthMarketDataField) {

	log.Println("GoCThostFtdcMdSpi.OnRtnDepthMarketData: ", pDepthMarketData.GetTradingDay(),
		pDepthMarketData.GetInstrumentID(),
		pDepthMarketData.GetExchangeID(),
		pDepthMarketData.GetExchangeInstID(),
		pDepthMarketData.GetLastPrice(),
		pDepthMarketData.GetPreSettlementPrice(),
		pDepthMarketData.GetPreClosePrice(),
		pDepthMarketData.GetPreOpenInterest(),
		pDepthMarketData.GetOpenPrice(),
		pDepthMarketData.GetHighestPrice(),
		pDepthMarketData.GetLowestPrice(),
		pDepthMarketData.GetVolume(),
		pDepthMarketData.GetTurnover(),
		pDepthMarketData.GetOpenInterest())

	//log.Printf("GoCThostFtdcMdSpi.OnRtnDepthMarketData: %+v\n", &pDepthMarketData)

}

func (p *GoCThostFtdcMdSpi) OnRtnForQuoteRsp(pForQuoteRsp goctp.CThostFtdcForQuoteRspField) {
	log.Printf("GoCThostFtdcMdSpi.OnRtnForQuoteRsp: %#v\n", pForQuoteRsp)
}

func NewCtpClient(path string) (*CtpTran, error) {
	var CtpClient CtpTran

	conf, err := ini.Load(path) //加载配置文件
	if err != nil {
		log.Println("load config file fail!")
		return nil, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&CtpClient.MdClient)     //解析成结构体
	err = conf.MapTo(&CtpClient.TraderClient) //解析成结构体
	CtpClient.MdClient.MdApi = goctp.CThostFtdcMdApiCreateFtdcMdApi()
	CtpClient.TraderClient.TraderApi = goctp.CThostFtdcTraderApiCreateFtdcTraderApi()
	CtpClient.MdClient.login = make(chan bool)
	CtpClient.TraderClient.login = make(chan bool)

	//	pMdSpi := goctp.NewDirectorCThostFtdcMdSpi(&GoCThostFtdcMdSpi{Client: CtpClient.MdClient})
	pTraderSpi := goctp.NewDirectorCThostFtdcTraderSpi(&GoCThostFtdcTraderSpi{Client: CtpClient.TraderClient})

	CtpClient.TraderClient.TraderApi.RegisterSpi(pTraderSpi)
	CtpClient.TraderClient.TraderApi.SubscribePublicTopic(0)
	CtpClient.TraderClient.TraderApi.SubscribePrivateTopic(0)
	CtpClient.TraderClient.TraderApi.RegisterFront(CtpClient.TraderClient.TraderFront)
	CtpClient.TraderClient.TraderApi.Init()
	/*
		CtpClient.MdClient.MdApi.RegisterSpi(pMdSpi)
		CtpClient.MdClient.MdApi.RegisterFront(CtpClient.MdClient.MdFront)
		CtpClient.MdClient.MdApi.Init()
	*/
	return &CtpClient, nil
}
